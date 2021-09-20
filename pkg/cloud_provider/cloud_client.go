package cloud_provider

import (
	"github.com/kmpp/pkg/cloud_provider/client"
	"github.com/kmpp/pkg/constant"
)

type CloudClient interface {
	ListDatacenter() ([]string, error)
	ListClusters() ([]interface{}, error)
	ListTemplates() ([]interface{}, error)
	ListFlavors() ([]interface{}, error)
	GetIpInUsed(network string) ([]string, error)
	UploadImage() error
	DefaultImageExist() (bool, error)
	CreateDefaultFolder() error
	ListDatastores() ([]client.DatastoreResult, error)
}

func NewCloudClient(vars map[string]interface{}) CloudClient {
	switch vars["provider"] {
	case constant.OpenStack:
		return client.NewOpenStackClient(vars)
	case constant.VSphere:
		return client.NewVSphereClient(vars)
	case constant.FusionCompute:
		return client.NewFusionComputeClient(vars)
	}
	return nil
}
