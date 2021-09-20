package controller

import (
	"github.com/kmpp/pkg/constant"
	"github.com/kmpp/pkg/controller/condition"
	"github.com/kmpp/pkg/controller/kolog"
	"github.com/kmpp/pkg/controller/page"
	"github.com/kmpp/pkg/dto"
	"github.com/kmpp/pkg/service"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

type RegionController struct {
	Ctx           context.Context
	RegionService service.RegionService
}

func NewRegionController() *RegionController {
	return &RegionController{
		RegionService: service.NewRegionService(),
	}
}

// List Region
// @Tags regions
// @Summary Show all regions
// @Description 获取区域列表
// @Accept  json
// @Produce  json
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /regions [get]
func (r RegionController) Get() (*page.Page, error) {

	p, _ := r.Ctx.Values().GetBool("page")
	var conditions condition.Conditions
	if r.Ctx.GetContentLength() > 0 {
		if err := r.Ctx.ReadJSON(&conditions); err != nil {
			return nil, err
		}
	}
	if p {
		num, _ := r.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := r.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return r.RegionService.Page(num, size, conditions)
	} else {
		var p page.Page
		items, err := r.RegionService.List(condition.TODO())
		if err != nil {
			return nil, err
		}
		p.Items = items
		p.Total = len(items)
		return &p, nil
	}
}

// Search Region
// @Tags regions
// @Summary Search regions
// @Description 过滤部署计划
// @Accept  json
// @Produce  json
// @Param conditions body condition.Conditions true "conditions"
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /regions/search [post]
func (r RegionController) PostSearch() (*page.Page, error) {

	p, _ := r.Ctx.Values().GetBool("page")
	var conditions condition.Conditions
	if r.Ctx.GetContentLength() > 0 {
		if err := r.Ctx.ReadJSON(&conditions); err != nil {
			return nil, err
		}
	}
	if p {
		num, _ := r.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := r.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return r.RegionService.Page(num, size, conditions)
	} else {
		var p page.Page
		items, err := r.RegionService.List(condition.TODO())
		if err != nil {
			return nil, err
		}
		p.Items = items
		p.Total = len(items)
		return &p, nil
	}
}

// Get Region
// @Tags regions
// @Summary Show a region
// @Description 获取单个区域
// @Accept  json
// @Produce  json
// @Param name path string true "区域名称"
// @Success 200 {object} dto.Region
// @Security ApiKeyAuth
// @Router /regions/{name} [get]
func (r RegionController) GetBy(name string) (dto.Region, error) {
	return r.RegionService.Get(name)
}

// Create Region
// @Tags regions
// @Summary Create a region
// @Description 创建区域
// @Accept  json
// @Produce  json
// @Param request body dto.RegionCreate true "request"
// @Success 200 {object} dto.Region
// @Security ApiKeyAuth
// @Router /regions [post]
func (r RegionController) Post() (*dto.Region, error) {
	var req dto.RegionCreate
	err := r.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return nil, err
	}

	operator := r.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.CREATE_REGION, req.Name)

	return r.RegionService.Create(req)
}

// Update Region
// @Tags regions
// @Summary Update a region
// @Description 更新区域
// @Accept  json
// @Produce  json
// @Param request body dto.RegionUpdate true "request"
// @Param name path string true "区域名称"
// @Success 200 {object} dto.Region
// @Security ApiKeyAuth
// @Router /regions/{name} [patch]
func (r RegionController) PatchBy(name string) (*dto.Region, error) {
	var req dto.RegionUpdate
	err := r.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return nil, err
	}
	operator := r.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.CREATE_REGION, name)
	return r.RegionService.Update(name, req)
}

// Delete Region
// @Tags regions
// @Summary Delete a region
// @Description 删除区域
// @Accept  json
// @Produce  json
// @Param name path string true "区域名称"
// @Security ApiKeyAuth
// @Router /regions/{name} [delete]
func (r RegionController) DeleteBy(name string) error {
	operator := r.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.DELETE_REGION, name)

	return r.RegionService.Delete(name)
}

func (r RegionController) PostBatch() error {
	var req dto.RegionOp
	err := r.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}
	err = r.RegionService.Batch(req)
	if err != nil {
		return err
	}

	operator := r.Ctx.Values().GetString("operator")
	delRegions := ""
	for _, item := range req.Items {
		delRegions += item.Name + ","
	}
	go kolog.Save(operator, constant.DELETE_REGION, delRegions)

	return err
}

// Get Datacenter List
// @Tags regions
// @Summary Get datacenter list
// @Description 获取数据中心
// @Accept  json
// @Produce  json
// @Param request body dto.RegionDatacenterRequest true "request"
// @Success 200 {object} dto.CloudRegionResponse
// @Security ApiKeyAuth
// @Router /regions/datacenter [post]
func (r RegionController) PostDatacenter() (*dto.CloudRegionResponse, error) {
	var req dto.RegionDatacenterRequest
	err := r.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	data, err := r.RegionService.ListDatacenter(req)
	if err != nil {
		return nil, err
	}
	return &dto.CloudRegionResponse{Result: data}, err
}
