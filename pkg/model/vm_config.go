package model

import (
	"encoding/json"
	"errors"

	"github.com/kmpp/pkg/db"
	"github.com/kmpp/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type VmConfig struct {
	common.BaseModel
	ID       string `json:"-"`
	Name     string `json:"name"`
	Cpu      int    `json:"cpu"`
	Memory   int    `json:"memory"`
	Disk     int    `json:"disk"`
	Provider string `json:"provider"`
}

func (v *VmConfig) BeforeCreate() error {
	v.ID = uuid.NewV4().String()
	return nil
}

func (v *VmConfig) BeforeDelete() error {
	var plans []Plan
	if err := db.DB.Find(&plans).Error; err != nil {
		return err
	}
	for _, p := range plans {
		planVars := map[string]string{}
		_ = json.Unmarshal([]byte(p.Vars), &planVars)
		if planVars["masterModel"] == v.Name {
			return errors.New("VM_CONFIG_DELETE_FAILED")
		}
		if planVars["workerModel"] == v.Name {
			return errors.New("VM_CONFIG_DELETE_FAILED")
		}
	}
	return nil
}
