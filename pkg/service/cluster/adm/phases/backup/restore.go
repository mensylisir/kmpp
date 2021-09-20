package backup

import (
	"github.com/kmpp/pkg/service/cluster/adm/phases"
	"github.com/kmpp/pkg/util/kobe"
	"io"
)

const (
	restoreCluster = "95-restore-cluster.yml"
)

type RestoreClusterPhase struct {
}

func (restore RestoreClusterPhase) Name() string {
	return "backupCluster"
}

func (restore RestoreClusterPhase) Run(b kobe.Interface, writer io.Writer) error {
	return phases.RunPlaybookAndGetResult(b, restoreCluster, "", writer)
}
