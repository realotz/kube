package config

import "github.com/go-kratos/kratos/v2/config/source"

type watcher struct {
}

func newWatcher(k *kube) (source.Watcher, error) {
	return &watcher{}, nil
}

func (w *watcher) Next() (*source.KeyValue, error) {
	return nil, nil
}

func (w *watcher) Close() error {
	return nil
}
