package registry

import (
	"github.com/go-kratos/kratos/v2/registry"
)

type watcher struct {
}

func newWatcher(k *kube) (registry.Watcher, error) {
	return &watcher{}, nil
}

func (w *watcher) Next() ([]*registry.ServiceInstance, error) {
	return nil, nil
}

func (w *watcher) Close() error {
	return nil
}
