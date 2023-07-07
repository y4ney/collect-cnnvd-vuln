package model

// HazardLevel 威胁等级
type HazardLevel struct {
	DictLabel string `json:"dictLabel,omitempty"`
	DictValue string `json:"dictValue,omitempty"`
}

// Vendor 供应商选择列表
type Vendor struct {
	Label string `json:"label,omitempty"`
	Value string `json:"value,omitempty"`
}
