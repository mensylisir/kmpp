package controller

import (
	"io/ioutil"

	"github.com/kmpp/pkg/constant"
	"github.com/kmpp/pkg/controller/kolog"
	"github.com/kmpp/pkg/dto"
	"github.com/kmpp/pkg/service"
	"github.com/kataras/iris/v12/context"
)

type LicenseController struct {
	Ctx            context.Context
	LicenseService service.LicenseService
}

func NewLicenseController() *LicenseController {
	return &LicenseController{
		LicenseService: service.NewLicenseService(),
	}
}

func (l *LicenseController) Get() (*dto.License, error) {
	return l.LicenseService.Get()
}
func (l *LicenseController) Post() (*dto.License, error) {
	f, _, err := l.Ctx.FormFile("file")
	if err != nil {
		return nil, err
	}
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	operator := l.Ctx.Values().GetString("operator")
	go kolog.Save(operator, constant.IMPORT_LICENCE, "-")

	return l.LicenseService.Save(string(bs))
}
