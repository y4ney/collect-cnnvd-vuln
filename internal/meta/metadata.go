package meta

import (
	"github.com/rs/zerolog/log"
	"github.com/y4ney/collect-cnnvd-vuln/internal/utils"
	"os"
	"path/filepath"
	"time"
)

const (
	OldestCnnvdId = "CNNVD-198801-001"
	MetadataFile  = "metadata.json"
)

type Data struct {
	TotalVuln   int    `json:"total_vuln"`
	LatestCnnvd string `json:"latest_cnnvd"`

	NextIncrementUpdate time.Time
	NextAllUpdate       time.Time
	UpdatedAt           time.Time

	DownloadedAt time.Time
}

func (d *Data) Init(dir string) {
	if err := utils.Mkdir(dir); err != nil {
		log.Fatal().Str("Directory", dir).Msgf("failed to mkdir:%v", err)
	}
	path := filepath.Join(dir, MetadataFile)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if _, err = os.Create(path); err != nil {
			log.Fatal().Str("Path", path).Msgf("failed to create file:%v", err)
		}
	}
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Fatal().Str("Path", path).Msgf("failed to create file:%v", err)
	}
	if fileInfo.Size() == 0 {
		d.Write(dir, 0, OldestCnnvdId)
	}
}

func (d *Data) Write(dir string, total int, latestCnnvd string) {
	var now = time.Now()
	d.NextIncrementUpdate = now.AddDate(0, 0, 1)
	d.NextAllUpdate = now.AddDate(0, 1, 0)
	d.UpdatedAt = now
	d.TotalVuln = total
	d.LatestCnnvd = latestCnnvd

	path := filepath.Join(dir, MetadataFile)
	if err := utils.WriteFile(path, d); err != nil {
		log.Fatal().Str("File", path).Interface("Metadata", d).Msg("failed to write metadata")
	}
}

func (d *Data) Read(dir string) {
	path := filepath.Join(dir, MetadataFile)
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Fatal().Str("Filepath", path).Msgf("failed to get the file info:%v", err)
	}
	if fileInfo.Size() == 0 {
		log.Fatal().Str("Filepath", path).Msg("metadata file ie empty,please init firstly")
	}

	if err := utils.ReadFile(path, d); err != nil {
		log.Fatal().Str("Filepath", path).Msg("failed to read metadata")
	}
}
