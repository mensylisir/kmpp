package dto

import "github.com/kmpp/pkg/model"

type VmConfig struct {
	model.VmConfig
}

type VmConfigOp struct {
	Operation string     `json:"operation"`
	Items     []VmConfig `json:"items"`
}

type VmConfigCreate struct {
	Name     string `json:"name" validate:"required"`
	Provider string `json:"provider"`
	Cpu      int    `json:"cpu" validate:"min=1,max=1000,required" en:"CPU" zh:"CPU"`
	Memory   int    `json:"memory"  validate:"min=1,max=1000,required" en:"Memory" zh:"内存"`
}

type VmConfigUpdate struct {
	Name     string `json:"name"`
	Provider string `json:"provider"`
	Cpu      int    `json:"cpu" validate:"min=1,max=1000,required" en:"CPU" zh:"CPU"`
	Memory   int    `json:"memory" validate:"min=1,max=1000,required" en:"Memory" zh:"内存"`
}
