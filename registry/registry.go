package registry

import "github.com/go-kratos/kratos/v2/registry"

type kube struct {
}

// NewRegistry new a kube registry.
func NewRegistry() registry.Registry {
	return &kube{}
}

// Register the registration.
func (k *kube) Register(service *registry.ServiceInstance) error {
	return nil
}

// Deregister the registration.
func (k *kube) Deregister(service *registry.ServiceInstance) error {
	return nil
}

// Service return the service instances in memory according to the service name.
func (k *kube) Service(name string) ([]*registry.ServiceInstance, error) {
	return nil, nil
}

// Watch creates a watcher according to the service name.
func (k *kube) Watch(name string) (registry.Watcher, error) {
	return newWatcher(k)
}
