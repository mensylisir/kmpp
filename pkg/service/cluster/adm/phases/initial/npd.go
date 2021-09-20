package initial

import (
	"github.com/kmpp/pkg/service/cluster/adm/phases"
	"github.com/kmpp/pkg/util/kobe"
	"io"
)

const (
	initNpd = "12-npd.yml"
)

type NpdPhase struct {
}

func (s NpdPhase) Name() string {
	return "Npd Init"
}

func (s NpdPhase) Run(b kobe.Interface, writer io.Writer) error {
	return phases.RunPlaybookAndGetResult(b, initNpd, "", writer)
}
