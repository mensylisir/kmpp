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

type ProjectController struct {
	Ctx            context.Context
	ProjectService service.ProjectService
}

func NewProjectController() *ProjectController {
	return &ProjectController{
		ProjectService: service.NewProjectService(),
	}
}

// List Project
// @Tags projects
// @Summary Show all projects
// @Description 获取项目列表
// @Accept  json
// @Produce  json
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /projects [get]
func (p ProjectController) Get() (*page.Page, error) {
	pa, _ := p.Ctx.Values().GetBool("page")
	sessionUser := p.Ctx.Values().Get("user")
	user, _ := sessionUser.(dto.SessionUser)
	if pa {
		num, _ := p.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := p.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return p.ProjectService.Page(num, size, user, condition.TODO())
	} else {
		var page page.Page
		items, err := p.ProjectService.List(user, condition.TODO())
		if err != nil {
			return &page, err
		}
		page.Items = items
		page.Total = len(items)
		return &page, nil
	}
}

func (p ProjectController) PostSearch() (*page.Page, error) {
	pa, _ := p.Ctx.Values().GetBool("page")
	var conditions condition.Conditions
	sessionUser := p.Ctx.Values().Get("user")
	user, _ := sessionUser.(dto.SessionUser)
	if p.Ctx.GetContentLength() > 0 {
		if err := p.Ctx.ReadJSON(&conditions); err != nil {
			return nil, err
		}
	}
	if pa {
		num, _ := p.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := p.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return p.ProjectService.Page(num, size, user, conditions)
	} else {
		var page page.Page
		items, err := p.ProjectService.List(user, condition.TODO())
		if err != nil {
			return &page, err
		}
		page.Items = items
		page.Total = len(items)
		return &page, nil
	}
}

// Get Project
// @Tags projects
// @Summary Show a project
// @Description 获取单个项目
// @Accept  json
// @Produce  json
// @Param name path string true "项目名称"
// @Success 200 {object} dto.Project
// @Security ApiKeyAuth
// @Router /projects/{name} [get]
func (p ProjectController) GetBy(name string) (*dto.Project, error) {
	return p.ProjectService.Get(name)
}

// Create Project
// @Tags projects
// @Summary Create a project
// @Description 创建项目
// @Accept  json
// @Produce  json
// @Param request body dto.ProjectCreate true "request"
// @Success 200 {object} dto.Project
// @Security ApiKeyAuth
// @Router /projects [post]
func (p ProjectController) Post() (*dto.Project, error) {
	var req dto.ProjectCreate
	err := p.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return nil, err
	}
	result, err := p.ProjectService.Create(req)
	if err != nil {
		return result, err
	}

	operator := p.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.CREATE_PROJECT, req.Name)

	return nil, err
}

// Update Project
// @Tags projects
// @Summary Update a project
// @Description 更新项目
// @Accept  json
// @Produce  json
// @Param request body dto.ProjectUpdate true "request"
// @Param name path string true "项目名称"
// @Success 200 {object} dto.Project
// @Security ApiKeyAuth
// @Router /projects/{name} [patch]
func (p ProjectController) PatchBy(name string) (*dto.Project, error) {
	var req dto.ProjectUpdate
	err := p.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return nil, err
	}
	result, err := p.ProjectService.Update(name, req)
	if err != nil {
		return nil, err
	}

	operator := p.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.UPDATE_PROJECT_INFO, name)

	return result, nil
}

// Delete Project
// @Tags projects
// @Summary Delete a project
// @Description 删除项目
// @Accept  json
// @Produce  json
// @Param name path string true "项目名称"
// @Success 200
// @Security ApiKeyAuth
// @Router /projects/{name} [delete]
func (p ProjectController) DeleteBy(name string) error {
	operator := p.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.DELETE_PROJECT, name)

	return p.ProjectService.Delete(name)
}

func (p ProjectController) PostBatch() error {
	var req dto.ProjectOp
	err := p.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}
	err = p.ProjectService.Batch(req)
	if err != nil {
		return err
	}

	operator := p.Ctx.Values().GetString("operator")
	delProjects := ""
	for _, item := range req.Items {
		delProjects += item.Name + ","
	}
	go kolog.Save(operator, constant.DELETE_PROJECT, delProjects)

	return err
}

func (p ProjectController) GetTree() ([]dto.ProjectResourceTree, error) {
	sessionUser := p.Ctx.Values().Get("user")
	user, _ := sessionUser.(dto.SessionUser)
	return p.ProjectService.GetResourceTree(user)
}
