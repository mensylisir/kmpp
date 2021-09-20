package dto

import "github.com/kmpp/pkg/model"

type Zone struct {
	model.Zone
	CloudVars  interface{} `json:"cloudVars"`
	RegionName string      `json:"regionName"`
	Provider   string      `json:"provider"`
	IpPoolName string      `json:"ipPoolName"`
	IpPool     IpPool      `json:"ipPool"`
}

type ZoneCreate struct {
	Name           string      `json:"name" validate:"required"`
	CloudVars      interface{} `json:"cloudVars" validate:"required"`
	RegionName     string      `json:"regionName" validate:"required"`
	IpPoolName     string      `json:"ipPoolName"`
	CredentialName string      `json:"credentialName"`
}

type ZoneOp struct {
	Operation string `json:"operation" validate:"required"`
	Items     []Zone `json:"items" validate:"required"`
}

type CloudZoneResponse struct {
	Result interface{} `json:"result"`
}

type CloudZoneRequest struct {
	CloudVars  interface{} `json:"cloudVars" validate:"required"`
	Datacenter string      `json:"datacenter"`
	RegionName string      `json:"RegionName"`
}

type ZoneUpdate struct {
	ID         string      `json:"id" validate:"required"`
	Name       string      `json:"name" validate:"required"`
	CloudVars  interface{} `json:"cloudVars" validate:"required"`
	RegionID   string      `json:"regionID" validate:"required"`
	IpPoolName string      `json:"ipPoolName" validate:"required"`
}

type CloudDatastore struct {
	Name      string `json:"name"`
	Capacity  int    `json:"capacity"`
	FreeSpace int    `json:"freeSpace"`
}
