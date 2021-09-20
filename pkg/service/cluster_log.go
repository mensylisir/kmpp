package service

import (
	"github.com/kmpp/pkg/constant"
	"github.com/kmpp/pkg/db"
	"github.com/kmpp/pkg/dto"
	"github.com/kmpp/pkg/model"
	"github.com/kmpp/pkg/repository"
	"time"
)

type ClusterLogService interface {
	List(clusterName string) ([]dto.ClusterLog, error)
	Save(clusterName string, clusterLog *model.ClusterLog) error
	Start(log *model.ClusterLog) error
	End(log *model.ClusterLog, success bool, message string) error
	GetRunningLogWithClusterNameAndType(clusterName string, logType string) (model.ClusterLog, error)
}

type clusterLogService struct {
	clusterLogRepo repository.ClusterLogRepository
}

func NewClusterLogService() ClusterLogService {
	return &clusterLogService{
		clusterLogRepo: repository.NewClusterLogRepository(),
	}
}

func (c *clusterLogService) List(clusterName string) ([]dto.ClusterLog, error) {
	mos, err := c.clusterLogRepo.List(clusterName)
	if err != nil {
		return nil, err
	}
	var items []dto.ClusterLog

	for _, mo := range mos {
		items = append(items, dto.ClusterLog{ClusterLog: mo})
	}
	return items, nil
}

func (c *clusterLogService) Save(clusterName string, clusterLog *model.ClusterLog) error {
	return c.clusterLogRepo.Save(clusterName, clusterLog)
}

func (c *clusterLogService) Start(log *model.ClusterLog) error {
	log.StartTime = time.Now()
	log.Status = constant.ClusterRunning
	return db.DB.Save(log).Error
}

func (c *clusterLogService) End(log *model.ClusterLog, success bool, message string) error {
	log.EndTime = time.Now()
	if success {
		log.Status = constant.ClusterLogStatusSuccess
	} else {
		log.Status = constant.ClusterLogStatusFailed
	}
	log.Message = message
	return db.DB.Save(log).Error
}

func (c *clusterLogService) GetRunningLogWithClusterNameAndType(clusterName string, logType string) (model.ClusterLog, error) {
	return c.clusterLogRepo.GetRunningLogWithClusterNameAndType(clusterName, logType)
}
