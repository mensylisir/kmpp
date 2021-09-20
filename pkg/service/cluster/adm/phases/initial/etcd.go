package initial

import (
	"github.com/kmpp/pkg/service/cluster/adm/phases"
	"github.com/kmpp/pkg/util/kobe"
	"io"
)

const (
	initEtcd = "06-etcd.yml"
)

type EtcdPhase struct {
	Upgrade bool
}

func (s EtcdPhase) Name() string {
	return "InitEtcd"
}

func (s EtcdPhase) Run(b kobe.Interface, writer io.Writer) error {
	var tag string
	if s.Upgrade {
		tag = "upgrade"
	}
	return phases.RunPlaybookAndGetResult(b, initEtcd, tag, writer)
}
