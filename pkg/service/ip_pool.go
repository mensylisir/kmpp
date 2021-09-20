package service

import (
	"github.com/kmpp/pkg/constant"
	"github.com/kmpp/pkg/controller/condition"
	"github.com/kmpp/pkg/controller/page"
	"github.com/kmpp/pkg/db"
	"github.com/kmpp/pkg/dto"
	"github.com/kmpp/pkg/model"
	"github.com/kmpp/pkg/model/common"
	"github.com/kmpp/pkg/repository"
	dbUtil "github.com/kmpp/pkg/util/db"
)

type IpPoolService interface {
	Get(name string) (dto.IpPool, error)
	Page(num, size int, conditions condition.Conditions) (*page.Page, error)
	Create(creation dto.IpPoolCreate) (dto.IpPool, error)
	Batch(op dto.IpPoolOp) error
	List(conditions condition.Conditions) ([]dto.IpPool, error)
	Delete(name string) error
}

type ipPoolService struct {
	ipPoolRepo repository.IpPoolRepository
	ipService  IpService
}

func NewIpPoolService() IpPoolService {
	return &ipPoolService{
		ipPoolRepo: repository.NewIpPoolRepository(),
		ipService:  NewIpService(),
	}
}

func (i ipPoolService) Get(name string) (dto.IpPool, error) {
	var ipPoolDTO dto.IpPool
	ipPool, err := i.ipPoolRepo.Get(name)
	if err != nil {
		return ipPoolDTO, err
	}
	ipPoolDTO.IpPool = ipPool
	return ipPoolDTO, nil
}

func (i ipPoolService) Page(num, size int, conditions condition.Conditions) (*page.Page, error) {
	var (
		p          page.Page
		ipPoolDTOS []dto.IpPool
		ipPools    []model.IpPool
	)
	d := db.DB.Model(model.IpPool{})
	if err := dbUtil.WithConditions(&d, model.IpPool{}, conditions); err != nil {
		return nil, err
	}

	if err := d.Count(&p.Total).Preload("Ips").Offset((num - 1) * size).Limit(size).Find(&ipPools).Error; err != nil {
		return nil, err
	}
	for _, mo := range ipPools {
		ipUsed := 0
		for _, ip := range mo.Ips {
			if ip.Status != constant.IpAvailable {
				ipUsed++
			}
		}
		ipPoolDTOS = append(ipPoolDTOS, dto.IpPool{
			IpPool: mo,
			IpUsed: ipUsed,
		})
	}
	p.Items = ipPoolDTOS
	return &p, nil
}

func (i ipPoolService) List(conditions condition.Conditions) ([]dto.IpPool, error) {
	var ipPoolDTOS []dto.IpPool
	var ipPools []model.IpPool
	d := db.DB.Model(model.IpPool{})
	if err := dbUtil.WithConditions(&d, model.IpPool{}, conditions); err != nil {
		return nil, err
	}
	err := d.Preload("Ips").Find(&ipPools).Error
	if err != nil {
		return ipPoolDTOS, err
	}
	for _, mo := range ipPools {
		ipUsed := 0
		for _, ip := range mo.Ips {
			if ip.Status != constant.IpAvailable {
				ipUsed++
			}
		}
		ipPoolDTOS = append(ipPoolDTOS, dto.IpPool{
			IpPool: mo,
			IpUsed: ipUsed,
		})
	}
	return ipPoolDTOS, nil
}

func (i ipPoolService) Create(creation dto.IpPoolCreate) (dto.IpPool, error) {
	var ipPoolDTO dto.IpPool
	ipPool := model.IpPool{
		BaseModel:   common.BaseModel{},
		Name:        creation.Name,
		Description: creation.Description,
		Subnet:      creation.Subnet,
	}
	tx := db.DB.Begin()
	err := tx.Create(&ipPool).Error
	if err != nil {
		tx.Rollback()
		return ipPoolDTO, err
	}
	err = i.ipService.Create(dto.IpCreate{
		IpStart:    creation.IpStart,
		IpEnd:      creation.IpEnd,
		Gateway:    creation.Gateway,
		IpPoolName: ipPool.Name,
		DNS1:       creation.DNS1,
		DNS2:       creation.DNS2,
	}, tx)
	if err != nil {
		tx.Rollback()
		return ipPoolDTO, err
	}
	tx.Commit()
	ipPoolDTO.IpPool = ipPool
	return ipPoolDTO, err
}

func (i ipPoolService) Delete(name string) error {
	ipPool, err := i.Get(name)
	if err != nil {
		return err
	}
	tx := db.DB.Begin()
	if err := tx.Delete(&ipPool).Error; err != nil {
		return err
	}
	if err := tx.Where("ip_pool_id = ?", ipPool.ID).Delete(&model.Ip{}).Error; err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (i ipPoolService) Batch(op dto.IpPoolOp) error {
	var opItems []model.IpPool
	for _, item := range op.Items {
		opItems = append(opItems, model.IpPool{
			BaseModel: common.BaseModel{},
			Name:      item.Name,
		})
	}
	tx := db.DB.Begin()
	switch op.Operation {
	case constant.BatchOperationDelete:
		for i := range opItems {
			var ipPool model.IpPool
			if err := tx.Where("name = ?", opItems[i].Name).First(&ipPool).Error; err != nil {
				tx.Rollback()
				return err
			}
			if err := tx.Delete(&ipPool).Error; err != nil {
				tx.Rollback()
				return err
			}

			if err := tx.Where("ip_pool_id = ?", ipPool.ID).Delete(&model.Ip{}).Error; err != nil {
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
