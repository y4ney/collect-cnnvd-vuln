package cnnvd

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"golang.org/x/xerrors"
	"strconv"
	"strings"
	"time"
)

const (
	Schema = "https"
	Domain = "www.cnnvd.org.cn"
)

type CnnvdClient interface {
	Fetch(retry int) (any, error)
}

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

func (c *CNNVD) GetDate() *time.Time {
	date, err := time.Parse("2006-01", fmt.Sprintf("%v-%02d", c.Year, c.Month))
	if err != nil {
		log.Fatal().Interface("CNNVD ID", c).Msgf("fail to get date:%v", err)
	}
	return &date
}

func (c *CNNVD) After(item *CNNVD) bool {
	var (
		cDate    = c.GetDate()
		itemDate = item.GetDate()
	)
	if cDate.Before(*itemDate) {
		return false
	}
	if cDate.After(*itemDate) {
		return true
	}
	if c.ID > item.ID {
		return true
	}
	return false
}

func (c *CNNVD) Equal(item *CNNVD) bool {
	return c.Year == item.Year && c.Month == item.Month && c.ID == item.ID
}
