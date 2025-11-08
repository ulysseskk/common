package thirdparty

import (
	"fmt"
	"sync"
)

const (
	NamespaceDefault = "default"
)

func Register[T any, C any](factoryMethod FactoryMethod[T, C]) *thirdPartyFactory[T, C] {
	return &thirdPartyFactory[T, C]{
		factoryMethod: factoryMethod,
		instances:     map[string]T{},
		lock:          sync.RWMutex{},
	}
}

type FactoryMethod[T any, C any] func(cfg C) (T, error)

type thirdPartyFactory[T any, C any] struct {
	factoryMethod FactoryMethod[T, C]
	singleton     bool
	instances     map[string]T
	lock          sync.RWMutex
}

func (t *thirdPartyFactory[T, C]) New(namespace string, cfg C) (result T, err error) {
	t.lock.Lock()
	if _, ok := t.instances[namespace]; ok {
		t.lock.Unlock()
		return result, fmt.Errorf("instance already exist")
	}
	t.lock.Unlock()
	result, err = t.factoryMethod(cfg)
	if err != nil {
		return result, err
	}
	t.lock.Lock()
	defer t.lock.Unlock()
	t.instances[namespace] = result
	return result, nil
}

func (t *thirdPartyFactory[T, C]) NewDefault(cfg C) (T, error) {
	return t.New(NamespaceDefault, cfg)
}

func (t *thirdPartyFactory[T, C]) Get(namespace string) T {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.instances[namespace]
}

func (t *thirdPartyFactory[T, C]) GetDefault() T {
	return t.Get(NamespaceDefault)
}
