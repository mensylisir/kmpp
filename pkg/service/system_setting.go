package service

import (
	"fmt"

	"github.com/kmpp/pkg/controller/condition"
	dbUtil "github.com/kmpp/pkg/util/db"

	"github.com/kmpp/pkg/controller/page"
	"github.com/kmpp/pkg/model/common"

	"github.com/kmpp/pkg/constant"
	"github.com/kmpp/pkg/db"
	"github.com/kmpp/pkg/dto"
	"github.com/kmpp/pkg/model"
	"github.com/kmpp/pkg/repository"
	"github.com/kmpp/pkg/util/message"
	"github.com/kmpp/pkg/util/message/client"
	"github.com/jinzhu/gorm"
)

type SystemSettingService interface {
	Get(name string) (dto.SystemSetting, error)
	GetLocalIPs() ([]model.SystemRegistry, error)
	List() (dto.SystemSettingResult, error)
	Create(creation dto.SystemSettingCreate) ([]dto.SystemSetting, error)
	ListByTab(tabName string) (dto.SystemSettingResult, error)
	CheckSettingByType(tabName string, creation dto.SystemSettingCreate) error
	ListRegistry(conditions condition.Conditions) ([]dto.SystemRegistry, error)
	PageRegistry(num, size int, conditions condition.Conditions) (*page.Page, error)
	GetRegistryByID(id string) (dto.SystemRegistry, error)
	GetRegistryByArch(arch string) (dto.SystemRegistry, error)
	CreateRegistry(creation dto.SystemRegistryCreate) (*dto.SystemRegistry, error)
	UpdateRegistry(arch string, creation dto.SystemRegistryUpdate) (*dto.SystemRegistry, error)
	BatchRegistry(op dto.SystemRegistryBatchOp) error
	DeleteRegistry(id string) error
}

type systemSettingService struct {
	systemSettingRepo  repository.SystemSettingRepository
	systemRegistryRepo repository.SystemRegistryRepository
	userRepo           repository.UserRepository
}

func NewSystemSettingService() SystemSettingService {
	return &systemSettingService{
		systemSettingRepo:  repository.NewSystemSettingRepository(),
		systemRegistryRepo: repository.NewSystemRegistryRepository(),
		userRepo:           repository.NewUserRepository(),
	}
}

func (s systemSettingService) Get(key string) (dto.SystemSetting, error) {
	var systemSettingDTO dto.SystemSetting
	mo, err := s.systemSettingRepo.Get(key)
	if err != nil {
		return systemSettingDTO, err
	}
	systemSettingDTO.SystemSetting = mo
	return systemSettingDTO, err
}

func (s systemSettingService) List() (dto.SystemSettingResult, error) {
	var systemSettingResult dto.SystemSettingResult
	vars := make(map[string]string)
	mos, err := s.systemSettingRepo.List()
	if err != nil {
		return systemSettingResult, err
	}
	for _, mo := range mos {
		vars[mo.Key] = mo.Value
	}
	systemSettingResult.Vars = vars
	return systemSettingResult, err
}

func (s systemSettingService) ListByTab(tabName string) (dto.SystemSettingResult, error) {
	var systemSettingResult dto.SystemSettingResult
	vars := make(map[string]string)
	mos, err := s.systemSettingRepo.ListByTab(tabName)
	if err != nil {
		return systemSettingResult, err
	}
	for _, mo := range mos {
		vars[mo.Key] = mo.Value
	}
	if len(mos) > 0 {
		systemSettingResult.Tab = tabName
	}
	systemSettingResult.Vars = vars
	return systemSettingResult, err
}

func (s systemSettingService) Create(creation dto.SystemSettingCreate) ([]dto.SystemSetting, error) {

	var result []dto.SystemSetting
	for k, v := range creation.Vars {
		systemSetting, err := s.systemSettingRepo.Get(k)
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				systemSetting.Key = k
				systemSetting.Value = v
				systemSetting.Tab = creation.Tab
				err := s.systemSettingRepo.Save(&systemSetting)
				if err != nil {
					return result, err
				}
				result = append(result, dto.SystemSetting{SystemSetting: systemSetting})
			} else {
				return result, err
			}
		} else if systemSetting.ID != "" {
			systemSetting.Value = v
			if systemSetting.Tab == "" {
				systemSetting.Tab = creation.Tab
			}
			err := s.systemSettingRepo.Save(&systemSetting)
			if err != nil {
				return result, err
			}
			result = append(result, dto.SystemSetting{SystemSetting: systemSetting})
		}
	}
	return result, nil
}

func (s systemSettingService) GetLocalIPs() ([]model.SystemRegistry, error) {
	var sysRepo []model.SystemRegistry
	if err := db.DB.Find(&sysRepo).Error; err != nil {
		return sysRepo, fmt.Errorf("can't found repo from system registry, err %s", err.Error())
	}
	return sysRepo, nil
}

func (s systemSettingService) CheckSettingByType(tabName string, creation dto.SystemSettingCreate) error {

	vars := make(map[string]interface{})
	for k, value := range creation.Vars {
		vars[k] = value
	}
	if tabName == constant.Email {
		vars["type"] = constant.Email
		vars["RECEIVERS"] = vars["SMTP_TEST_USER"]
		vars["TITLE"] = "KubeOperator测试邮件"
		vars["CONTENT"] = "此邮件由 KubeOperator 发送，用于测试邮件发送，请勿回复"
	} else if tabName == constant.DingTalk {
		vars["type"] = constant.DingTalk
		vars["RECEIVERS"] = vars["DING_TALK_TEST_USER"]
		vars["TITLE"] = "KubeOperator测试消息"
		vars["CONTENT"] = "此邮件由 KubeOperator 发送，用于测试消息发送"
	} else if tabName == constant.WorkWeiXin {
		vars["type"] = constant.WorkWeiXin
		vars["CONTENT"] = "此邮件由 KubeOperator 发送，用于测试消息发送"
		vars["RECEIVERS"] = vars["WORK_WEIXIN_TEST_USER"]
	}
	c, err := message.NewMessageClient(vars)
	if err != nil {
		return err
	}
	if tabName == constant.WorkWeiXin {
		token, err := client.GetToken(vars)
		if err != nil {
			return err
		}
		vars["TOKEN"] = token
	}
	err = c.SendMessage(vars)
	if err != nil {
		return err
	}
	return nil
}

func (s systemSettingService) ListRegistry(conditions condition.Conditions) ([]dto.SystemRegistry, error) {
	var systemRegistryDto []dto.SystemRegistry
	var mos []model.SystemRegistry
	d := db.DB.Model(model.SystemRegistry{})
	if err := dbUtil.WithConditions(&d, model.User{}, conditions); err != nil {
		return nil, err
	}
	if err := d.Order("architecture").
		Find(&mos).Error; err != nil {
		return nil, err
	}
	for _, mo := range mos {
		systemRegistryDto = append(systemRegistryDto, dto.SystemRegistry{
			SystemRegistry: mo,
		})
	}
	return systemRegistryDto, nil
}

func (s systemSettingService) GetRegistryByID(id string) (dto.SystemRegistry, error) {
	r, err := s.systemRegistryRepo.Get(id)
	if err != nil {
		return dto.SystemRegistry{}, err
	}
	systemRegistryDto := dto.SystemRegistry{
		SystemRegistry: model.SystemRegistry{
			ID:                 r.ID,
			Hostname:           r.Hostname,
			Protocol:           r.Protocol,
			Architecture:       r.Architecture,
			RepoPort:           r.RepoPort,
			RegistryPort:       r.RegistryPort,
			RegistryHostedPort: r.RegistryHostedPort,
		},
	}
	return systemRegistryDto, nil
}

func (s systemSettingService) GetRegistryByArch(arch string) (dto.SystemRegistry, error) {
	r, err := s.systemRegistryRepo.GetByArch(arch)
	if err != nil {
		return dto.SystemRegistry{}, err
	}
	systemRegistryDto := dto.SystemRegistry{
		SystemRegistry: model.SystemRegistry{
			ID:           r.ID,
			Hostname:     r.Hostname,
			Protocol:     r.Protocol,
			Architecture: r.Architecture,
		},
	}
	return systemRegistryDto, nil
}

func (s systemSettingService) PageRegistry(num, size int, conditions condition.Conditions) (*page.Page, error) {
	var (
		p                 page.Page
		systemRegistryDto []dto.SystemRegistry
		mos               []model.SystemRegistry
	)

	d := db.DB.Model(model.SystemRegistry{})
	if err := dbUtil.WithConditions(&d, model.SystemRegistry{}, conditions); err != nil {
		return nil, err
	}
	if err := d.
		Count(&p.Total).
		Order("architecture").
		Offset((num - 1) * size).
		Limit(size).
		Find(&mos).Error; err != nil {
		return nil, err
	}
	for _, mo := range mos {
		systemRegistryDto = append(systemRegistryDto, dto.SystemRegistry{SystemRegistry: mo})
	}
	p.Items = systemRegistryDto
	return &p, nil
}

func (s systemSettingService) CreateRegistry(creation dto.SystemRegistryCreate) (*dto.SystemRegistry, error) {
	systemRegistry := model.SystemRegistry{
		ID:                 creation.ID,
		Architecture:       creation.Architecture,
		Protocol:           creation.Protocol,
		Hostname:           creation.Hostname,
		RepoPort:           creation.RepoPort,
		RegistryPort:       creation.RegistryPort,
		RegistryHostedPort: creation.RegistryHostedPort,
	}
	err := s.systemRegistryRepo.Save(&systemRegistry)
	if err != nil {
		return nil, err
	}
	return &dto.SystemRegistry{SystemRegistry: systemRegistry}, nil
}

func (s systemSettingService) UpdateRegistry(arch string, creation dto.SystemRegistryUpdate) (*dto.SystemRegistry, error) {
	systemRegistry := model.SystemRegistry{
		ID:                 creation.ID,
		Architecture:       arch,
		Protocol:           creation.Protocol,
		Hostname:           creation.Hostname,
		RepoPort:           creation.RepoPort,
		RegistryPort:       creation.RegistryPort,
		RegistryHostedPort: creation.RegistryHostedPort,
	}
	err := s.systemRegistryRepo.Save(&systemRegistry)
	if err != nil {
		return nil, err
	}
	return &dto.SystemRegistry{SystemRegistry: systemRegistry}, nil
}

func (s systemSettingService) BatchRegistry(op dto.SystemRegistryBatchOp) error {
	var deleteItems []model.SystemRegistry
	for _, item := range op.Items {
		deleteItems = append(deleteItems, model.SystemRegistry{
			BaseModel:    common.BaseModel{},
			ID:           item.ID,
			Architecture: item.Architecture,
		})
	}
	err := s.systemRegistryRepo.Batch(op.Operation, deleteItems)
	if err != nil {
		return err
	}
	return nil
}

func (s systemSettingService) DeleteRegistry(id string) error {
	err := s.systemRegistryRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
