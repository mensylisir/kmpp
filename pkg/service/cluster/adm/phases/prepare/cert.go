package prepare

import (
	"github.com/kmpp/pkg/service/cluster/adm/phases"
	"github.com/kmpp/pkg/util/kobe"
	"io"
)

const (
	prepareCertificates = "05-certificates.yml"
)

type CertificatesPhase struct {
}

func (c CertificatesPhase) Name() string {
	return "GenerateCertificates"
}

func (c CertificatesPhase) Run(b kobe.Interface, writer io.Writer) error {
	return phases.RunPlaybookAndGetResult(b, prepareCertificates, "", writer)
}
