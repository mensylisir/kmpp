package dto

import "github.com/kmpp/pkg/model"

type ClusterTool struct {
	model.ClusterTool
	Vars map[string]interface{} `json:"vars"`
}
