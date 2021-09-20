package job

import (
	"math"
	"sync"
	"time"

	"github.com/kmpp/pkg/db"
	"github.com/kmpp/pkg/dto"
	"github.com/kmpp/pkg/logger"
	"github.com/kmpp/pkg/model"
	"github.com/kmpp/pkg/repository"
	"github.com/kmpp/pkg/service"
)

type ClusterBackup struct {
	cLusterBackupFileService        service.CLusterBackupFileService
	clusterBackupStrategyRepository repository.ClusterBackupStrategyRepository
}

func NewClusterBackup() *ClusterBackup {
	return &ClusterBackup{
		cLusterBackupFileService:        service.NewClusterBackupFileService(),
		clusterBackupStrategyRepository: repository.NewClusterBackupStrategyRepository(),
	}
}

func (c *ClusterBackup) Run() {
	logger.Log.Infof("---------- start backup cron job -----------")
	var wg sync.WaitGroup
	clusterBackupStrategies, _ := c.clusterBackupStrategyRepository.List()
	for _, clusterBackupStrategy := range clusterBackupStrategies {
		if clusterBackupStrategy.Status == "ENABLE" {
			var backupFiles []model.ClusterBackupFile
			db.DB.Where("cluster_id = ?", clusterBackupStrategy.ClusterID).Order("created_at ASC").Find(&backupFiles)
			if len(backupFiles) > 0 {
				lastBackupFile := backupFiles[len(backupFiles)-1]
				backupDate := lastBackupFile.CreatedAt
				now := time.Now()
				sumD := now.Sub(backupDate)
				day := int(math.Floor(sumD.Hours() / 24))
				if day < clusterBackupStrategy.Cron {
					continue
				}
			}
			var cluster model.Cluster
			db.DB.Where("id = ?", clusterBackupStrategy.ClusterID).Find(&cluster)
			if len(backupFiles) >= clusterBackupStrategy.SaveNum {
				var deleteFileNum = len(backupFiles) + 1 - clusterBackupStrategy.SaveNum
				for i := 0; i < deleteFileNum; i++ {
					logger.Log.Infof("delete backup file %s", backupFiles[i].Name)
					err := c.cLusterBackupFileService.Delete(backupFiles[i].Name)
					if err != nil {
						logger.Log.Errorf("delete cluster [%s] backup file error : %s", cluster.Name, err.Error())
					}
				}
			}
			db.DB.Where("cluster_id = ?", clusterBackupStrategy.ClusterID).Order("created_at ASC").Find(&backupFiles)
			logger.Log.Infof("length %s", len(backupFiles))
			if len(backupFiles) < clusterBackupStrategy.SaveNum {
				wg.Add(1)
				go func() {
					defer wg.Done()
					logger.Log.Infof("backup cluster [%s]", cluster.Name)
					if cluster.ID != "" {
						err := c.cLusterBackupFileService.Backup(dto.ClusterBackupFileCreate{ClusterName: cluster.Name})
						if err != nil {
							logger.Log.Errorf("backup cluster error: %s", err.Error())
						} else {
							logger.Log.Infof("backup cluster [%s] success", cluster.Name)
						}
					}
				}()
			}
		}
	}
	wg.Wait()
	logger.Log.Infof("---------- backup cron job end -----------")
}
