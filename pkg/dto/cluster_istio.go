package dto

import "github.com/kmpp/pkg/model"

type ClusterIstio struct {
	ClusterIstio model.ClusterIstio     `json:"cluster_istio"`
	Operation    string                 `json:"operation"`
	Enable       bool                   `json:"enable"`
	Vars         map[string]interface{} `json:"vars"`
}
