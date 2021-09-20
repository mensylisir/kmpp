package controller

import (
	"github.com/kmpp/pkg/constant"
	"github.com/kmpp/pkg/controller/kolog"
	"github.com/kmpp/pkg/dto"
	"github.com/kmpp/pkg/model"
	"github.com/kmpp/pkg/service"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

type ManifestController struct {
	Ctx             context.Context
	ManifestService service.ClusterManifestService
}

func NewManifestController() *ManifestController {
	return &ManifestController{
		ManifestService: service.NewClusterManifestService(),
	}
}

func (m *ManifestController) Get() ([]dto.ClusterManifest, error) {
	return m.ManifestService.List()
}

func (m *ManifestController) GetActive() ([]dto.ClusterManifest, error) {
	return m.ManifestService.ListActive()
}

// List Manifest
// @Tags manifest
// @Summary Show all manifest
// @Description 获取Kubernetes版本列表
// @Accept  json
// @Produce  json
// @Success 200 {object} []dto.ClusterManifestGroup
// @Security ApiKeyAuth
// @Router /manifest [get]
func (m *ManifestController) GetGroup() ([]dto.ClusterManifestGroup, error) {
	return m.ManifestService.ListByLargeVersion()
}

// Update Manifest
// @Tags manifest
// @Summary Update a manifest
// @Description 更新 Kubernetes 版本状态
// @Accept  json
// @Produce  json
// @Param request body dto.ClusterManifestUpdate true "request"
// @Param name path string true "Kubernetes 版本"
// @Success 200 {object} dto.ClusterManifestUpdate
// @Security ApiKeyAuth
// @Router /manifest/{name} [patch]
func (m ManifestController) PatchBy(name string) (model.ClusterManifest, error) {
	var req dto.ClusterManifestUpdate
	err := m.Ctx.ReadJSON(&req)

	if err != nil {
		return model.ClusterManifest{}, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return model.ClusterManifest{}, err
	}

	operator := m.Ctx.Values().GetString("operator")
	if req.IsActive {
		go kolog.Save(operator, constant.ENABLE_VERSION, req.Name)
	} else {
		go kolog.Save(operator, constant.DISABLE_VERSION, req.Name)
	}

	return m.ManifestService.Update(req)
}
