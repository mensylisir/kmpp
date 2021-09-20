package dto

import "github.com/kmpp/pkg/model"

type ProjectResource struct {
	model.ProjectResource
	ResourceName string `json:"resourceName"`
}

type ProjectResourceOp struct {
	Operation string            `json:"operation" validate:"required"`
	Items     []ProjectResource `json:"items" validate:"required"`
}

type ProjectResourceTree struct {
	ID       int                   `json:"id"`
	Label    string                `json:"label"`
	Type     string                `json:"type"`
	Children []ProjectResourceTree `json:"children"`
}

type ProjectResourceCreate struct {
	ResourceType string   `json:"resourceType" validate:"required"`
	Names        []string `json:"names" validate:"required"`
}
