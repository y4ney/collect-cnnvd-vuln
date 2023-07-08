package cnnvd

import (
	"github.com/y4ney/collect-cnnvd-vuln/internal/model"
	"github.com/y4ney/collect-cnnvd-vuln/internal/utils"
	"golang.org/x/xerrors"
	"path/filepath"
)

const (
	HazardLevelPath = "web/dictionaries/type/hazardLevel"
	HazardLevelFile = "hazard_level.json"
)

// ReqHazardLevel 威胁等级请求参数
type ReqHazardLevel struct {
}

// ResHazardLevel 威胁等级响应参数
type ResHazardLevel struct {
	ResCode                      // 响应码
	Data    []*model.HazardLevel `json:"data,omitempty"` // 威胁等级列表
}

func (r *ReqHazardLevel) Fetch(retry int) ([]*model.HazardLevel, error) {
	http := utils.HTTP{URL: utils.URL(Schema, Domain, HazardLevelPath), Method: utils.Get, Retry: retry}
	var res ResHazardLevel
	if err := http.Fetch(&res); err != nil {
		return nil, xerrors.Errorf("failed to fetch hazard level:%w", err)
	}

	var hazardLevel []*model.HazardLevel
	for _, data := range res.Data {
		hazardLevel = append(hazardLevel, data)
	}
	return hazardLevel, nil
}

func (r *ReqHazardLevel) Save(data []*model.HazardLevel, dir string) error {
	path := filepath.Join(dir, HazardLevelFile)
	err := utils.WriteFile(path, data)
	if err != nil {
		return xerrors.Errorf("failed to save hazard level:%w", err)
	}
	return nil
}
