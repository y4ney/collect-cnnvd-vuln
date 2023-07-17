package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/y4ney/collect-cnnvd-vuln/internal/model"
	"golang.org/x/xerrors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func WriteFile(filepath string, data any) error {
	d, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return xerrors.Errorf("failed to marshal:%w", err)
	}
	if err = os.WriteFile(filepath, d, os.ModePerm); err != nil {
		return xerrors.Errorf("failed to write to %s:%w", filepath, err)
	}

	return nil
}

func ReadFile(filepath string, data any) error {
	d, err := os.ReadFile(filepath)
	if err != nil {
		return xerrors.Errorf("failed to read %s:%w", filepath, err)
	}

	if err = json.Unmarshal(d, data); err != nil {
		return xerrors.Errorf("failed to unmarshal %s:%w", filepath, err)
	}

	return nil
}

func DeleteFile(filepath string) error {
	// 若文件不存在，则正常退出
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return nil
	}

	// 若文件存在，则删除
	err = os.Remove(filepath)
	if err != nil {
		return xerrors.Errorf("failed to remove %s:%w", filepath, err)
	}
	return nil
}

func SizeOfFile(filepath string) (string, error) {
	// 获取文件信息
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return "", xerrors.Errorf("failed to get file info:%w", err)
	}
	return formatFileSize(fileInfo.Size()), nil
}

func Mkdir(dir string) error {
	_, err := os.Stat(dir)
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}

	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func CacheDir() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		os.TempDir()
	}
	return filepath.Join(cacheDir, "y4ney", "cnnvd")
}
func GetCNNVDFromFile(dir string) ([]*model.VulDetail, error) {
	var items []*model.VulDetail
	buffer := &bytes.Buffer{}
	err := FileWalk(dir, func(r io.Reader, _ string) error {
		var item model.VulDetail
		if _, err := buffer.ReadFrom(r); err != nil {
			return xerrors.Errorf("failed to read file: %w", err)
		}
		if err := json.Unmarshal(buffer.Bytes(), &item); err != nil {
			return xerrors.Errorf("failed to decode vuln detail JSON: %w", err)
		}
		buffer.Reset()
		items = append(items, &item)
		return nil
	})
	if err != nil {
		return nil, xerrors.Errorf("error in vuln detail walk: %w", err)
	}
	return items, nil
}

// FileWalk 遍历文件，并执行walkFn的操作
func FileWalk(root string, walkFn func(r io.Reader, path string) error) error {
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		} else if d.IsDir() {
			// 若是目录，则正常退出
			return nil
		}

		// 获取文件信息，若文件大小为0，则正常退出
		info, err := d.Info()
		if err != nil {
			return xerrors.Errorf("file info error: %w", err)
		}
		if info.Size() == 0 {
			log.Warn().Str("path", path).Msg("invalid size")
			return nil
		}

		if filepath.Ext(path) != ".json" {
			return nil
		}

		// 打开文件，执行walkFn操作
		f, err := os.Open(path)
		if err != nil {
			return xerrors.Errorf("failed to open file: %w", err)
		}
		defer f.Close()
		if err = walkFn(f, path); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return xerrors.Errorf("failed to file walk: %w", err)
	}
	return nil
}

// formatFileSize 格式化文件大小为可读格式
func formatFileSize(size int64) string {
	const (
		B  = 1
		KB = 1024 * B
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case size < KB:
		return fmt.Sprintf("%d B", size)
	case size < MB:
		return fmt.Sprintf("%.2f KB", float64(size)/KB)
	case size < GB:
		return fmt.Sprintf("%.2f MB", float64(size)/MB)
	default:
		return fmt.Sprintf("%.2f GB", float64(size)/GB)
	}
}

func FormatCnnvdRef(refLink string) string {
	var Links []string
	for _, link := range strings.Split(refLink, "\n") {
		if strings.Contains(link, "链接:http") {

			Links = append(Links, strings.TrimPrefix(strings.Trim(strings.Trim(link, "\n"), "\r"), "链接:"))
		}
	}
	return strings.Join(Links, ",")
}
