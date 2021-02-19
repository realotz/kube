package config

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/config"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type watcher struct {
	k       *kube
	ctx     context.Context
	cancel  context.CancelFunc
	watcher watch.Interface
}

func newWatcher(k *kube) (config.Watcher, error) {
	ctx, cancel := context.WithCancel(context.Background())
	w, err := k.client.CoreV1().ConfigMaps(k.opts.Namespace).Watch(ctx, metav1.ListOptions{
		LabelSelector: k.opts.LabelSelector,
		FieldSelector: k.opts.FieldSelector,
	})
	if err != nil {
		cancel()
		return nil, err
	}
	return &watcher{
		k:       k,
		watcher: w,
		ctx:     ctx,
		cancel:  cancel,
	}, nil
}

func (w *watcher) Next() ([]*config.KeyValue, error) {
	ch := <-w.watcher.ResultChan()
	if ch.Object == nil {
		return nil, fmt.Errorf("kube config watcher close")
	}
	cm := ch.Object.(*v1.ConfigMap)
	if ch.Type == "DELETED" {
		return nil, fmt.Errorf("kube configmap delete %s", cm.Name)
	}
	return w.k.configMap(*cm), nil
}

func (w *watcher) Close() error {
	w.cancel()
	return nil
}
