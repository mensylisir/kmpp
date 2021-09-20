package controller

import (
	"github.com/kmpp/pkg/constant"
	"github.com/kmpp/pkg/controller/condition"
	"github.com/kmpp/pkg/controller/kolog"
	"github.com/kmpp/pkg/controller/page"
	"github.com/kmpp/pkg/dto"
	"github.com/kmpp/pkg/service"
	"github.com/kmpp/pkg/util/validator_error"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

type VmConfigController struct {
	Ctx             context.Context
	VmConfigService service.VmConfigService
}

func NewVmConfigController() *VmConfigController {
	return &VmConfigController{
		VmConfigService: service.NewVmConfigService(),
	}
}

// List VmConfigs
// @Tags vmConfigs
// @Summary Show all vmConfigs
// @Description 获取虚拟机配置列表
// @Accept  json
// @Produce  json
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /vmconfigs [get]
func (v VmConfigController) Get() (*page.Page, error) {
	p, _ := v.Ctx.Values().GetBool("page")
	var conditions condition.Conditions
	if v.Ctx.GetContentLength() > 0 {
		if err := v.Ctx.ReadJSON(&conditions); err != nil {
			return nil, err
		}
	}
	if p {
		num, _ := v.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := v.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return v.VmConfigService.Page(num, size, condition.TODO())
	} else {
		var p page.Page
		items, err := v.VmConfigService.List(condition.TODO())
		if err != nil {
			return nil, err
		}
		p.Items = items
		p.Total = len(items)
		return &p, nil
	}
}

// Search VmConfigs
// @Tags vmConfigs
// @Summary Search vmConfigs
// @Description 过滤虚拟机配置
// @Accept  json
// @Produce  json
// @Param conditions body condition.Conditions true "conditions"
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /vmconfigs/search [post]
func (v VmConfigController) PostSearch() (*page.Page, error) {
	p, _ := v.Ctx.Values().GetBool("page")
	var conditions condition.Conditions
	if v.Ctx.GetContentLength() > 0 {
		if err := v.Ctx.ReadJSON(&conditions); err != nil {
			return nil, err
		}
	}
	if p {
		num, _ := v.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := v.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return v.VmConfigService.Page(num, size, conditions)
	} else {
		var p page.Page
		items, err := v.VmConfigService.List(conditions)
		if err != nil {
			return nil, err
		}
		p.Items = items
		p.Total = len(items)
		return &p, nil
	}
}

// Get VmConfig
// @Tags vmConfigs
// @Summary Get a vmConfig
// @Description 获取单个虚拟机配置
// @Accept  json
// @Produce  json
// @Param name path string true "虚拟机配置名称"
// @Success 200 {object} dto.VmConfig
// @Security ApiKeyAuth
// @Router /vmconfigs/{name} [get]
func (v VmConfigController) GetBy(name string) (*dto.VmConfig, error) {
	return v.VmConfigService.Get(name)
}

// Create VmConfig
// @Tags vmConfigs
// @Summary Create a vmConfig
// @Description 创建虚拟机配置
// @Accept  json
// @Produce  json
// @Param request body dto.VmConfigCreate true "request"
// @Success 200 {object} dto.VmConfig
// @Security ApiKeyAuth
// @Router /vmconfigs [post]
func (v VmConfigController) Post() (*dto.VmConfig, error) {
	var req dto.VmConfigCreate
	err := v.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	validator_error.RegisterTagNameFunc(v.Ctx, validate)
	err = validate.Struct(req)
	if err != nil {
		return nil, validator_error.Tr(v.Ctx, validate, err)
	}

	operator := v.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.CREATE_VM_CONFIG, req.Name)

	return v.VmConfigService.Create(req)
}

// Update VmConfig
// @Tags vmConfigs
// @Summary Update a vmConfig
// @Description 更新虚拟机配置
// @Accept  json
// @Produce  json
// @Param request body dto.VmConfigUpdate true "request"
// @Param name path string true "虚拟机配置名称"
// @Success 200 {object} dto.VmConfig
// @Security ApiKeyAuth
// @Router /vmconfigs/{name} [patch]
func (v VmConfigController) PatchBy(name string) (*dto.VmConfig, error) {
	var req dto.VmConfigUpdate
	err := v.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	validator_error.RegisterTagNameFunc(v.Ctx, validate)
	err = validate.Struct(req)
	if err != nil {
		return nil, validator_error.Tr(v.Ctx, validate, err)
	}
	result, err := v.VmConfigService.Update(name, req)
	if err != nil {
		return nil, err
	}

	operator := v.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.UPDATE_VM_CONFIG, name)

	return result, nil
}

// Delete VmConfig
// @Tags vmConfigs
// @Summary Delete a vmConfig
// @Description 删除虚拟机配置
// @Accept  json
// @Produce  json
// @Param name path string true "虚拟机配置名称"
// @Security ApiKeyAuth
// @Router /vmconfigs/{name} [delete]
func (v VmConfigController) DeleteBy(name string) error {
	operator := v.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.DELETE_VM_CONFIG, name)
	return v.VmConfigService.Delete(name)
}

func (v VmConfigController) PostBatch() error {
	var req dto.VmConfigOp
	err := v.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}
	err = v.VmConfigService.Batch(req)
	if err != nil {
		return err
	}

	operator := v.Ctx.Values().GetString("operator")
	delConfs := ""
	for _, item := range req.Items {
		delConfs += item.Name + ","
	}
	go kolog.Save(operator, constant.DELETE_VM_CONFIG, delConfs)

	return err
}
