package cmd

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/y4ney/collect-cnnvd-vuln/internal/boltdb"
	"github.com/y4ney/collect-cnnvd-vuln/internal/utils"
)

const (
	BoltDB           = "boltdb"
	repoURL          = "http://192.168.80.36/%s/%s.git"
	defaultRepoOwner = "cnsp-source"
	defaultRepoName  = "vuln-list"
	branch           = "master"
)

var (
	Trivy    bool
	Database string
	CnnvdDir string
	desDir   string
)

// saveCmd represents the save command
var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "保存 CNNVD 漏洞信息至数据库中",
	RunE:  runSave,
}

func init() {
	saveCmd.Flags().BoolVarP(&Trivy, "trivy", "t", false, "fill cnnvd vuln by trivy")
	saveCmd.Flags().StringVar(&Database, "database", BoltDB, "save cnnvd vuln to database(boltdb)")
	saveCmd.Flags().StringVarP(&CnnvdDir, "cnnvd-dir", "c", ".", "specify the cnnvd dir")
	saveCmd.Flags().StringVarP(&desDir, "des-dir", "d", ".", "specify the des dir")
}
func runSave(_ *cobra.Command, _ []string) error {

	if err := boltdb.BuildCnnvd(CnnvdDir, desDir); err != nil {
		log.Error().Msg(err.Error())
		return err
	}

	return nil
}
func GetCnnvdVuln() {
	url := fmt.Sprintf(repoURL, defaultRepoOwner, defaultRepoName)
	git := utils.Git{URL: url, Dir: Dir}
	if err := git.Clone(); err != nil {
		log.Fatal().Str("repository", url).Str("directory", Dir).
			Msgf("failed to git clone:%v", err)
	}
}
