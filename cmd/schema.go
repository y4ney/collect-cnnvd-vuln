package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/invopop/jsonschema"
	"github.com/spf13/cobra"
	"github.com/y4ney/collect-cnnvd-vuln/internal/model"
	"github.com/y4ney/collect-cnnvd-vuln/internal/utils"
	"golang.org/x/xerrors"
)

const (
	VulnListSchema    = "vuln-list"
	VulnDetailSchema  = "vuln-detail"
	HazardLevelSchema = "hazard-level"
	VendorSchema      = "vendor"
	ProductSchema     = "product"
	VulnTypeSchema    = "vuln-type"
)

var Type string

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "打印 CNNVD 漏洞信息的json文件结构",
	RunE:  runGenerateSchema,
}

func init() {
	schemaCmd.Flags().StringVarP(&Type, "type", "t", VulnDetailSchema,
		fmt.Sprintf("指定数据类型(只支持 %s, %s, %s, %s, %s 和 %s)",
			VulnListSchema, VulnDetailSchema, HazardLevelSchema, VendorSchema, ProductSchema, VulnTypeSchema),
	)

	utils.BindFlags(schemaCmd)
}
func runGenerateSchema(_ *cobra.Command, _ []string) error {
	var schema *jsonschema.Schema
	switch Type {
	case HazardLevelSchema:
		schema = jsonschema.Reflect(&model.HazardLevel{})
	case VendorSchema:
		schema = jsonschema.Reflect(&model.Vendor{})
	case ProductSchema:
		schema = jsonschema.Reflect(&model.Product{})
	case VulnTypeSchema:
		schema = jsonschema.Reflect(&model.VulnType{})
	case VulnListSchema:
		schema = jsonschema.Reflect(&model.Record{})
	case VulnDetailSchema:
		schema = jsonschema.Reflect(&model.VulDetail{})
	default:
		return xerrors.Errorf("type %s is not supported", Type)
	}
	enc := json.NewEncoder(out)
	enc.SetIndent("", "\t")

	return enc.Encode(schema)
}
