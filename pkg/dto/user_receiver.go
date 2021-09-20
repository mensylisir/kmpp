package dto

import "github.com/kmpp/pkg/model"

type UserReceiver struct {
	model.UserReceiver
}

type UserReceiverDTO struct {
	ID     string            `json:"id"`
	UserID string            `json:"userId"`
	Vars   map[string]string `json:"vars"`
}
