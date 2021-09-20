package kolog

import (
	"github.com/kmpp/pkg/dto"
	"github.com/kmpp/pkg/logger"
	"github.com/kmpp/pkg/service"
)

func Save(name, operation, operationInfo string) {
	lS := service.NewSystemLogService()
	logInfo := dto.SystemLogCreate{
		Name:          name,
		Operation:     operation,
		OperationInfo: operationInfo,
	}
	if err := lS.Create(logInfo); err != nil {
		logger.Log.Errorf("save system logs failed, error: %s", err.Error())
	}
}
