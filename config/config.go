package config

import "github.com/go-kratos/kratos/v2/config/source"

type kube struct {
}

// NewSource new a kube config source.
func NewSource() source.Source {
	return &kube{}
}

func (k *kube) Load() ([]*source.KeyValue, error) {
	return nil, nil
}
func (k *kube) Watch() (source.Watcher, error) {
	return newWatcher(k)
}
