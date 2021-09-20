package repository

import (
	"github.com/kmpp/pkg/db"
	"github.com/kmpp/pkg/model"
)

type MessageRepository interface {
	Save(message *model.Message) error
}

type messageRepository struct {
}

func NewMessageRepository() MessageRepository {
	return &messageRepository{}
}

func (m messageRepository) Save(message *model.Message) error {
	if db.DB.NewRecord(message) {
		return db.DB.Create(&message).Error
	} else {
		return db.DB.Save(&message).Error
	}
}
