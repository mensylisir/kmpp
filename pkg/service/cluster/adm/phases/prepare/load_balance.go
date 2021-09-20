package prepare

import (
	"github.com/kmpp/pkg/service/cluster/adm/phases"
	"github.com/kmpp/pkg/util/kobe"
	"io"
)

const (
	prepareLoadBalancer = "04-load-balancer.yml"
)

type LoadBalancerPhase struct {
}

func (s LoadBalancerPhase) Name() string {
	return "Install Load Balancer"
}

func (s LoadBalancerPhase) Run(b kobe.Interface, writer io.Writer) error {
	return phases.RunPlaybookAndGetResult(b, prepareLoadBalancer, "", writer)
}
