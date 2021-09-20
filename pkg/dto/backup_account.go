package dto

import "github.com/kmpp/pkg/model"

type BackupAccount struct {
	model.BackupAccount
	CredentialVars interface{} `json:"credentialVars"`
	Projects       string      `json:"projects"`
	Clusters       string      `json:"clusters"`
}

type BackupAccountOp struct {
	Operation string          `json:"operation" validate:"required"`
	Items     []BackupAccount `json:"items" validate:"required"`
}

type BackupAccountRequest struct {
	Name           string      `json:"name" validate:"required"`
	CredentialVars interface{} `json:"credentialVars" validate:"required"`
	Bucket         string      `json:"bucket" validate:"required"`
	Type           string      `json:"type" validate:"required"`
	Projects       []string    `json:"projects"`
	Clusters       []string    `json:"clusters"`
}

type BackupAccountUpdate struct {
	ID             string      `json:"id" validate:"required"`
	Name           string      `json:"name" validate:"required"`
	CredentialVars interface{} `json:"credentialVars" validate:"required"`
	Bucket         string      `json:"bucket" validate:"required"`
	Type           string      `json:"type" validate:"required"`
	Projects       []string    `json:"projects"`
}

type CloudStorageRequest struct {
	CredentialVars interface{} `json:"credentialVars" validate:"required"`
	Type           string      `json:"type" validate:"required"`
}
