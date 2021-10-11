package framework

import (
	"errors"
	"sync"
)

type Container interface {
	Bind(provider ServiceProvider) error

	IsBind(key string) bool

	Make(key string) (interface{}, error)

	MustMake(key string) interface{}

	MakeNew(key string, params []interface{}) (interface{}, error)
}

type GinxContainer struct {
	Container

	providers map[string]ServiceProvider
	instances map[string]interface{}

	lock sync.RWMutex
}

func NewGinxContainer() *GinxContainer {
	return &GinxContainer{
		providers: map[string]ServiceProvider{},
		instances: map[string]interface{}{},
		lock:      sync.RWMutex{},
	}
}

func (ginx *GinxContainer) PrintProviders() []string {
	ret := make([]string, 0, len(ginx.providers))
	for _, provider := range ginx.providers {
		name := provider.Name()
		ret = append(ret, name)
	}
	return ret
}

func (ginx *GinxContainer) Bind(provider ServiceProvider) error {
	ginx.lock.Lock()
	defer ginx.lock.Unlock()
	key := provider.Name()

	ginx.providers[key] = provider

	if !provider.IsDefer() {
		if err := provider.Boot(ginx); err != nil {
			return err
		}
		params := provider.Params(ginx)
		method := provider.Register(ginx)
		instance, err := method(params...)
		if err != nil {
			return errors.New(err.Error())
		}
		ginx.instances[key] = instance
	}
	return nil
}

func (ginx *GinxContainer) IsBind(key string) bool {
	return ginx.findServiceProvider(key) != nil
}

func (ginx *GinxContainer) findServiceProvider(key string) ServiceProvider {
	ginx.lock.RLock()
	defer ginx.lock.RUnlock()
	if sp, ok := ginx.providers[key]; ok {
		return sp
	}
	return nil
}

func (ginx *GinxContainer) Make(key string) (interface{}, error) {
	return ginx.make(key, nil, false)
}

func (ginx *GinxContainer) MustMake(key string) interface{} {
	serv, err := ginx.make(key, nil, false)
	if err != nil {
		panic(err)
	}
	return serv
}

func (ginx *GinxContainer) MakeNew(key string, params []interface{}) (interface{}, error) {
	return ginx.make(key, params, true)
}

func (ginx *GinxContainer) make(key string, params []interface{}, forceNew bool) (interface{}, error) {
	ginx.lock.RLock()
	defer ginx.lock.RUnlock()
	sp := ginx.findServiceProvider(key)
	if sp == nil {
		return nil, errors.New("contract " + key + " have not register")
	}

	if forceNew {
		ginx.newInstance(sp, params)
	}

	if ins, ok := ginx.instances[key]; ok {
		return ins, nil
	}
	ins, err := ginx.newInstance(sp, params)
	if err != nil {
		return nil, err
	}
	ginx.instances[key] = ins
	return ins, nil
}

func (ginx *GinxContainer) newInstance(sp ServiceProvider, params []interface{}) (interface{}, error) {
	if err := sp.Boot(ginx); err != nil {
		return nil, err
	}
	if params == nil {
		params = sp.Params(ginx)
	}
	method := sp.Register(ginx)
	ins, err := method(params...)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return ins, nil
}
