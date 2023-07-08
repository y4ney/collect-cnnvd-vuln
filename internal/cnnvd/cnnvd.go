package cnnvd

import (
	"fmt"
	"github.com/y4ney/collect-cnnvd-vuln/internal/utils"
	"golang.org/x/xerrors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	Schema = "https"
	Domain = "www.cnnvd.org.cn"
)

// ResCode 响应码
type ResCode struct {
	Code    int    `json:"code,omitempty"`    // 代码
	Success bool   `json:"success,omitempty"` // 是否成功
	Message string `json:"message,omitempty"` // 消息
	Time    string `json:"time,omitempty"`    // 时间
}

type CNNVD struct {
	Year  int
	Month int
	ID    int
}

func NewCNNVD(str string) (*CNNVD, error) {
	s := strings.Split(str, "-")
	if len(s) != 3 {
		return nil, xerrors.Errorf("invalid CNNVD-ID format: %s", str)
	}

	id, err := strconv.Atoi(s[2])
	if err != nil {
		return nil, xerrors.Errorf("failed to convert %s's id:%w", str, err)
	}
	year, err := strconv.Atoi(s[1][:4])
	if err != nil {
		return nil, xerrors.Errorf("failed to convert %s's year:%w", str, err)
	}
	month, err := strconv.Atoi(s[1][4:])
	if err != nil {
		return nil, xerrors.Errorf("failed to convert %s's month:%w", str, err)
	}
	return &CNNVD{Year: year, Month: month, ID: id}, nil
}

func (c *CNNVD) FormatCNNVD() (string, error) {
	if c.Year == 0 || c.Month < 0 || c.Month > 12 {
		return "", xerrors.Errorf("invalid CNNVD-ID format")
	}
	cnnvd := fmt.Sprintf("CNNVD-%d", c.Year)
	if c.Month != 0 {
		cnnvd += fmt.Sprintf("%02d", c.Month)
	}
	return cnnvd, nil
}

func (c *CNNVD) GetDate() (*time.Time, error) {
	date, err := time.Parse("2006-01", fmt.Sprintf("%v-%v", c.Year, c.Month))
	if err != nil {
		return nil, xerrors.Errorf("fail to get date:%w", err)
	}
	return &date, nil
}

func (c *CNNVD) Before(item *CNNVD) (bool, error) {
	cDate, err := c.GetDate()
	if err != nil {
		return false, err
	}
	itemDate, err := item.GetDate()
	if err != nil {
		return false, err
	}
	if cDate.After(*itemDate) {
		return false, nil
	}
	if cDate.Before(*itemDate) {
		return true, nil
	}
	if c.ID < item.ID {
		return true, nil
	}
	return false, nil

}

func (c *CNNVD) After(item *CNNVD) (bool, error) {
	cDate, err := c.GetDate()
	if err != nil {
		return false, err
	}
	itemDate, err := item.GetDate()
	if err != nil {
		return false, err
	}
	if cDate.Before(*itemDate) {
		return false, nil
	}
	if cDate.After(*itemDate) {
		return true, nil
	}
	if c.ID > item.ID {
		return true, nil
	}
	return false, nil
}

func (c *CNNVD) Equal(item *CNNVD) bool {
	return c.Year == item.Year && c.Month == item.Month && c.ID == item.ID
}

func LatestCNNVD(str1, str2 string) (string, error) {
	cnnvd1, err := NewCNNVD(str1)
	if err != nil {
		return "", xerrors.Errorf("fail to new %s:%w\n", str1, err)
	}
	cnnvd2, err := NewCNNVD(str2)
	if err != nil {
		return "", xerrors.Errorf("fail to new %s:%w\n", str2, err)
	}
	flag, err := cnnvd1.After(cnnvd2)
	if err != nil {
		return "", err
	}
	if flag {
		return str1, nil
	}
	return str2, nil
}

// SaveCNNVDPerYear 存储每年的漏洞
func SaveCNNVDPerYear(dirPath string, cnnvdID string, data interface{}) error {
	// 创建 cnnvd 对象
	cnnvd, err := NewCNNVD(cnnvdID)
	if err != nil {
		return xerrors.Errorf("failed to new %s:%w", cnnvdID, err)
	}

	// 根据年和月创建目录
	yearDir := filepath.Join(dirPath, strconv.Itoa(cnnvd.Year))
	monthDir := filepath.Join(yearDir, strconv.Itoa(cnnvd.Month))
	if err = os.MkdirAll(monthDir, os.ModePerm); err != nil {
		return err
	}

	// 写入文件
	filePath := filepath.Join(monthDir, fmt.Sprintf("%s.json", cnnvdID))
	if err = utils.WriteFile(filePath, data); err != nil {
		return xerrors.Errorf("failed to write %s: %w", filePath, err)
	}
	return nil
}
