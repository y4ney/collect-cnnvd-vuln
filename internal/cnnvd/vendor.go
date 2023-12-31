package cnnvd

import (
	"github.com/y4ney/collect-cnnvd-vuln/internal/model"
	"github.com/y4ney/collect-cnnvd-vuln/internal/utils"
	"golang.org/x/xerrors"
)

const (
	VendorPath = "web/homePage/getVendorSelectList"
)

// ReqVendor 供应商选择列表请求参数
type ReqVendor struct {
	VendorKeyword string `json:"vendorKeyword"` // 供应商关键词
}

// ResVendor 供应商选择列表响应参数
type ResVendor struct {
	ResCode                 // 响应码
	Data    []*model.Vendor `json:"data,omitempty"` // 供应商选择列表
}

func (r *ReqVendor) Fetch(retry int) ([]*model.Vendor, error) {
	// 获取供应商信息
	http := utils.HTTP{URL: utils.URL(Schema, Domain, VendorPath), Method: utils.Post, Retry: retry, Body: r}
	var res ResVendor
	if err := http.Fetch(&res); err != nil {
		return nil, xerrors.Errorf("failed to fetch vendor:%w", err)
	}

	var vendors []*model.Vendor
	for _, data := range res.Data {
		vendors = append(vendors, data)
	}

	return vendors, nil
}
