package dto

import "github.com/kmpp/pkg/model"

type UserMessageDTO struct {
	model.UserMessage
	MsgContent  interface{} `json:"msgContent"`
	ClusterName string      `json:"clusterName"`
}

type UserMessageOp struct {
	Operation string           `json:"operation"`
	Items     []UserMessageDTO `json:"items"`
}

type UnReadMessage struct {
	Info    int `json:"info"`
	Warning int `json:"warning"`
}
