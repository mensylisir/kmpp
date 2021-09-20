package controller

import (
	"github.com/kmpp/pkg/constant"
	"github.com/kmpp/pkg/controller/kolog"
	"github.com/kmpp/pkg/controller/page"
	"github.com/kmpp/pkg/dto"
	"github.com/kmpp/pkg/service"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

type ClusterResourceController struct {
	Ctx                    context.Context
	ClusterResourceService service.ClusterResourceService
}

func NewClusterResourceController() *ClusterResourceController {
	return &ClusterResourceController{
		ClusterResourceService: service.NewClusterResourceService(),
	}
}

// List clusterResource
// @Tags clusterResources
// @Summary Show all clusterResources
// @Description 获取集群资源列表
// @Accept  json
// @Produce  json
// @Param project path string true "项目名称"
// @Param cluster path string true "集群名称"
// @Success 200 {object} page.Page
// @Security ApiKeyAuth
// @Router /projects/{project}/clusters/{cluster}/resources [get]
func (c ClusterResourceController) Get() (*page.Page, error) {
	pa, _ := c.Ctx.Values().GetBool("page")
	resourceType := c.Ctx.URLParam("resourceType")
	clusterName := c.Ctx.Params().GetString("cluster")
	if pa {
		num, _ := c.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := c.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return c.ClusterResourceService.Page(num, size, clusterName, resourceType)
	} else {
		var page page.Page
		resource, err := c.ClusterResourceService.List(clusterName, resourceType)
		page.Items = resource
		return &page, err
	}
}

// Create ClusterResource
// @Tags clusterResources
// @Summary Create a clusterResource
// @Description 授权资源到集群
// @Accept  json
// @Produce  json
// @Param request body dto.ClusterResourceCreate true "request"
// @Param project path string true "项目名称"
// @Param cluster path string true "集群名称"
// @Success 200 {Array} []dto.ClusterResource
// @Security ApiKeyAuth
// @Router /projects/{project}/clusters/{cluster}/resources [post]
func (c ClusterResourceController) Post() ([]dto.ClusterResource, error) {
	clusterName := c.Ctx.Params().GetString("cluster")
	var req dto.ClusterResourceCreate
	if err := c.Ctx.ReadJSON(&req); err != nil {
		return nil, err
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, err
	}

	operator := c.Ctx.Values().GetString("operator")
	resource := clusterName + "("
	for _, re := range req.Names {
		resource += (re + ",")
	}
	switch req.ResourceType {
	case "PLAN":
		go kolog.Save(operator, constant.BIND_CLUSTER_RESOURCE_PLAN, resource+")")
	case "BACKUP_ACCOUNT":
		go kolog.Save(operator, constant.BIND_CLUSTER_RESOURCE_BACKUP, resource+")")
	case "HOST":
		go kolog.Save(operator, constant.BIND_CLUSTER_RESOURCE_HOST, resource+")")
	}

	return c.ClusterResourceService.Create(clusterName, req)
}

// Delete ClusterResource
// @Tags clusterResources
// @Summary Delete clusterResource
// @Description 取消集群资源授权
// @Accept  json
// @Produce  json
// @Param resourceType query string true  "资源类型（HOST,PLAN,BACKUP_ACCOUNT）"
// @Param project path string true "项目名称"
// @Param cluster path string true "集群名称"
// @Param name path string true "资源名称"
// @Security ApiKeyAuth
// @Router /projects/{project}/clusters/{cluster}/resources/{name} [delete]
func (c ClusterResourceController) DeleteBy(name string) error {
	resourceType := c.Ctx.URLParam("resourceType")
	clusterName := c.Ctx.Params().GetString("cluster")

	operator := c.Ctx.Values().GetString("operator")
	switch resourceType {
	case "PLAN":
		go kolog.Save(operator, constant.UNBIND_CLUSTER_RESOURCE_PLAN, clusterName+"-"+name)
	case "BACKUP_ACCOUNT":
		go kolog.Save(operator, constant.UNBIND_CLUSTER_RESOURCE_BACKUP, clusterName+"-"+name)
	case "HOST":
		go kolog.Save(operator, constant.UNBIND_CLUSTER_RESOURCE_HOST, clusterName+"-"+name)
	}

	return c.ClusterResourceService.Delete(name, resourceType, clusterName)
}

func (c ClusterResourceController) GetList() (interface{}, error) {
	resourceType := c.Ctx.URLParam("resourceType")
	projectName := c.Ctx.Params().GetString("project")
	clusterName := c.Ctx.Params().GetString("cluster")
	return c.ClusterResourceService.GetResources(resourceType, projectName, clusterName)
}
