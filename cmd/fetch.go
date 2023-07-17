package cmd

import (
	"fmt"
	"github.com/cheggaaa/pb"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/y4ney/collect-cnnvd-vuln/internal/cnnvd"
	"github.com/y4ney/collect-cnnvd-vuln/internal/model"
	"github.com/y4ney/collect-cnnvd-vuln/internal/utils"
	"golang.org/x/xerrors"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

const (
	HazardLevelFile = "hazard_level.json"
	ProductFile     = "product.json"
	VendorFile      = "vendor.json"
	VulnTypeFile    = "vuln_type.json"
	MetadataFile    = "metadata.json"

	MinYear  = 1988
	MinMonth = 1
	MaxMonth = 12

	IncrementalUpdate = "increment"
	AllUpdate         = "all"
	Specific          = "specific"
	OldestCnnvdId     = "CNNVD-198801-001"
)

var (
	Year     int
	Month    int
	Dir      string
	Restart  bool
	Complete bool
	Github   bool
)

var (
	fetchCmd = &cobra.Command{
		Use:   "fetch",
		Short: "收集 CNNVD 漏洞信息",
		Run:   runFetch,
	}
	ThisYear     = time.Now().Year()
	ThisMonth, _ = utils.FormatMonth(time.Now().Month())
)

func init() {
	fetchCmd.Flags().IntVarP(&Year, "year", "y", ThisYear, "指定你想收集的年份的漏洞，仅在 --type=specific  时有效")
	fetchCmd.Flags().IntVarP(&Month, "month", "m", ThisMonth, "指定你想收集的月份的漏洞，仅在 --type=specific  时有效")
	fetchCmd.Flags().StringVarP(&Dir, "dir", "d", utils.CacheDir(), "指定数据的缓存目录")
	fetchCmd.Flags().StringVarP(&Type, "type", "t", AllUpdate, fmt.Sprintf("指定更新的类型，仅支持 %s、%s 和 %s", IncrementalUpdate, AllUpdate, Specific))
	fetchCmd.Flags().BoolVar(&Restart, "restart", false, "重新从CNNVD中收集漏洞信息")
	fetchCmd.Flags().BoolVarP(&Complete, "complete", "c", true, "收集漏洞详情")
	fetchCmd.Flags().IntVarP(&Retry, "retry", "r", RetryNum, "指定发送请求的次数")
	fetchCmd.Flags().BoolVarP(&Github, "github", "g", true, "是否上传到 Github 仓库中")

}
func runFetch(cmd *cobra.Command, _ []string) {
	var (
		metaPath      = filepath.Join(Dir, MetadataFile)
		totalVuln     = 0
		latestVulnId  = OldestCnnvdId
		latestVuln, _ = cnnvd.NewCNNVD(latestVulnId)
	)

	var newCnnvd *cnnvd.CNNVD
	getVulnExtraInfo()
	for _, feed := range getFeeds() {
		total, latestCnnvdId := getVulnInfo(feed)
		totalVuln += total
		newCnnvd, _ = cnnvd.NewCNNVD(latestCnnvdId)
		if newCnnvd.After(latestVuln) {
			latestVuln = newCnnvd
			latestVulnId = latestCnnvdId
		}
		log.Debug().Interface("Keyword", feed).Str("Latest CNNVD", latestVulnId).Int("Total Vuln", totalVuln).
			Msg("success to fetch vuln info")
	}
	utils.WriteMetaFile(metaPath, totalVuln, latestVulnId)

	v := viper.New()
	v.AutomaticEnv()
	git := utils.Git{
		URL:        v.GetString("CNNVD_URL"),
		Dir:        Dir,
		RemoteName: v.GetString("REMOTE_NAME"),
		Name:       v.GetString("NAME"),
		Email:      v.GetString("EMAIL"),
		Token:      v.GetString("TOKEN"),
	}
	if err := git.Add(); err != nil {
		log.Fatal().Interface("Git", git).Msgf("failed to git add:%v", err)
	}
	if err := git.Commit(); err != nil {
		log.Fatal().Interface("Git", git).Msgf("failed to git commit:%v", err)
	}
	if err := git.Push(); err != nil {
		log.Fatal().Interface("Git", git).Msgf("failed to git push:%v", err)
	}
}

func getFeeds() []*cnnvd.CNNVD {
	var (
		feeds   []*cnnvd.CNNVD
		month   int
		cnnvdId cnnvd.CNNVD
	)

	switch Type {
	case AllUpdate:
		for y := MinYear; y <= ThisYear; y++ {
			if y != ThisYear {
				month = MaxMonth
			} else {
				month = ThisMonth
			}
			for m := MinMonth; m <= month; m++ {
				feeds = append(feeds, &cnnvd.CNNVD{Year: y, Month: m})
			}
		}
	case Specific:
		if Month != ThisMonth {
			feeds = append(feeds, &cnnvd.CNNVD{Year: Year, Month: Month})

		} else {
			for m := MinMonth; m <= MaxMonth; m++ {
				feeds = append(feeds, &cnnvd.CNNVD{Year: Year, Month: m})
			}
		}

	case IncrementalUpdate:
		var metadata model.Metadata
		path := filepath.Join(Dir, MetadataFile)
		fileInfo, err := os.Stat(path)
		if err != nil {
			log.Fatal().Str("Path", path).Msgf("failed to get the size:%v", err)
		}
		if fileInfo.Size() == 0 {
			log.Fatal().Str("Path", path).Msg("metadata file ie empty,please run with --type=specific")
		}
		err = utils.ReadFile(path, &metadata)
		if err != nil {
			log.Fatal().Str("Path", path).Msgf("failed to read metadata:%v", err)
		}
		if time.Now().After(metadata.NextIncrementUpdate) {
			lastCnnvd, _ := cnnvd.NewCNNVD(metadata.LatestCnnvd)
			var start, end int
			for y := lastCnnvd.Year; y <= ThisYear; y++ {
				if y == lastCnnvd.Year {
					start = lastCnnvd.Month
				} else {
					start = MinMonth
				}
				if y == ThisYear {
					end = ThisMonth
				} else {
					end = MaxMonth
				}
				for m := start; m < end; m++ {
					cnnvdId.Year = y
					cnnvdId.Month = m
					feeds = append(feeds, &cnnvd.CNNVD{Year: y, Month: m})
				}
			}
		}
	}
	return feeds
}

func getVulnExtraInfo() {
	log.Info().Str("Directory", Dir).Msg("Saving CNNVD extra data...")
	// 下载漏洞等级
	bar := pb.StartNew(3)
	var reqHazardLevel cnnvd.ReqHazardLevel
	hazardLevel, err := reqHazardLevel.Fetch(Retry)
	if err != nil {
		log.Fatal().Msgf("failed to fetch hazard level:%w", err)
	}
	err = utils.WriteFile(filepath.Join(Dir, HazardLevelFile), hazardLevel)
	if err != nil {
		log.Fatal().Str("Hazard level file", filepath.Join(Dir, HazardLevelFile)).
			Msgf("failed to write hazard level:%w", err)
	}
	bar.Increment()

	// 下载产品信息
	var reqProduct cnnvd.ReqProduct
	product, err := reqProduct.Fetch(Retry)
	if err != nil {
		log.Fatal().Msgf("failed to fetch product:%w", err)
	}
	err = utils.WriteFile(filepath.Join(Dir, ProductFile), product)
	if err != nil {
		log.Fatal().Str("Product file", filepath.Join(Dir, ProductFile)).
			Msgf("failed to write product:%w", err)
	}
	bar.Increment()

	// 下载供应商信息
	var reqVendor cnnvd.ReqVendor
	vendor, err := reqVendor.Fetch(Retry)
	if err != nil {
		log.Fatal().Msgf("failed to fetch vendor:%w", err)
	}
	err = utils.WriteFile(filepath.Join(Dir, VendorFile), vendor)
	if err != nil {
		log.Fatal().Str("Vendor file", filepath.Join(Dir, VendorFile)).
			Msgf("failed to write vendor:%w", err)
	}
	bar.Increment()

	// 下载漏洞类型
	var reqVulnType cnnvd.ReqVulType
	vulnType, err := reqVulnType.Fetch(Retry)
	if err != nil {
		log.Fatal().Msgf("failed to fetch vuln type:%w", err)
	}
	err = utils.WriteFile(filepath.Join(Dir, VulnTypeFile), vulnType)
	if err != nil {
		log.Fatal().Str("Vuln type file", filepath.Join(Dir, VulnTypeFile)).
			Msgf("failed to write vuln type:%w", err)
	}
	bar.Finish()
}

func getVulnInfo(cnnvdId *cnnvd.CNNVD) (int, string) {
	var (
		bar        *pb.ProgressBar
		vulns      []*model.Record
		reqDetail  cnnvd.ReqVulDetail
		vulnDetail *model.VulDetail
		newCnnvdId *cnnvd.CNNVD
	)

	// 获取漏洞的总数和循环获取漏洞列表的数字
	keyword, err := cnnvdId.FormatCNNVD()
	if err != nil {
		log.Fatal().Interface("CNNVD ID", cnnvdId).Msgf("failed to format cnnvd:%v", err)
	}
	reqList := cnnvd.ReqVulList{PageSize: cnnvd.MaxPageSize, Keyword: keyword}
	loopNum, total, err := reqList.GetPageInfo(Retry)
	if err != nil {
		log.Fatal().Interface("Request", reqList).Msgf("failed to get page num:%w", err)
	}
	log.Info().Str("Keyword", keyword).Int("Total", total).Msg("Saving CNNVD data...")

	// 循环漏洞列表，获取漏洞信息
	if !verbose {
		bar = pb.StartNew(total)
	}
	latestCnnvdId := OldestCnnvdId
	latestCnnvd, _ := cnnvd.NewCNNVD(latestCnnvdId)
	for i := 1; i <= loopNum; i++ {
		reqList = cnnvd.ReqVulList{PageIndex: 1, PageSize: cnnvd.MaxPageSize, Keyword: keyword}
		vulns, err = reqList.Fetch(Retry)
		if err != nil {
			log.Fatal().Interface("request", reqList).Msgf("failed to get vuln list:%w", err)
		}
		for _, vuln := range vulns {
			//如果指定了 Complete ，还需要收集漏洞详情
			var data interface{}
			if !Complete {
				data = vuln
			} else {
				reqDetail = cnnvd.ReqVulDetail{Id: vuln.Id, VulType: vuln.VulType, CnnvdCode: vuln.CnnvdCode}
				vulnDetail, err = reqDetail.Fetch(Retry)
				if err != nil {
					log.Fatal().Interface("request", reqDetail).Msgf("failed to get vuln detail:%w", err)
				}
				data = vulnDetail
			}

			// 保存漏洞信息
			if err = save(vuln.CnnvdCode, data); err != nil {
				log.Fatal().Str("CNNVD ID", vuln.CnnvdCode).Msgf("failed to write print vuln:%w", err)
			}

			// 获取最新的cnnvd 编号
			newCnnvdId, _ = cnnvd.NewCNNVD(vuln.CnnvdCode)
			if newCnnvdId.After(latestCnnvd) {
				latestCnnvd = newCnnvdId
				latestCnnvdId = vuln.CnnvdCode
			}
			if !verbose {
				bar.Increment()
			}
		}
	}
	if !verbose {
		bar.Finish()
	}
	return total, latestCnnvdId
}

func save(vulnId string, data interface{}) error {
	// new 一个CNNVD ID
	cnnvdId, err := cnnvd.NewCNNVD(vulnId)
	if err != nil {
		return xerrors.Errorf("failed to new cnnvd:%w", err)
	}
	// 创建目录
	dir := filepath.Join(Dir, strconv.Itoa(cnnvdId.Year), strconv.Itoa(cnnvdId.Month))
	if err := utils.Mkdir(dir); err != nil {
		return xerrors.Errorf("failed to mkdir %s:%w", dir, err)
	}
	// 创建文件
	filename := filepath.Join(dir, fmt.Sprintf("%s.json", vulnId))
	err = utils.WriteFile(filename, data)
	if err != nil {
		return xerrors.Errorf("failed to write %s:%w", filename, err)
	}
	if verbose {
		log.Debug().Str("file", filename).Msg("success to save cnnvd info")
	}
	return nil
}
