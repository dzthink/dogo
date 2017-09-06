//依赖管理容器
package dogo

import (
	"reflect"
	"sync"
	"errors"
	"strings"
)

const(
	SCOPE_STATEFUL = "stateful"
	SCOPE_STATELESS = "stateless"
)

type dependDefinition struct {
	isRef bool
	dataType string
	data interface{}
}

type instanceDefinition struct {
	scope string
	proto reflect.Type
	dependsOn map[string]*dependDefinition
	isLazy bool
}

type core struct {
	typeMap map[string]reflect.Type
	dependencies map[string]*instanceDefinition
	instances map[string]interface{}
	mutex *sync.Mutex
}

func NewCore() *core {
	return nil
}


func(c *core)make(name string) (interface{}, error) {
	var (
		def *instanceDefinition
		ins interface{}
		ok bool
	)
	if def, ok = c.dependencies[name];!ok {
		return nil, errors.New(name + " definition not found");
	}

	if ins, ok = c.instances[name];!ok {

	}

	//todo 处理过程中需要修改singleton， 涉及到锁处理
	return ins, nil;
}

func(c *core)makeInstance(id *instanceDefinition) (interface{}, error) {
	return nil, nil
}
