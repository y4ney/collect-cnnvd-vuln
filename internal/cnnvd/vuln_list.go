package cnnvd

import (
	"github.com/y4ney/collect-cnnvd-vuln/internal/model"
	"github.com/y4ney/collect-cnnvd-vuln/internal/utils"
	"golang.org/x/xerrors"
	"math"
)

const (
	FirstPage    = 1
	MaxPageSize  = 50
	VulnListPath = "web/homePage/cnnvdVulList"
	VulListFile  = "vuln_list"
)

// ReqVulList cnnvd漏洞列表请求参数
type ReqVulList struct {
	PageIndex   int    `json:"pageIndex"`   // 分页索引
	PageSize    int    `json:"pageSize"`    // 分页大小
	Keyword     string `json:"keyword"`     // 关键字
	HazardLevel string `json:"hazardLevel"` // 漏洞等级
	VulType     string `json:"vulType"`     // 漏洞类型
	Vendor      string `json:"vendor"`      // 供应商
	Product     string `json:"product"`     // 产品
	DateType    string `json:"dateType"`    // 数据类型
}

// ResVulList 供应商选择列表响应参数
type ResVulList struct {
	ResCode          // 响应码
	Data    VulnList `json:"data,omitempty"` // cnnvd漏洞列表
}

// VulnList cnnvd漏洞列表
type VulnList struct {
	Total     int             `json:"total,omitempty"`
	Records   []*model.Record `json:"records,omitempty"`
	PageIndex int             `json:"pageIndex,omitempty"`
	PageSize  int             `json:"pageSize,omitempty"`
}

func NewReqVulList(keyword string) *ReqVulList {
	return &ReqVulList{
		Keyword:  keyword,
		PageSize: MaxPageSize,
	}
}

func (r *ReqVulList) Fetch(retry int) ([]*model.Record, error) {
	// 获取漏洞列表
	http := utils.HTTP{URL: utils.URL(Schema, Domain, VulnListPath), Method: utils.Post, Retry: retry, Body: r}
	var res ResVulList
	if err := http.Fetch(&res); err != nil {
		return nil, xerrors.Errorf("failed to fetch vuln list:%w", err)
	}

	// 处理应答
	var vulns []*model.Record
	for _, record := range res.Data.Records {
		vulns = append(vulns, record)
	}
	return vulns, nil
}

func (r *ReqVulList) GetPageInfo(retry int) (loopNum int, total int, err error) {
	http := utils.HTTP{URL: utils.URL(Schema, Domain, VulnListPath), Method: utils.Post, Retry: retry, Body: r}

	var res ResVulList
	if err := http.Fetch(&res); err != nil {
		return 0, 0, xerrors.Errorf("failed to get page num:%w", err)
	}

	return int(math.Ceil(float64(res.Data.Total) / float64(res.Data.PageSize))), res.Data.Total, nil
}
