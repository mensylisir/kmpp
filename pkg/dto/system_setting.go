package dto

import "github.com/kmpp/pkg/model"

type SystemSetting struct {
	model.SystemSetting
}

type SystemSettingCreate struct {
	Vars map[string]string `json:"vars" validate:"required"`
	Tab  string            `json:"tab" validate:"required"`
}

type SystemSettingUpdate struct {
	Vars map[string]string `json:"vars" validate:"required"`
	Tab  string            `json:"tab" validate:"required"`
}

type SystemSettingResult struct {
	Vars map[string]string `json:"vars" validate:"required"`
	Tab  string            `json:"tab" validate:"required"`
}
