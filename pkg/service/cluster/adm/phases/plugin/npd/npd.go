package npd

import (
	"github.com/kmpp/pkg/service/cluster/adm/phases"
	"github.com/kmpp/pkg/util/kobe"
	"io"
)

const (
	npdPlaybook = "12-npd.yml"
)

type NpdPhase struct {
}

func (NpdPhase) Name() string {
	return "Npd"
}

func (c NpdPhase) Run(b kobe.Interface, writer io.Writer) error {
	return phases.RunPlaybookAndGetResult(b, npdPlaybook, "", writer)
}
