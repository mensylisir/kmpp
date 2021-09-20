package constant

import (
	"errors"
	"path"
)

const (
	PageNumQueryKey  = "pageNum"
	PageSizeQueryKey = "pageSize"

	BatchOperationUpdate = "update"
	BatchOperationCreate = "create"
	BatchOperationDelete = "delete"

	LocalRepositoryDomainName = "registry.kubeoperator.io"

	DefaultResourceName = "kubeoperator"
	StatusPending       = "Pending"
	StatusRunning       = "Running"
	StatusNotReady      = "NotReady"
	StatusUpgrading     = "Upgrading"
	StatusSuccess       = "Success"
	StatusFailed        = "Failed"
	StatusLost          = "Lost"
	StatusCreating      = "Creating"
	StatusInitializing  = "Initializing"
	StatusTerminating   = "Terminating"
	StatusWaiting       = "Waiting"

	DefaultPassword = "kubeoperator@admin123"
)

var (
	NotSupportedBatchOperation = errors.New("not supported operation")
)

var (
	ResourceDir          = "resource"
	ChartsDir            = path.Join(ResourceDir, "charts")
	DefaultDataDir       = "/home/mensyli1/Documents/Workspace/kmpp/data"
	DefaultAnsibleLogDir = path.Join(DefaultDataDir, "ansible")
	BackupDir            = path.Join(DefaultDataDir, "backup")
	DefaultRepositoryDir = path.Join(DefaultDataDir, "git")
)
