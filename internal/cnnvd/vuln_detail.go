package cnnvd

import (
	"github.com/y4ney/collect-cnnvd-vuln/internal/model"
	"github.com/y4ney/collect-cnnvd-vuln/internal/utils"
	"golang.org/x/xerrors"
	"path/filepath"
)

const (
	VulnDetailPath = "web/cnnvdVul/getCnnnvdDetailOnDatasource"
	VulDetailFile  = "vul_detail"
)

// ReqVulDetail cnnvd漏洞详情请求参数
type ReqVulDetail struct {
	Id        string `json:"id"`        // 漏洞id
	VulType   string `json:"vulType"`   // 漏洞类型
	CnnvdCode string `json:"cnnvdCode"` // cnnvd编号
}

// ResVulDetail cnnvd漏洞详情响应参数
type ResVulDetail struct {
	ResCode            // 响应码
	Data    *VulDetail `json:"data,omitempty"` // 漏洞详情数据
}

// VulDetail 漏洞详情数据
type VulDetail struct {
	model.CNNVDDetail `json:"cnnvdDetail,omitempty"` // CNNVD详情
	ReceviceVulDetail string                         `json:"receviceVulDetail,omitempty"` // 接收到的漏洞详情
}

func (r *ReqVulDetail) Fetch(retry int) (*VulDetail, error) {
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

func (r *ReqVulDetail) Save(data *VulDetail, dir string) error {
	// 创建目录
	path := filepath.Join(dir, VulDetailFile)
	if err := utils.Mkdir(path); err != nil {
		return xerrors.Errorf("failed to mkdir %s:%w", path, err)
	}

	// 写入文件
	if err := SaveCNNVDPerYear(path, data.CnnvdCode, data); err != nil {
		return xerrors.Errorf("failed to save %s:%w", data.CnnvdCode, err)
	}
	return nil
}
