package service

import (
	"errors"

	"github.com/kmpp/pkg/db"
	"github.com/kmpp/pkg/dto"
	"github.com/kmpp/pkg/model"
	"github.com/kmpp/pkg/repository"
	"github.com/kmpp/pkg/util/license"
)

type LicenseService interface {
	Save(content string) (*dto.License, error)
	Get() (*dto.License, error)
}

type licenseService struct {
	licenseRepo repository.LicenseRepository
}

func NewLicenseService() LicenseService {
	return &licenseService{
		licenseRepo: repository.NewLicenseRepository(),
	}
}

var (
	errFormatLicense   = errors.New("parse license error")
	errVerification    = errors.New("license is invalid")
	errLicenseNotFound = errors.New("license not found")
)

func (l *licenseService) Save(content string) (*dto.License, error) {
	resp, err := license.Parse(content)
	if err != nil {
		return nil, errFormatLicense
	}
	if resp.Status != "valid" {
		return nil, errVerification
	}
	var lcs model.License
	notFound := db.DB.First(&lcs).RecordNotFound()
	if notFound {
		lcs.Content = content
		if err := db.DB.Create(&lcs).Error; err != nil {
			return nil, err
		}
	}
	lcs.Content = content
	if err := db.DB.Save(&lcs).Error; err != nil {
		return nil, err
	}
	return resp, nil
}

func (l *licenseService) Get() (*dto.License, error) {
	var ls dto.License
	lc, err := l.licenseRepo.Get()
	if err != nil {
		ls.Status = "invalid"
		ls.Message = errLicenseNotFound.Error()
		return &ls, nil
	}
	resp, err := license.Parse(lc.Content)
	if err != nil {
		return nil, errFormatLicense
	}
	return resp, err
}
