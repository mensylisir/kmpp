package prepare

import (
	"github.com/kmpp/pkg/service/cluster/adm/phases"
	"github.com/kmpp/pkg/util/kobe"
	"io"
)

const (
	prepareKubernetesComponents = "03-kubernetes-component.yml"
)

type KubernetesComponentPhase struct {
}

func (s KubernetesComponentPhase) Name() string {
	return "Prepare Kubernetes Component"
}

func (s KubernetesComponentPhase) Run(b kobe.Interface, writer io.Writer) error {
	return phases.RunPlaybookAndGetResult(b, prepareKubernetesComponents, "", writer)
}
