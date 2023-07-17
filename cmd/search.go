package cmd

import (
	"fmt"
	"github.com/cheggaaa/pb"
	"github.com/olekukonko/tablewriter"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/y4ney/collect-cnnvd-vuln/internal/cnnvd"
	"github.com/y4ney/collect-cnnvd-vuln/internal/utils"
	"os"
)

const (
	QueryVulnList = "vuln"
	QueryVendor   = "vendor"
	QueryProduct  = "product"
	RetryNum      = 5
)

var (
	Keyword string
	Retry   int

	PageIndex   int
	PageSize    int
	HazardLevel string
	Product     string
	Vendor      string
)

var (
	level     = map[string]string{"超危": "1", "高危": "2", "中危": "3", "低危": "4"}
	severity  = map[int]string{0: "未知", 1: "超危", 2: "高危", 3: "中危", 4: "低危"}
	searchCmd = &cobra.Command{
		Use:   "search",
		Short: "搜索 CNNVD 漏洞信息",
		Run:   runSearchVuln,
	}
)

func init() {
	searchCmd.Flags().StringVarP(&Type, "type", "t", QueryVulnList, fmt.Sprintf("指定类型，仅支持 %s, %s 和 %s", QueryVulnList, QueryVendor, QueryProduct))
	searchCmd.Flags().StringVarP(&Keyword, "keyword", "k", "", "指定搜索的关键词")
	searchCmd.Flags().IntVarP(&Retry, "retry", "r", RetryNum, "指定请求的重试次数")
	searchCmd.Flags().IntVar(&PageIndex, "page-index", cnnvd.FirstPage, "指定当前的页码，仅在 --type=vuln 时有效")
	searchCmd.Flags().IntVar(&PageSize, "page-size", cnnvd.MaxPageSize, "指定页数大小 仅在 --type=vuln 时有效")
	searchCmd.Flags().StringVar(&HazardLevel, "hazard-level", "", "指定威胁等级，仅支持 超危、高危、中危和低危, 仅在 --type=vuln 时有效")
	// TODO 优化
	searchCmd.Flags().StringVar(&Product, "product", "", "指定商品编号，仅在 --type=vuln 时有效，请先通过 --type=product 来获取商品编号")
	searchCmd.Flags().StringVar(&Vendor, "vendor", "", "指定供应商编号，仅在 --type=vuln 时有效，请先通过 --type=vendor 获取供应商编号 ")
	utils.BindFlags(searchCmd)
}
func runSearchVuln(_ *cobra.Command, _ []string) {
	switch Type {
	case QueryProduct:
		searchProduct()
	case QueryVendor:
		searchVendor()
	case QueryVulnList:
		searchVuln()
	default:
		log.Error().Msgf("type %s is not supported", Type)
	}
}

func searchProduct() {
	c := cnnvd.ReqProduct{ProductKeyword: Keyword}
	products, err := c.Fetch(Retry)
	if err != nil {
		log.Fatal().Str("keyword", Keyword).Int("retry", Retry).
			Msgf("failed to search products:%v", err)
	}
	if len(products) == 0 {
		log.Info().Str("keyword", Keyword).Int("retry", Retry).
			Msg("there is no record")
		return
	}
	var data [][]string
	for _, product := range products {
		data = append(data, []string{product.Label, product.Value})
	}
	printInfo([]string{"LABEL", "VALUE"}, data)
}

func searchVendor() {
	c := cnnvd.ReqVendor{VendorKeyword: Keyword}
	vendors, err := c.Fetch(Retry)
	if err != nil {
		log.Fatal().Str("keyword", Keyword).Int("retry", Retry).
			Msgf("failed to search vendors:%v", err)
	}
	if len(vendors) == 0 {
		log.Info().Str("keyword", Keyword).Int("retry", Retry).
			Msg("there is no record")
		return
	}
	var data [][]string
	for _, product := range vendors {
		data = append(data, []string{product.Label, product.Value})
	}
	printInfo([]string{"LABEL", "VALUE"}, data)
}

func searchVuln() {
	c := cnnvd.ReqVulList{
		PageIndex:   PageIndex,
		PageSize:    PageSize,
		Keyword:     Keyword,
		HazardLevel: level[HazardLevel],
		Vendor:      Vendor,
		Product:     Product,
	}
	vulns, err := c.Fetch(Retry)
	if err != nil {
		log.Fatal().Interface("request", c).Int("retry", Retry).
			Msgf("failed to search vulns:%v", err)
	}
	log.Info().Interface("request", c).Int("retry", Retry).
		Msg("success to request... ...")
	if len(vulns) == 0 {
		log.Info().Interface("request", c).Int("retry", Retry).
			Msg("there is no record")
		return
	}
	bar := pb.StartNew(len(vulns))
	var data [][]string
	for _, vuln := range vulns {
		detailC := cnnvd.ReqVulDetail{Id: vuln.Id, VulType: vuln.VulType, CnnvdCode: vuln.CnnvdCode}
		detail, err := detailC.Fetch(Retry)
		log.Debug().Interface("request", detailC).Int("retry", Retry).
			Msg("success to request... ...")
		if err != nil {
			log.Fatal().Interface("request", detailC).Int("retry", Retry).
				Msgf("failed to search vuln detail:%v", err)
		}
		data = append(data, []string{severity[detail.HazardLevel], detail.CnnvdCode, detail.CveCode, detail.VulName,
			detail.VulTypeName, detail.AffectedVendor, detail.AffectedProduct, detail.UpdateTime})
		bar.Increment()
	}
	bar.Finish()
	printInfo([]string{"SEVERITY", "CNNVD ID", "CVE ID", "NAME", "TYPE", "VENDOR", "PRODUCT", "UPDATE TIME"}, data)

}

func printInfo(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	for _, row := range data {
		colors := make([]tablewriter.Colors, len(row))
		for i, cell := range row {
			switch cell {
			case severity[0]:
				colors[i] = tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlackColor}
			case severity[1]:
				colors[i] = tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor}
			case severity[2]:
				colors[i] = tablewriter.Colors{tablewriter.Bold, tablewriter.FgRedColor}
			case severity[3]:
				colors[i] = tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiYellowColor}
			case severity[4]:
				colors[i] = tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlueColor}
			default:
				colors[i] = tablewriter.Colors{}
			}
		}
		table.Rich(row, colors)
	}
	table.Render()
}
