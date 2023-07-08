package cnnvd

import (
	"github.com/y4ney/collect-cnnvd-vuln/internal/model"
	"github.com/y4ney/collect-cnnvd-vuln/internal/utils"
	"golang.org/x/xerrors"
	"path/filepath"
)

const (
	ProductFile = "product.json"
	ProductPath = "web/homePage/getProductSelectList"
)

// ReqProduct 产品选择列表请求参数
type ReqProduct struct {
	ProductKeyword string `json:"productKeyword"` // 产品关键词
}

// ResProduct 产品选择列表响应参数
type ResProduct struct {
	ResCode                  // 响应码
	Data    []*model.Product `json:"data,omitempty"` // 产品列表
}

func (r *ReqProduct) Fetch(retry int) ([]*model.Product, error) {
	// 获取产品信息
	http := utils.HTTP{URL: utils.URL(Schema, Domain, ProductPath), Method: utils.Post, Retry: retry, Body: r}
	var res ResProduct
	if err := http.Fetch(&res); err != nil {
		return nil, xerrors.Errorf("failed to fetch product:%w", err)
	}

	var products []*model.Product
	for _, data := range res.Data {
		products = append(products, data)
	}

	_ = utils.WriteFile("./testdata/product.json", products)
	return products, nil
}

func (r *ReqProduct) Save(data []*model.Product, dir string) error {
	path := filepath.Join(dir, ProductFile)
	err := utils.WriteFile(path, data)
	if err != nil {
		return xerrors.Errorf("failed to save product:%w", err)
	}
	return nil
}
