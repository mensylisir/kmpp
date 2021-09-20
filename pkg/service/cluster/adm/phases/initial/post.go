package initial

import (
	"github.com/kmpp/pkg/service/cluster/adm/phases"
	"github.com/kmpp/pkg/util/kobe"
	"io"
)

const (
	initPost = "15-post.yml"
)

type PostPhase struct {
}

func (s PostPhase) Name() string {
	return "Post Init"
}

func (s PostPhase) Run(b kobe.Interface, writer io.Writer) error {
	return phases.RunPlaybookAndGetResult(b, initPost, "", writer)
}
