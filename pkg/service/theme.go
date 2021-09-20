package service

import (
	"github.com/kmpp/pkg/db"
	"github.com/kmpp/pkg/dto"
	"github.com/kmpp/pkg/model"
)

type ThemeService interface{
	GetConsumerTheme() (*dto.Theme, error)
}

func NewThemeService() ThemeService {
	return &themeService{}
}

type themeService struct {
}

func (*themeService) GetConsumerTheme() (*dto.Theme, error) {
	var theme model.Theme
	if err := db.DB.First(&theme).Error; err != nil {
		return nil, err
	}
	return &dto.Theme{Theme: theme}, nil
}
