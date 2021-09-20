package controller

import (
	"github.com/kataras/iris/v12/context"
	"github.com/kmpp/pkg/dto"
	"github.com/kmpp/pkg/service"
)

type GrpcproxyController struct {
	Ctx context.Context
	GrpcproxyService service.GrpcproxyService
}

func NewGrpcproxyController() *GrpcproxyController {
	return &GrpcproxyController{
		GrpcproxyService:          service.NewGrpcproxyService(),
	}
}

func (c *GrpcproxyController) Get() []string {
	return c.GrpcproxyService.GetActiveConns()
}

func (c *GrpcproxyController) PostService() ([]string, error) {
	var req dto.ActiveList
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	return c.GrpcproxyService.GetLists(req.Host, "", req.UseTLS, req.Restart)
}

func (c *GrpcproxyController) PostFunction() ([]string, error) {
	var req dto.ActiveList
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	return c.GrpcproxyService.GetLists(req.Host, req.Service, req.UseTLS, req.Restart)
}

func (c *GrpcproxyController) PostFunctionParam() (interface{}, error) {
	var req dto.ActiveList
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	return c.GrpcproxyService.DescribeFunction(req.Host, req.FunName, req.UseTLS)
}

func (c *GrpcproxyController) PostFunctionInvoke() (interface{}, error) {
	var req dto.ActiveList
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	metadataHeader := c.Ctx.GetHeader("Metadata")
	return c.GrpcproxyService.InvokeFunction(req.Host, req.FunName, req.UseTLS, metadataHeader, req.Body)
}

