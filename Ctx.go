package dogo

import (
	"reflect"
	"sync"
	"errors"
)

type BluePrintField struct {
	Name string //field name
	ValueType int //value type
	Value interface{}
}
type Blueprint struct {
	Type reflect.Type
	TypeStr string
	Fields []*BluePrintField
}

type Ctx struct {
	typeInstance map[reflect.Type]interface{} //type to instance map
	idInstance map[string]interface{} //id to instance map
	idBlueprint map[string]*Blueprint //id to type definition map
	typeBlueprint map[reflect.Type]*Blueprint //type to type definition map
	types map[string]reflect.Type //typeid to reflect type map
	mutex *sync.Mutex
}

func(ctx *Ctx)GetInstanceWithId(id string) (interface{}, error) {
	var(
		ins interface{}
		bp *Blueprint
		err error
		exist bool
	)
	if ins, exist = ctx.idInstance[id];!exist {
		defer ctx.mutex.Unlock()
		ctx.mutex.Lock()
		if ins, exist = ctx.idInstance[id]; exist {
			return ins, nil
		}
		if bp, exist= ctx.idBlueprint[id]; !exist {
			return nil, errors.New("No definition found for " + id)
		}
		ins, err = ctx.buildInstance(bp)
		if err != nil {
			return nil ,err
		}
		ctx.idInstance[id] = ins
		ctx.typeInstance[bp.Type] = ins
	}
	return ins, nil
}

func(ctx *Ctx)GetInstanceWithType(t reflect.Type)(interface{}, error) {
	return nil, nil
}

func(ctx *Ctx)NewInstanceWithId(id string)(interface{}, error) {
	return nil, nil
}

func(ctx *Ctx)NewInstanceWithType(id string)(interface{}, error) {
	return nil, nil
}

func(ctx *Ctx)buildInstance(bp *Blueprint) (interface{}, error) {
	return nil, nil
}

func(ctx *Ctx)RegBlueprint(id string, bp *Blueprint) {
	tStr := bp.TypeStr
	t, typeExist := ctx.types[tStr];
	if !typeExist {
		panic("type " + tStr + "not exist")
	}
	bp.Type = t
	if _, exist := ctx.typeBlueprint[t]; !exist {
		ctx.typeBlueprint[t] = bp
	}

	if _, exist := ctx.idBlueprint[id]; !exist {
		ctx.idBlueprint[id] = bp
	}

}

func(ctx *Ctx)RegType(t interface{}) {
	tType := reflect.TypeOf(t)
	if tType.Kind() != reflect.Struct {
		panic("Only struct can be registerd")
	}
	fullName := tType.PkgPath() + "/" + tType.Name()
	ctx.types[fullName] = tType
}