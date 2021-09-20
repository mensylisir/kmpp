package model

import (
	"errors"

	"github.com/kmpp/pkg/db"
	"github.com/kmpp/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

var (
	DeleteBackupAccountFailedByProject = "DELETE_BACKUP_ACCOUNT_FAILED_BY_PROJECT"
)

type BackupAccount struct {
	common.BaseModel
	ID         string `json:"id" gorm:"type:varchar(64)"`
	Name       string `json:"name" gorm:"type:varchar(256)"`
	Bucket     string `json:"bucket" gorm:"type:varchar(256)"`
	Credential string `json:"credential" gorm:"type:text(65535)"`
	Type       string `json:"type" gorm:"type:varchar(64)"`
	Status     string `json:"status" gorm:"type:varchar(64)"`
}

func (b *BackupAccount) BeforeCreate() (err error) {
	b.ID = uuid.NewV4().String()
	return err
}

func (b *BackupAccount) BeforeDelete() (err error) {
	var backupAccounts []ProjectResource
	err = db.DB.Where(ProjectResource{ResourceID: b.ID}).Find(&backupAccounts).Error
	if err != nil {
		return err
	}
	if len(backupAccounts) > 0 {
		return errors.New(DeleteBackupAccountFailedByProject)
	}
	return err
}
