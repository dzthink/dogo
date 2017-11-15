package dogo

import (
	"reflect"
	"sync"
)

const (
	ValueTypeConst  = iota
	ValueTypeConfig
	ValueTypeRef
	ValueTypeAutoWired
)
type BluePrintField struct {
	Name string //field name
	ValueType int //value type
	Value interface{}
}
type Blueprint struct {
	Type reflect.Type
	TypeStr string
	Fields map[string]*BluePrintField
}

type Ctx struct {
	typeInstance map[reflect.Type]interface{} //type to instance map
	idInstance map[string]interface{} //id to instance map
	idBlueprint map[string]*Blueprint //id to type definition map
	typeBlueprint map[reflect.Type]*Blueprint //type to type definition map
	types map[string]reflect.Type //typeid to reflect type map
	mutex *sync.Mutex
}

func(ctx *Ctx)GetInstanceWithId(id string) (interface{}) {
	var(
		ins interface{}
		bp *Blueprint
		exist bool
	)
	if ins, exist = ctx.idInstance[id];!exist {
		defer ctx.mutex.Unlock()
		ctx.mutex.Lock()
		if ins, exist = ctx.idInstance[id]; exist {
			return ins
		}
		if bp, exist= ctx.idBlueprint[id]; !exist {
			return nil
		}
		ins = ctx.buildInstance(bp)

		ctx.idInstance[id] = ins
		ctx.typeInstance[bp.Type] = ins
	}
	return ins
}

func(ctx *Ctx)GetInstanceWithType(t reflect.Type)(interface{}) {
	var(
		ins interface{}
		bp *Blueprint
		exist bool
	)
	if ins, exist = ctx.typeInstance[t];!exist {
		defer ctx.mutex.Unlock()
		ctx.mutex.Lock()
		if ins, exist = ctx.typeInstance[t]; exist {
			return ins
		}
		if bp, exist= ctx.typeBlueprint[t]; !exist {
			bp = &Blueprint{
				Type : t,
				TypeStr : "",
				Fields:make(map[string]*BluePrintField),
			}
		}
		ins = ctx.buildInstance(bp)

		ctx.typeInstance[bp.Type] = ins
	}
	return ins
}

func(ctx *Ctx)NewInstanceWithId(id string)(interface{}) {
	var(
		bp *Blueprint
		exist bool
	)
	defer ctx.mutex.Unlock()
	ctx.mutex.Lock()
	if bp, exist= ctx.idBlueprint[id]; !exist {
		return nil
	}
	return ctx.buildInstance(bp)
}

func(ctx *Ctx)NewInstanceWithType(t reflect.Type)(interface{}) {
	var(
		bp *Blueprint
		exist bool
	)
	defer ctx.mutex.Unlock()
	ctx.mutex.Lock()

	if bp, exist= ctx.typeBlueprint[t]; !exist {
		bp = &Blueprint{
			Type : t,
			TypeStr : "",
			Fields:make(map[string]*BluePrintField),
		}
	}

	return ctx.buildInstance(bp)

}

func(ctx *Ctx)mergeBlueprintField(t reflect.Type, fields map[string]*BluePrintField) {

}

func(ctx *Ctx)initField(fieldValue reflect.Value) {
	//Init map/slice
	switch fieldValue.Kind() {
	case reflect.Map:
		if !fieldValue.IsValid() || fieldValue.IsNil() {
			fieldValue.Set(reflect.MakeMap(fieldValue.Type()))
		}
	case reflect.Slice:
		if !fieldValue.IsValid() || fieldValue.IsNil() {
			fieldValue.Set(reflect.MakeSlice(fieldValue.Type(), 0, 0))
		}
	}
}

func(ctx *Ctx)injectField(fieldValue reflect.Value, bpField *BluePrintField) {
	switch bpField.ValueType {
	case ValueTypeConst:
	case ValueTypeConfig:
	case ValueTypeRef:
	case ValueTypeAutoWired:
	}
}

func(ctx *Ctx)buildInstance(bp *Blueprint) (interface{}) {
	t := bp.Type
	ctx.mergeBlueprintField(t, bp.Fields)
	ins := reflect.New(t)
	if initType, exist := t.MethodByName("Init");exist && initType.Type.NumIn() == 0 {
		ins.MethodByName("Init").Call(make([]reflect.Value, 0))
	}
	for index := 0; index < t.NumField(); index++ {
		ctx.initField(ins.Field(index).Elem())
		if bpField, ok := bp.Fields[t.Field(index).Name]; ok {
			ctx.injectField(ins.Field(index).Elem(), bpField)
		}
	}
	return ins
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