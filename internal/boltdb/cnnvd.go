package boltdb

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/y4ney/collect-cnnvd-vuln/internal/model"
	"github.com/y4ney/collect-cnnvd-vuln/internal/utils"
	bolt "go.etcd.io/bbolt"
	"golang.org/x/xerrors"
	"os"
	"path/filepath"
)

const (
	CnnvdDB    = "cnnvd.db"
	vulnDetail = "vuln-detail"
)

func BuildCnnvd(srcDir, desDir string) error {
	// 防止数据库越来越大，先删除原来的数据库文件
	path := filepath.Join(desDir, CnnvdDB)
	if err := utils.DeleteFile(path); err != nil {
		return err
	}

	// 打开数据库
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return xerrors.Errorf("failed to open %s:%w", CnnvdDB, err)
	}
	defer db.Close()

	// 插入数据
	if err = insert(db, srcDir, desDir); err != nil {
		return xerrors.Errorf("failed to insert %s:%w ", CnnvdDB, err)
	}

	size, err := utils.SizeOfFile(path)
	if err != nil {
		return err
	}
	log.Info().Str("Database", filepath.Join(desDir, CnnvdDB)).Str("Size", size).
		Msg("success to save into DB")
	return nil
}

// Insert 更新 CNNVD
func insert(db *bolt.DB, dir string, desDir string) error {
	// 获取CNNVD数据
	items, err := utils.GetCNNVDFromFile(dir)
	if err != nil {
		return xerrors.Errorf("failed to get cnnvd from %s:%w", dir, err)
	}
	log.Info().Int("Total CNNVD", len(items)).Msg("success to get cnnvd vuln")

	// 根据trivy优化
	vulns, err := optimizeByTrivy(items)
	if err != nil {
		return xerrors.Errorf("failed to optimize by trivy:%w", err)
	}
	log.Info().Int("Total CNNVD", len(vulns)).Msg("success to optimize cnnvd vuln")

	// 将items存入数据库中
	if err = save(db, vulns); err != nil {
		return xerrors.Errorf("error in CNNVD save: %w", err)
	}

	return nil
}

// save 多个漏洞存储。将更新操作作为原子操作
func save(db *bolt.DB, items []*model.VulDetail) error {
	var aFlag int
	err := db.Batch(func(tx *bolt.Tx) error {
		for _, item := range items {
			if err := put(tx, []string{vulnDetail}, item.CveCode, item); err != nil {
				return xerrors.Errorf("failed to put vuln detail: %w", err)
			}
			aFlag++
		}
		return nil
	})
	if err != nil {
		return xerrors.Errorf("error in batch update: %w", err)
	}
	return nil
}

func optimizeByTrivy(items []*model.VulDetail) ([]*model.VulDetail, error) {
	dir := "."
	path := filepath.Join(dir, DB, DBFile)

	// 下载trivy db
	meta, err := DownloadTrivyDB(dir)
	if err != nil {
		return nil, xerrors.Errorf("failed to download trivy DB:%w", err)
	}
	log.Debug().Interface("metadata", meta).Msg("success to download trivy DB")

	// 从 trivy db 中获取 CVE编号
	CvdIds, err := GetCvdIdFromTrivyDB(path)
	if err != nil {
		return nil, xerrors.Errorf("failed to get CVD ID from trivy DB:%w", err)
	}
	log.Debug().Int("Total CVE", len(CvdIds)).Msg("success to get CVE ID from Trivy")

	// 过滤
	var vulns []*model.VulDetail
	var links []string
	for _, item := range items {
		if CvdIds[item.CveCode] {
			item.ReferUrl = utils.FormatCnnvdRef(item.ReferUrl)
			links = append(links, item.ReferUrl)
			vulns = append(vulns, item)
		}
	}
	err = utils.WriteFile("./test.json", links)
	if err != nil {
		panic(err)
	}

	// 删除trivy db
	if err = os.RemoveAll(filepath.Join(dir, DB)); err != nil {
		return nil, xerrors.Errorf("failed to remove %s:%w", path, err)
	}

	return vulns, nil
}

// bktNames 单个漏洞存储。将value存储在bktNames中，并使用key作为键值
func put(tx *bolt.Tx, bktNames []string, key string, value interface{}) error {
	// 若桶的名字为空，则异常退出
	if len(bktNames) == 0 {
		return xerrors.Errorf("empty bucket name")
	}

	// 创建第一个桶（若不存在）
	bkt, err := tx.CreateBucketIfNotExists([]byte(bktNames[0]))
	if err != nil {
		return xerrors.Errorf("failed to create '%s' bucket: %w", bktNames[0], err)
	}

	// 递归创建桶（若不存在）
	for _, bktName := range bktNames[1:] {
		bkt, err = bkt.CreateBucketIfNotExists([]byte(bktName))
		if err != nil {
			return xerrors.Errorf("failed to create a bucket: %w", err)
		}
	}
	// 需要将value序列化为json格式
	v, err := json.Marshal(value)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal JSON: %w", err)
	}
	//存储到桶中，其中键值为key
	return bkt.Put([]byte(key), v)
}
