package model

import (
	"github.com/kmpp/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type ClusterManifest struct {
	common.BaseModel
	ID          string `json:"-"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	CoreVars    string `json:"coreVars"`
	NetworkVars string `json:"networkVars"`
	ToolVars    string `json:"toolVars"`
	StorageVars string `json:"storageVars"`
	OtherVars   string `json:"otherVars"`
	IsActive    bool   `json:"isActive"`
}

func (m *ClusterManifest) BeforeCreate() (err error) {
	m.ID = uuid.NewV4().String()
	return nil
}

type VersionHelp struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
