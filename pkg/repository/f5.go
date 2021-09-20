package repository

import (
	"github.com/kmpp/pkg/db"
	"github.com/kmpp/pkg/model"
)

type F5Repository interface {
	Save(f5 *model.F5Setting) error
	Get(name string) (*model.F5Setting, error)
}

func NewF5Repository() F5Repository {
	return &f5Repository{}
}

type f5Repository struct {
}

func (f f5Repository) Save(f5 *model.F5Setting) error {
	if db.DB.NewRecord(f5) {
		return db.DB.Create(&f5).Error
	} else {
		return db.DB.Save(&f5).Error
	}
}

func (f f5Repository) Get(clusterID string) (*model.F5Setting, error) {
	var f5 model.F5Setting
	if err := db.DB.Where("cluster_id = ?", clusterID).First(&f5).Error; err != nil {
		return nil, err
	}
	return &f5, nil
}
