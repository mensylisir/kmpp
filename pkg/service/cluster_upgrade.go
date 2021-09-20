package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"github.com/kmpp/pkg/constant"
	"github.com/kmpp/pkg/db"
	"github.com/kmpp/pkg/dto"
	"github.com/kmpp/pkg/logger"
	"github.com/kmpp/pkg/model"
	"github.com/kmpp/pkg/repository"
	"github.com/kmpp/pkg/service/cluster/adm"
	"github.com/kmpp/pkg/util/ansible"
)

type ClusterUpgradeService interface {
	Upgrade(upgrade dto.ClusterUpgrade) error
}

func NewClusterUpgradeService() ClusterUpgradeService {
	return &clusterUpgradeService{
		clusterService:    NewClusterService(),
		clusterStatusRepo: repository.NewClusterStatusRepository(),
		messageService:    NewMessageService(),
	}
}

type clusterUpgradeService struct {
	clusterService    ClusterService
	clusterStatusRepo repository.ClusterStatusRepository
	messageService    MessageService
}

func (c *clusterUpgradeService) Upgrade(upgrade dto.ClusterUpgrade) error {
	loginfo, _ := json.Marshal(upgrade)
	logger.Log.WithFields(logrus.Fields{"cluster_upgrade_info": string(loginfo)}).Debugf("start to upgrade the cluster %s", upgrade.ClusterName)

	cluster, err := c.clusterService.Get(upgrade.ClusterName)
	if err != nil {
		return fmt.Errorf("can not get cluster %s error %s", upgrade.ClusterName, err.Error())
	}
	if !(cluster.Source == constant.ClusterSourceLocal) {
		return errors.New("CLUSTER_IS_NOT_LOCAL")
	}
	if cluster.Status != constant.StatusRunning && cluster.Status != constant.StatusFailed {
		return fmt.Errorf("cluster status error %s", cluster.Status)
	}

	tx := db.DB.Begin()
	//从错误后继续
	if cluster.Cluster.Status.Phase == constant.StatusFailed && cluster.Cluster.Status.PrePhase == constant.StatusUpgrading {
		if err := tx.Model(&model.ClusterStatusCondition{}).
			Where("cluster_status_id = ? AND status = ?", cluster.StatusID, constant.ConditionFalse).
			Updates(map[string]interface{}{
				"Status":  constant.ConditionUnknown,
				"Message": "",
			}).Error; err != nil {
			return fmt.Errorf("reset status error %s", err.Error())
		}
	} else {
		if err := tx.Delete(&model.ClusterStatusCondition{}, "cluster_status_id = ?", cluster.StatusID).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("reset contidion err %s", err.Error())
		}
	}
	// 修改状态
	cluster.Cluster.Status.PrePhase = cluster.Status
	cluster.Cluster.Status.Phase = constant.StatusUpgrading
	if err := tx.Save(&cluster.Cluster.Status).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("change status err %s", err.Error())
	}
	// 创建日志
	logId, writer, err := ansible.CreateAnsibleLogWriter(cluster.Name)
	if err != nil {
		_ = c.messageService.SendMessage(constant.System, false, GetContent(constant.ClusterUpgrade, false, err.Error()), cluster.Cluster.Name, constant.ClusterUpgrade)
		return fmt.Errorf("create log error %s", err.Error())
	}
	cluster.LogId = logId
	if err := tx.Save(&cluster.Cluster).Error; err != nil {
		tx.Rollback()
		_ = c.messageService.SendMessage(constant.System, false, GetContent(constant.ClusterUpgrade, false, err.Error()), cluster.Cluster.Name, constant.ClusterUpgrade)
		return fmt.Errorf("save cluster error %s", err.Error())
	}
	cluster.Spec.UpgradeVersion = upgrade.Version
	if err := tx.Save(&cluster.Spec).Error; err != nil {
		tx.Rollback()
		_ = c.messageService.SendMessage(constant.System, false, GetContent(constant.ClusterUpgrade, false, err.Error()), cluster.Cluster.Name, constant.ClusterUpgrade)
		return fmt.Errorf("save cluster spec error %s", err.Error())
	}
	// 更新工具版本状态
	if err := c.updateToolVersion(tx, upgrade.Version, cluster.ID); err != nil {
		tx.Rollback()
		_ = c.messageService.SendMessage(constant.System, false, GetContent(constant.ClusterUpgrade, false, err.Error()), cluster.Cluster.Name, constant.ClusterUpgrade)
		return err
	}

	tx.Commit()

	logger.Log.Infof("update db data of cluster %s successful, now start to upgrade cluster", cluster.Name)
	go c.do(&cluster.Cluster, writer)
	return nil
}

func (c *clusterUpgradeService) do(cluster *model.Cluster, writer io.Writer) {
	status, err := c.clusterService.GetStatus(cluster.Name)
	if err != nil {
		logger.Log.Errorf("can not get cluster %s status, error: %s", cluster.Name, err.Error())
	}
	cluster.Status = status.ClusterStatus
	ctx, cancel := context.WithCancel(context.Background())
	admCluster := adm.NewCluster(*cluster, writer)
	statusChan := make(chan adm.Cluster)
	go c.doUpgrade(ctx, *admCluster, statusChan)
	for {
		cluster := <-statusChan
		// 保存进度
		_ = c.clusterStatusRepo.Save(&cluster.Status)
		switch cluster.Status.Phase {
		case constant.StatusRunning:
			_ = c.messageService.SendMessage(constant.System, true, GetContent(constant.ClusterUpgrade, true, ""), cluster.Name, constant.ClusterUpgrade)
			cluster.Spec.Version = cluster.Spec.UpgradeVersion
			db.DB.Save(&cluster.Spec)
			cancel()
			return
		case constant.StatusFailed:
			_ = c.messageService.SendMessage(constant.System, false, GetContent(constant.ClusterUpgrade, false, cluster.Status.Message), cluster.Name, constant.ClusterUpgrade)
			cancel()
			return
		}
	}
}
func (c clusterUpgradeService) doUpgrade(ctx context.Context, cluster adm.Cluster, statusChan chan adm.Cluster) {
	ad := adm.NewClusterAdm()
	for {
		resp, err := ad.OnUpgrade(cluster)
		if err != nil {
			cluster.Status.Message = err.Error()
		}
		cluster.Status = resp.Status
		select {
		case <-ctx.Done():
			return
		case statusChan <- cluster:
		}
		time.Sleep(5 * time.Second)
	}
}

func (c clusterUpgradeService) updateToolVersion(tx *gorm.DB, version, clusterID string) error {
	var (
		tools    []model.ClusterTool
		manifest model.ClusterManifest
		toolVars []model.VersionHelp
	)
	if err := tx.Where("name = ?", version).First(&manifest).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("get manifest error %s", err.Error())
	}
	if err := tx.Where("cluster_id = ?", clusterID).Find(&tools).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("get tools error %s", err.Error())
	}
	if err := json.Unmarshal([]byte(manifest.ToolVars), &toolVars); err != nil {
		return fmt.Errorf("unmarshal manifest.toolvar error %s", err.Error())
	}
	for _, tool := range tools {
		for _, item := range toolVars {
			if tool.Name == item.Name {
				if tool.Version != item.Version {
					if tool.Status == constant.ClusterWaiting {
						tool.Version = item.Version
					} else {
						tool.HigherVersion = item.Version
					}
					if err := tx.Save(&tool).Error; err != nil {
						return fmt.Errorf("update tool version error %s", err.Error())
					}
				}
				break
			}
		}
	}
	return nil
}
