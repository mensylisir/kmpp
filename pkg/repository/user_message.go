package repository

import (
	"github.com/kmpp/pkg/constant"
	"github.com/kmpp/pkg/db"
	"github.com/kmpp/pkg/model"
)

type UserMessageRepository interface {
	Page(num int, size int, userId string) (int, []model.UserMessage, error)
	Batch(operation string, items []model.UserMessage) error
	Save(message *model.UserMessage) error
	ListUnreadMsg(userId string) ([]model.UserMessage, error)
}

func NewUserMessageRepository() UserMessageRepository {
	return &userMessageRepository{}
}

type userMessageRepository struct {
}

func (u userMessageRepository) Page(num int, size int, userId string) (int, []model.UserMessage, error) {
	var total int
	var userMessages []model.UserMessage
	if err := db.DB.Model(&model.UserMessage{}).
		Where("user_id = ?", userId).
		Count(&total).
		Order("created_at desc").
		Offset((num - 1) * size).
		Limit(size).
		Preload("Message").
		Preload("Message.Cluster").
		Find(&userMessages).
		Error; err != nil {
		return total, nil, err
	}
	return total, userMessages, nil
}

func (u userMessageRepository) ListUnreadMsg(userId string) ([]model.UserMessage, error) {
	var userMessages []model.UserMessage
	if err := db.DB.
		Where("user_id = ? AND read_status = ?", userId, constant.UnRead).
		Preload("Message").
		Find(&userMessages).
		Error; err != nil {
		return nil, err
	}
	return userMessages, nil
}

func (u userMessageRepository) Save(message *model.UserMessage) error {
	if db.DB.NewRecord(message) {
		return db.DB.Create(&message).Error
	} else {
		return db.DB.Save(&message).Error
	}
}

func (u userMessageRepository) Batch(operation string, items []model.UserMessage) error {

	tx := db.DB.Begin()
	switch operation {
	case constant.BatchOperationDelete:
		for i := range items {
			var userMessage model.UserMessage
			if err := tx.Where("id = ?", items[i].ID).First(&userMessage).Error; err != nil {
				tx.Rollback()
				return err
			}
			if err := tx.Delete(&userMessage).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	case constant.BatchOperationUpdate:
		for i := range items {
			item := items[i]
			item.ReadStatus = constant.Read
			if err := tx.Save(&item).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	default:
		return constant.NotSupportedBatchOperation
	}
	tx.Commit()
	return nil
}
