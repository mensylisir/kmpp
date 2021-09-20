package dto

import "github.com/kmpp/pkg/model"

type CisTask struct {
	model.CisTask
}

type CisResult struct {
	model.CisTaskResult
}

type CisBatch struct {
	Items     []CisTask
	Operation string
}
