package trivy

import (
	"context"
	"fmt"
	"github.com/aquasecurity/trivy-db/pkg/metadata"
	"github.com/aquasecurity/trivy/pkg/db"
	"github.com/aquasecurity/trivy/pkg/fanal/types"
	bolt "go.etcd.io/bbolt"
	"golang.org/x/xerrors"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultDBRepository = "ghcr.io/aquasecurity/trivy-db" // 默认的数据库仓库
	vulnBucket          = "vulnerability"
	DBFile              = "trivy.db"
	DB                  = "db"
)

func DownloadTrivyDB(dir string) (*metadata.Metadata, error) {
	// 创建客户端
	client := db.NewClient(dir, false, db.WithDBRepository(defaultDBRepository))

	// 判断数据是否需要更新
	needsUpdate, err := client.NeedsUpdate("dev", false)

	// 下载漏洞库
	if needsUpdate {
		if err := client.Download(context.Background(), dir, types.RegistryOptions{}); err != nil {
			return nil, xerrors.Errorf("failed to download vulnerability DB: %v", err)
		}
	}

	// 获取元数据
	m := metadata.NewClient(dir)
	meta, err := m.Get()
	if err != nil {
		return nil, xerrors.Errorf("something wrong with DB: %w", err)
	}
	return &meta, nil
}

func GetCvdIdFromTrivyDB(dir string) ([]string, error) {
	// 构建数据库文件路径
	path := filepath.Join(dir, DB, DBFile)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, xerrors.Errorf("%s is not exist:%w", path, err)
	}

	// 打开数据库文件
	trivyDB, err := bolt.Open(path, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer trivyDB.Close()

	// 读取数据
	var CveIds []string
	if err = trivyDB.View(func(tx *bolt.Tx) error {
		// 打开 bucket（如果不存在则创建）
		bucket := tx.Bucket([]byte(vulnBucket))
		if bucket == nil {
			return fmt.Errorf("bucket %s is not exist", vulnBucket)
		}

		// 遍历并获取 bucket 中的全部key
		if err = bucket.ForEach(func(k, v []byte) error {
			if strings.Contains(string(k), "CVE") {
				CveIds = append(CveIds, string(k))
			}
			return nil
		}); err != nil {
			return xerrors.Errorf("failed to for each %s:%w", vulnBucket, err)
		}

		return nil
	}); err != nil {
		return nil, xerrors.Errorf("failed to view %s:%w", trivyDB, err)
	}
	return CveIds, nil
}
