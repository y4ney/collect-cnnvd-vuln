package cnnvd

import (
	"github.com/y4ney/collect-cnnvd-vuln/internal/model"
	"github.com/y4ney/collect-cnnvd-vuln/internal/utils"
	"golang.org/x/xerrors"
)

const (
	VulnDetailPath = "web/cnnvdVul/getCnnnvdDetailOnDatasource"
	VulDetailFile  = "vuln_detail"
)

// ReqVulDetail cnnvd漏洞详情请求参数
type ReqVulDetail struct {
	Id        string `json:"id"`        // 漏洞id
	VulType   string `json:"vulType"`   // 漏洞类型
	CnnvdCode string `json:"cnnvdCode"` // cnnvd编号
}

// ResVulDetail cnnvd漏洞详情响应参数
type ResVulDetail struct {
	ResCode                  // 响应码
	Data    *model.VulDetail `json:"data,omitempty"` // 漏洞详情数据
}

func (r *ReqVulDetail) Fetch(retry int) (*model.VulDetail, error) {
	if r.Id == "" || r.VulType == "" || r.CnnvdCode == "" {
		return nil, xerrors.New("please specify id,vuln type and cnnvd code")
	}
	http := utils.HTTP{URL: utils.URL(Schema, Domain, VulnDetailPath), Method: utils.Post, Retry: retry, Body: r}
	var res ResVulDetail
	if err := http.Fetch(&res); err != nil {
		return nil, xerrors.Errorf("failed to request %s detail:%w", r.CnnvdCode, err)
	}
	return res.Data, nil
}
