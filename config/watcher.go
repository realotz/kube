package config

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/config/source"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type watcher struct {
	k       *kube
	watcher watch.Interface
	ctx     context.Context
	close   context.CancelFunc
}

func newWatcher(k *kube) (source.Watcher, error) {
	if k.client == nil {
		if err := k.initKubeClient(); err != nil {
			return nil, err
		}
	}
	ctx, closeFunc := context.WithCancel(context.Background())
	cmWatcher, err := k.client.CoreV1().ConfigMaps(k.op.Namespace).Watch(ctx, metav1.ListOptions{
		LabelSelector: k.op.LabelSelector,
		FieldSelector: k.op.FieldSelector,
	})
	if err != nil {
		return nil, err
	}
	return &watcher{
		k:       k,
		watcher: cmWatcher,
		ctx:     ctx,
		close:   closeFunc,
	}, nil
}

func (w *watcher) Next() ([]*source.KeyValue, error) {
	c := <-w.watcher.ResultChan()
	if c.Object == nil {
		return nil, fmt.Errorf("kube config watcher close")
	}
	sv := c.Object.(*v1.ConfigMap)
	if c.Type == "DELETED" {
		return nil, fmt.Errorf("kube configmap delete %s", sv.Name)
	}
	return w.k.configMapKV(*sv), nil
}

func (w *watcher) Close() error {
	w.close()
	return nil
}
