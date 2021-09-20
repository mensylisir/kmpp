package service

import (
	"github.com/kmpp/pkg/controller/condition"
	"github.com/kmpp/pkg/controller/page"
	"github.com/kmpp/pkg/db"
	"github.com/kmpp/pkg/dto"
	"github.com/kmpp/pkg/model"
	dbUtil "github.com/kmpp/pkg/util/db"
)

type SystemLogService interface {
	Create(creation dto.SystemLogCreate) error
	Page(num, size int, conditions condition.Conditions) (*page.Page, error)
}

type systemLogService struct{}

func NewSystemLogService() SystemLogService {
	return &systemLogService{}
}

func (s systemLogService) Create(creation dto.SystemLogCreate) error {
	log := model.SystemLog{
		Name:          creation.Name,
		Operation:     creation.Operation,
		OperationInfo: creation.OperationInfo,
	}

	if db.DB.NewRecord(log) {
		return db.DB.Create(&log).Error
	} else {
		return db.DB.Save(&log).Error
	}
}

func (u systemLogService) Page(num, size int, conditions condition.Conditions) (*page.Page, error) {
	var (
		p         page.Page
		logOfDTOs []dto.SystemLog
		mos       []model.SystemLog
	)
	d := db.DB.Model(model.SystemLog{})
	if err := dbUtil.WithConditions(&d, model.SystemLog{}, conditions); err != nil {
		return nil, err
	}
	if err := d.
		Count(&p.Total).
		Order("created_at DESC").
		Offset((num - 1) * size).
		Limit(size).
		Find(&mos).Error; err != nil {
		return nil, err
	}

	for _, mo := range mos {
		logOfDTOs = append(logOfDTOs, dto.SystemLog{SystemLog: mo})
	}
	p.Items = logOfDTOs
	return &p, nil
}
