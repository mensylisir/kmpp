package dto

import "github.com/kmpp/pkg/model"

type UserNotificationConfig struct {
	model.UserNotificationConfig
}

type UserNotificationConfigDTO struct {
	ID     string            `json:"id"`
	UserID string            `json:"userId"`
	Vars   map[string]string `json:"vars"`
	Type   string            `json:"type"`
}
