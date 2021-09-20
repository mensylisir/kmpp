package repository

import (
	"github.com/kmpp/pkg/db"
	"github.com/kmpp/pkg/model"
)

type LicenseRepository interface {
	Save(content string) error
	Get() (model.License, error)
}

func NewLicenseRepository() LicenseRepository {
	return &licenseRepository{}
}

type licenseRepository struct {
}

func (l licenseRepository) Save(content string) error {
	var license model.License
	if notFound := db.DB.First(&license).RecordNotFound(); notFound {
		license.Content = content
		return db.DB.Create(&license).Error
	} else {
		license.Content = content
		return db.DB.Save(&license).Error
	}
}

func (l licenseRepository) Get() (model.License, error) {
	var license model.License
	err := db.DB.First(&license).Error
	if err != nil {
		return license, err
	}
	return license, nil
}
