package config

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/config/source"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// New SourceOption
type SourceOption struct {
	// kube namespace
	Namespace string
	// kube labelSelector example `app=test`
	LabelSelector string
	// kube fieldSelector example `app=test`
	FieldSelector string
	// set KubeConfig out-of-cluster Use outside cluster
	KubeConfig string
}

type kube struct {
	op     SourceOption
	client *kubernetes.Clientset
}

// NewSource new a kube config source.
func NewSource(op SourceOption) source.Source {
	return &kube{
		op: op,
	}
}

// init kube client
func (k *kube) initKubeClient() (err error) {
	var config *rest.Config
	if k.op.KubeConfig != "" {
		if config, err = clientcmd.BuildConfigFromFlags("", k.op.KubeConfig); err != nil {
			return err
		}
	} else {
		if config, err = rest.InClusterConfig(); err != nil {
			return err
		}
	}
	k.client, err = kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	return nil
}

// loadKubeConfig
func (k *kube) loadKubeConfig() (kvs []*source.KeyValue, err error) {
	cmList, err := k.client.CoreV1().ConfigMaps(k.op.Namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: k.op.LabelSelector,
		FieldSelector: k.op.FieldSelector,
	})
	if err != nil {
		return nil, err
	}
	for _, cm := range cmList.Items {
		kvs = append(kvs, k.configMapKV(cm)...)
	}
	return kvs, nil
}

// configMapKV
func (k *kube) configMapKV(cm v1.ConfigMap) (kvs []*source.KeyValue) {
	for name, val := range cm.Data {
		kvs = append(kvs, &source.KeyValue{
			Key:       fmt.Sprintf("%s/%s/%s", k.op.Namespace, cm.Name, name),
			Value:     string2byte(val),
			Format:    format(name),
			Timestamp: cm.GetCreationTimestamp().Time,
		})
	}
	return kvs
}

// Load
func (k *kube) Load() ([]*source.KeyValue, error) {
	if k.op.Namespace == "" {
		return nil, errors.New("SourceOption namespace not full")
	}
	if err := k.initKubeClient(); err != nil {
		return nil, err
	}
	return k.loadKubeConfig()
}

// Watch
func (k *kube) Watch() (source.Watcher, error) {
	return newWatcher(k)
}
