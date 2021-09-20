package job

import (
	"sync"

	"github.com/kmpp/pkg/constant"
	"github.com/kmpp/pkg/db"
	"github.com/kmpp/pkg/logger"
	"github.com/kmpp/pkg/model"
	"github.com/kmpp/pkg/service"
)

type RefreshHostInfo struct {
	hostService service.HostService
}

func NewRefreshHostInfo() *RefreshHostInfo {
	return &RefreshHostInfo{
		hostService: service.NewHostService(),
	}
}

func (r *RefreshHostInfo) Run() {
	var hosts []model.Host
	var wg sync.WaitGroup
	sem := make(chan struct{}, 2) // 信号量
	db.DB.Find(&hosts)
	for _, host := range hosts {
		if host.Status == constant.ClusterCreating || host.Status == constant.ClusterInitializing || host.Status == constant.ClusterSynchronizing {
			continue
		}
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			_, err := r.hostService.Sync(name)
			if err != nil {
				logger.Log.Errorf("gather host info error: %s", err.Error())
			}
		}(host.Name)
	}
	wg.Wait()
}
