package dto

import "github.com/kmpp/pkg/model"

type ClusterMember struct {
	model.ClusterMember
	ClusterName string `json:"clusterName"`
	Username    string `json:"username"`
	Email       string `json:"email"`
}

type UsersResponse struct {
	Items []string `json:"items"`
}

type ClusterMemberCreate struct {
	Usernames []string `json:"usernames" validate:"required"`
}
