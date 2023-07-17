package cnnvd

import (
	"github.com/y4ney/collect-cnnvd-vuln/internal/model"
	"github.com/y4ney/collect-cnnvd-vuln/internal/utils"
	"golang.org/x/xerrors"
)

const (
	VulnTypePath = "web/homePage/vulTypeList"
)

// ReqVulType 漏洞类型列表请求参数
type ReqVulType struct {
}

// ResVulType 漏洞类型列表响应参数
type ResVulType struct {
	ResCode                   // 响应码
	Data    []*model.VulnType `json:"data,omitempty"` // 漏洞类型列表
}

func (r *ReqVulType) Fetch(retry int) ([]*model.VulnType, error) {
	// 获取产品信息
	http := utils.HTTP{URL: utils.URL(Schema, Domain, VulnTypePath), Method: utils.Post, Retry: retry, Body: r}
	var res ResVulType
	if err := http.Fetch(&res); err != nil {
		return nil, xerrors.Errorf("failed to fetch vuln type:%w", err)
	}

	var vulnTypes []*model.VulnType
	for _, data := range res.Data {
		vulnTypes = append(vulnTypes, data)
	}

	_ = utils.WriteFile("./testdata/vuln-type.json", vulnTypes)
	return vulnTypes, nil
}
