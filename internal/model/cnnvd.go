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

// Product 产品列表
type Product struct {
	Label string `json:"label,omitempty" gorm:"uniqueIndex size:256"`
	Value string `json:"value,omitempty"`
}

// VulnType 漏洞类型列表
type VulnType struct {
	Id      string     `json:"id,omitempty"`
	Pid     string     `json:"pid,omitempty"`
	Label   string     `json:"label,omitempty"`
	Value   string     `json:"value,omitempty"`
	VulType []VulnType `json:"children,omitempty"`
}

// Record 漏洞列表记录
type Record struct {
	Id          string `json:"id,omitempty"`
	VulName     string `json:"vulName,omitempty"`
	CnnvdCode   string `json:"cnnvdCode,omitempty"`
	CveCode     string `json:"cveCode,omitempty"`
	HazardLevel int64  `json:"hazardLevel,omitempty"`
	CreateTime  string `json:"createTime,omitempty"`
	PublishTime string `json:"publishTime,omitempty"`
	UpdateTime  string `json:"updateTime,omitempty"`
	TypeName    string `json:"typeName,omitempty"`
	VulType     string `json:"vulType,omitempty"`
}

// VulDetail 漏洞详情数据
type VulDetail struct {
	CNNVDDetail       `json:"cnnvdDetail,omitempty"`              // CNNVD详情
	ReceviceVulDetail string `json:"receviceVulDetail,omitempty"` // 接收到的漏洞详情
}

// CNNVDDetail CNNVD详情
type CNNVDDetail struct {
	Id                 string `json:"id,omitempty"`
	VulName            string `json:"vulName,omitempty"`
	CnnvdCode          string `json:"cnnvdCode,omitempty"`
	CveCode            string `json:"cveCode,omitempty"`
	PublishTime        string `json:"publishTime,omitempty"`
	IsOfficial         int    `json:"isOfficial,omitempty"`
	Vendor             string `json:"vendor,omitempty"`
	HazardLevel        int    `json:"hazardLevel,omitempty"`
	VulType            string `json:"vulType,omitempty"`
	VulTypeName        string `json:"vulTypeName,omitempty"`
	VulDesc            string `json:"vulDesc"`
	AffectedProduct    string `json:"affectedProduct,omitempty"`
	AffectedVendor     string `json:"affectedVendor,omitempty"`
	ProductDesc        string `json:"productDesc,omitempty"`
	AffectedSystem     string `json:"affectedSystem,omitempty"`
	ReferUrl           string `json:"referUrl,omitempty"`
	PatchId            string `json:"patchId,omitempty"`
	Patch              string `json:"patch,omitempty"`
	Deleted            string `json:"deleted,omitempty"`
	Version            string `json:"version,omitempty"`
	CreateUid          string `json:"createUid,omitempty"`
	CreateUname        string `json:"createUname,omitempty"`
	CreateTime         string `json:"createTime,omitempty"`
	UpdateUid          string `json:"updateUid,omitempty"`
	UpdateUname        string `json:"updateUname,omitempty"`
	UpdateTime         string `json:"updateTime,omitempty"`
	CnnvdFiledShow     string `json:"cnnvdFiledShow,omitempty"`
	CveVulVO           string `json:"cveVulVO,omitempty"`
	CveFiledShow       string `json:"cveFiledShow,omitempty"`
	IbmVulVO           string `json:"ibmVulVO,omitempty"`
	IbmFiledShow       string `json:"ibmFiledShow,omitempty"`
	IcsCertVulVO       string `json:"icsCertVulVO,omitempty"`
	IcsCertFiledShow   string `json:"icsCertFiledShow,omitempty"`
	MicrosoftVulVO     string `json:"microsoftVulVO,omitempty"`
	MicrosoftFiledShow string `json:"microsoftFiledShow,omitempty"`
	HuaweiVulVO        string `json:"huaweiVulVO,omitempty"`
	HuaweiFiledShow    string `json:"huaweiFiledShow,omitempty"`
	NvdVulVO           string `json:"nvdVulVO,omitempty"`
	NvdFiledShow       string `json:"nvdFiledShow,omitempty"`
	Varchar1           string `json:"varchar1,omitempty"`
}
