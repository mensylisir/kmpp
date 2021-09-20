package controller

import (
	"github.com/kmpp/pkg/constant"
	"github.com/kmpp/pkg/controller/kolog"
	"github.com/kmpp/pkg/dto"
	"github.com/kmpp/pkg/service"
	"github.com/kmpp/pkg/util/validator_error"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

type ClusterBackupStrategyController struct {
	Ctx                          context.Context
	CLusterBackupStrategyService service.CLusterBackupStrategyService
}

func NewClusterBackupStrategyController() *ClusterBackupStrategyController {
	return &ClusterBackupStrategyController{
		CLusterBackupStrategyService: service.NewCLusterBackupStrategyService(),
	}
}

// Get Cluster Backup Strategy By ClusterName
// @Tags backupStrategy
// @Summary Get Cluster Backup Strategy
// @Description Get Cluster Backup Strategy
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.ClusterBackupStrategy
// @Security ApiKeyAuth
// @Router /cluster/backup/strategy/{clusterName}/ [get]
func (c ClusterBackupStrategyController) GetStrategyBy(clusterName string) (*dto.ClusterBackupStrategy, error) {
	cb, err := c.CLusterBackupStrategyService.Get(clusterName)
	if err != nil {
		return nil, err
	}
	return cb, nil
}

// Create/Update Cluster Backup Strategy
// @Tags backupStrategy
// @Summary Create a Backup Strategy
// @Description create a Backup Strategy
// @Accept  json
// @Produce  json
// @Param request body dto.ClusterBackupStrategyRequest true "request"
// @Success 200 {object} dto.ClusterBackupStrategy
// @Security ApiKeyAuth
// @Router /cluster/backup/strategy/ [post]
func (c ClusterBackupStrategyController) PostStrategy() (*dto.ClusterBackupStrategy, error) {
	var req dto.ClusterBackupStrategyRequest
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	validator_error.RegisterTagNameFunc(c.Ctx, validate)
	err = validate.Struct(req)
	if err != nil {
		return nil, validator_error.Tr(c.Ctx, validate, err)
	}
	cb, err := c.CLusterBackupStrategyService.Save(req)
	if err != nil {
		return nil, err
	}
	operator := c.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.CREATE_CLUSTER_BACKUP_STRATEGY, req.ClusterName)
	return cb, nil
}
