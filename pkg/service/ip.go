package service

import (
	"errors"
	"strconv"
	"strings"

	"github.com/kmpp/pkg/controller/condition"
	dbUtil "github.com/kmpp/pkg/util/db"

	"github.com/kmpp/pkg/constant"
	"github.com/kmpp/pkg/controller/page"
	"github.com/kmpp/pkg/db"
	"github.com/kmpp/pkg/dto"
	"github.com/kmpp/pkg/model"
	"github.com/kmpp/pkg/util/ipaddr"
	"github.com/jinzhu/gorm"
)

type IpService interface {
	Get(ip string) (dto.Ip, error)
	Create(create dto.IpCreate, tx *gorm.DB) error
	Page(num, size int, ipPoolName string, conditions condition.Conditions) (*page.Page, error)
	Batch(op dto.IpOp) error
	Update(name string, update dto.IpUpdate) (*dto.Ip, error)
	Sync(ipPoolName string) error
	Delete(address string) error
	List(ipPoolName string, conditions condition.Conditions) ([]dto.Ip, error)
}

type ipService struct {
}

func NewIpService() IpService {
	return &ipService{}
}

func (i ipService) Get(ip string) (dto.Ip, error) {
	var ipDTO dto.Ip
	var ipM model.Ip
	if err := db.DB.Where("address = ?", ip).First(&ipM).Error; err != nil {
		return ipDTO, err
	}
	ipDTO = dto.Ip{
		Ip: ipM,
	}
	return ipDTO, nil
}

func (i ipService) List(ipPoolName string, conditions condition.Conditions) ([]dto.Ip, error) {
	var ipDTOS []dto.Ip
	var mos []model.Ip
	var ipPool model.IpPool
	err := db.DB.Where("name = ?", ipPoolName).First(&ipPool).Error
	if err != nil {
		return nil, err
	}
	d := db.DB.Model(model.Ip{})
	if err := dbUtil.WithConditions(&d, model.Ip{}, conditions); err != nil {
		return nil, err
	}
	if err := d.Order("inet_aton(address)").
		Find(&mos).Error; err != nil {
		return nil, err
	}
	for _, mo := range mos {
		ipDTOS = append(ipDTOS, dto.Ip{
			Ip: mo,
		})
	}
	return ipDTOS, nil
}

func (i ipService) Create(create dto.IpCreate, tx *gorm.DB) error {
	if tx == nil {
		tx = db.DB.Begin()
	}
	var ipPool model.IpPool
	if err := tx.Where("name = ?", create.IpPoolName).First(&ipPool).Error; err != nil {
		return err
	}
	cs := strings.Split(ipPool.Subnet, "/")
	mask, _ := strconv.Atoi(cs[1])
	startIp := strings.Replace(create.IpStart, " ", "", -1)
	endIp := strings.Replace(create.IpEnd, " ", "", -1)
	if !(ipaddr.CheckIP(startIp) && ipaddr.CheckIP(endIp)) {
		return errors.New("IP_INVALID")
	}
	ips := ipaddr.GenerateIps(cs[0], mask, startIp, endIp)
	if len(ips) == 0 {
		return errors.New("IP_NULL")
	}
	for _, ip := range ips {
		var old model.Ip
		tx.Where("address = ?", ip).First(&old)
		if old.ID != "" {
			tx.Rollback()
			return errors.New("IP_EXISTS")
		}
		insert := model.Ip{
			Address:  ip,
			Gateway:  create.Gateway,
			DNS1:     create.DNS1,
			DNS2:     create.DNS2,
			IpPoolID: ipPool.ID,
			Status:   constant.IpAvailable,
		}
		err := tx.Create(&insert).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		go func() {
			err := ipaddr.Ping(insert.Address)
			if err == nil {
				insert.Status = constant.IpReachable
				db.DB.Save(&insert)
			}
		}()
	}
	tx.Commit()
	return nil
}

func (i ipService) Page(num, size int, ipPoolName string, conditions condition.Conditions) (*page.Page, error) {
	var (
		p      page.Page
		ipDTOS []dto.Ip
		ips    []model.Ip
		ipPool model.IpPool
	)

	err := db.DB.Where("name = ?", ipPoolName).First(&ipPool).Error
	if err != nil {
		return nil, err
	}

	d := db.DB.Model(model.Ip{})
	if err := dbUtil.WithConditions(&d, model.Ip{}, conditions); err != nil {
		return nil, err
	}

	if err := d.
		Where("ip_pool_id = ?", ipPool.ID).
		Count(&p.Total).
		Order("inet_aton(address)").
		Offset((num - 1) * size).
		Limit(size).
		Find(&ips).Error; err != nil {
		return nil, err
	}

	for _, mo := range ips {
		ipDTOS = append(ipDTOS, dto.Ip{
			Ip: mo,
		})
	}
	p.Items = ipDTOS
	return &p, nil
}

func (i ipService) Batch(op dto.IpOp) error {
	tx := db.DB.Begin()
	switch op.Operation {
	case constant.BatchOperationDelete:
		for i := range op.Items {
			var ip model.Ip
			if err := tx.Where("address = ?", op.Items[i].Address).First(&ip).Error; err != nil {
				tx.Rollback()
				return err
			}
			if err := tx.Delete(&ip).Error; err != nil {
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

func (i ipService) Update(address string, update dto.IpUpdate) (*dto.Ip, error) {
	tx := db.DB.Begin()
	var ip model.Ip
	err := tx.Where("address = ?", address).First(&ip).Error
	if err != nil {
		return nil, err
	}
	switch update.Operation {
	case "LOCK":
		ip.Status = constant.IpLock
	case "UNLOCK":
		ip.Status = constant.IpAvailable
	default:
		break
	}
	err = tx.Save(&ip).Error
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return &dto.Ip{Ip: ip}, err
}

func (i ipService) Delete(address string) error {
	ip, err := i.Get(address)
	if err != nil {
		return err
	}
	if err := db.DB.Delete(&ip).Error; err != nil {
		return err
	}
	return nil
}

func (i ipService) Sync(ipPoolName string) error {
	var ipPool model.IpPool
	err := db.DB.Where("name = ?", ipPoolName).First(&ipPool).Error
	if err != nil {
		return err
	}
	var ips []model.Ip
	err = db.DB.Where("ip_pool_id = ?", ipPool.ID).Find(&ips).Error
	if err != nil {
		return err
	}
	for i := range ips {
		if ips[i].Status == constant.IpLock {
			continue
		}
		var host model.Host
		db.DB.Model(model.Host{}).Where("ip = ?", ips[i].Address).Find(&host)
		if host.ID != "" {
			if ips[i].Status == constant.IpUsed {
				continue
			} else {
				ips[i].Status = constant.IpUsed
				db.DB.Save(&ips[i])
				continue
			}
		}
		go func(i int) {
			err := ipaddr.Ping(ips[i].Address)
			if err == nil && ips[i].Status != constant.IpReachable {
				ips[i].Status = constant.IpReachable
				db.DB.Save(&ips[i])
			}
			if err != nil && ips[i].Status == constant.IpReachable {
				ips[i].Status = constant.IpAvailable
				db.DB.Save(&ips[i])
			}
		}(i)
	}
	return nil
}
