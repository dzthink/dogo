package dogo

import (
	"reflect"
	"sync"
	"strings"
)

const (
	ValueTypeConst  = iota
	ValueTypeConfig
	ValueTypeRef
	ValueTypeAutoWired
)

type TypeMeta struct {
	Alias string
	Abstract reflect.Type
	Implement reflect.Type
}

type BluePrintField struct {
	Name string //field name
	ValueType int //value type
	Value interface{}
}

type Blueprint struct {
	Type reflect.Type
	TypeAlias string
	Fields map[string]*BluePrintField
}

type Ctx struct {
	typeInstance map[reflect.Type]interface{} //type to instance map
	idInstance map[string]interface{} //id to instance map
	idBlueprint map[string]*Blueprint //id to type definition map
	typeBlueprint map[reflect.Type]*Blueprint //type to type definition map
	typesMeta []*TypeMeta
	mutex *sync.Mutex
}

func NewCtx() *Ctx {
	return &Ctx{
		typeInstance:make(map[reflect.Type]interface{}),
		idInstance:make(map[string]interface{}),
		idBlueprint:make(map[string]*Blueprint),
		typeBlueprint:make(map[reflect.Type]*Blueprint),
		typesMeta : make([]*TypeMeta, 0 ,50),
		mutex:new(sync.Mutex),
	}
}

func(ctx *Ctx)GetInstanceWithId(id string) (interface{}) {
	var(
		ins interface{}
		exist bool
	)
	if ins, exist = ctx.idInstance[id];!exist {
		defer ctx.mutex.Unlock()
		ctx.mutex.Lock()
		ins = ctx.buildInstanceWithId(id, true)
	}
	return ins
}

func(ctx *Ctx)GetInstanceWithType(t reflect.Type)(interface{}) {
	var(
		ins interface{}
		exist bool
	)
	if ins, exist = ctx.typeInstance[t];!exist {
		defer ctx.mutex.Unlock()
		ctx.mutex.Lock()
		ins = ctx.buildInstanceWithType(t, true)
	}
	return ins
}

func(ctx *Ctx)NewInstanceWithId(id string)(interface{}) {
	defer ctx.mutex.Unlock()
	ctx.mutex.Lock()
	return ctx.buildInstanceWithId(id, false)
}

func(ctx *Ctx)NewInstanceWithType(t reflect.Type)(interface{}) {

	defer ctx.mutex.Unlock()
	ctx.mutex.Lock()

	return ctx.buildInstanceWithType(t, false)
}

func(ctx *Ctx)mergeBlueprintField(t reflect.Type, fields map[string]*BluePrintField) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	for index := 0; index < t.NumField(); index++ {
		fieldStruct := t.Field(index)
		if _, ok := fields[fieldStruct.Name];ok {
			continue
		}
		tag := t.Field(index).Tag
		if tagValue :=tag.Get("Ref");!strings.EqualFold(tagValue, "") {
			fields[fieldStruct.Name] = &BluePrintField{
				Name : fieldStruct.Name,
				ValueType : ValueTypeRef,
				Value : tagValue,
			}
			continue
		}

		if tagValue :=tag.Get("Value");!strings.EqualFold(tagValue, "") {
			fields[fieldStruct.Name] = &BluePrintField{
				Name : fieldStruct.Name,
				ValueType : ValueTypeConst,
				Value : tagValue,
			}
			continue
		}

		if tagValue :=tag.Get("Config");!strings.EqualFold(tagValue, "") {
			fields[fieldStruct.Name] = &BluePrintField{
				Name : fieldStruct.Name,
				ValueType : ValueTypeConfig,
				Value : tagValue,
			}
			continue
		}

		if tagValue :=tag.Get("Autowired");!strings.EqualFold(tagValue, "") {
			fields[fieldStruct.Name] = &BluePrintField{
				Name : fieldStruct.Name,
				ValueType : ValueTypeAutoWired,
				Value : nil,
			}
			continue
		}
	}
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
	switch fieldValue.Kind() {

	//1. 字面值情况,类型可能为所有普通类型
	//2. 不支持配置类型
	//3. 引用类型, 类型为字符串
	//4. 不支持autowired
	case reflect.Slice:
		for _, e := range bpField.Value.([]interface{}) {
			if bpField.ValueType == ValueTypeRef {
				fieldValue.Set(reflect.Append(fieldValue, reflect.ValueOf(ctx.buildInstanceWithId(e.(string), false))))
			} else {
				fieldValue.Set(reflect.ValueOf(e))
			}
		}
		break
		//1. 字面值情况,类型可能为所有普通类型
		//2. 不支持配置类型
		//3. 引用类型, 类型为字符串
		//4. 不支持autowired
	case reflect.Map:
		for k, v := range bpField.Value.(map[string]interface{}) {
			if bpField.ValueType == ValueTypeRef {
				fieldValue.SetMapIndex(reflect.ValueOf(k),
					reflect.ValueOf(ctx.buildInstanceWithId(v.(string), false)))
			} else {
				fieldValue.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
			}
		}
		break
		//1. 字面值不支持
		//2. 配置不支持
		//3. 引用类型, 类型为字符串
		//4. autowired
	case reflect.Struct:
		if bpField.ValueType == ValueTypeRef {
			fieldValue.Set(reflect.ValueOf(ctx.buildInstanceWithId(bpField.Value.(string), false)).Elem())
		} else if bpField.ValueType == ValueTypeAutoWired {
			fieldValue.Set(reflect.ValueOf(ctx.buildInstanceWithType(fieldValue.Type(), false)).Elem())
		}
		break
	case reflect.Ptr:
	case reflect.Interface:
		if bpField.ValueType == ValueTypeRef {
			fieldValue.Set(reflect.ValueOf(ctx.buildInstanceWithId(bpField.Value.(string), false)))
		} else if bpField.ValueType == ValueTypeAutoWired {
			fieldValue.Set(reflect.ValueOf(ctx.buildInstanceWithType(fieldValue.Type(), false)))
		}
		break
	default:
		if bpField.ValueType == ValueTypeConst {
			fieldValue.Set(reflect.ValueOf(bpField.Value).Convert(fieldValue.Type()))
		} else if bpField.ValueType == ValueTypeConfig {
			//todo decide to implement the feature or not
			// on one hand it's convenient to inject config value to component, on the other hand it will make the di
			//system couple a configuration manager which is bad
		}
	}


}

func(ctx *Ctx)buildInstanceWithId(id string, save bool) (interface{}) {
	var(
		ins interface{}
		exist bool
		bp *Blueprint
	)
	if ins, exist = ctx.idInstance[id]; exist {
		return ins
	}
	if bp, exist= ctx.idBlueprint[id]; !exist {
		return nil
	}
	ins = ctx.buildInstance(bp)
	if save {
		ctx.idInstance[id] = ins
		ctx.typeInstance[bp.Type] = ins
	}
	return ins
}

func(ctx *Ctx)buildInstanceWithType(t reflect.Type, save bool) (interface{}) {
	var(
		ins interface{}
		exist bool
		bp *Blueprint
	)
	if ins, exist = ctx.typeInstance[t]; exist {
		return ins
	}
	if bp, exist= ctx.typeBlueprint[t]; !exist {
		typeMeta := ctx.searchTypeByType(t)
		if typeMeta == nil {
			return nil
		}
		bp = &Blueprint{
			Type : typeMeta.Implement,
			TypeAlias : typeMeta.Alias,
			Fields:make(map[string]*BluePrintField),
		}
	}
	ins = ctx.buildInstance(bp)

	if save {
		ctx.typeInstance[bp.Type] = ins
	}
	return ins
}
func(ctx *Ctx)buildInstance(bp *Blueprint) (interface{}) {
	t := bp.Type
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	ctx.mergeBlueprintField(t, bp.Fields)
	ins := reflect.New(t)
	if initType, exist := t.MethodByName("Init");exist && initType.Type.NumIn() == 0 {
		ins.MethodByName("Init").Call(make([]reflect.Value, 0))
	}
	for index := 0; index < t.NumField(); index++ {
		if bpField, ok := bp.Fields[t.Field(index).Name]; ok {
			fieldValue := ins.Elem().FieldByName(t.Field(index).Name)
			ctx.injectField(fieldValue, bpField)
		}
	}
	return ins.Interface()
}

func(ctx *Ctx)RegBlueprint(id string, bp *Blueprint) {
	typeAlias := ctx.searchTypeByAlias(bp.TypeAlias)
	if typeAlias == nil {
		panic("type " + bp.TypeAlias + "not exist")
	}
	bp.Type = typeAlias.Implement
	if _, exist := ctx.typeBlueprint[typeAlias.Abstract]; !exist {
		ctx.typeBlueprint[typeAlias.Abstract] = bp
	}

	if _, exist := ctx.idBlueprint[id]; !exist {
		ctx.idBlueprint[id] = bp
	}

}

func(ctx *Ctx)RegType(meta *TypeMeta) {
	ctx.typesMeta = append(ctx.typesMeta, meta)
}

func(ctx *Ctx)RegTypes(metas []*TypeMeta) {
	ctx.typesMeta = append(ctx.typesMeta, metas...)
}

func(ctx *Ctx)searchTypeByAlias(alias string)(*TypeMeta) {
	for _, meta := range ctx.typesMeta {
		if strings.EqualFold(meta.Alias, alias) {
			return meta
		}
	}
	return nil
}

func(ctx *Ctx)searchTypeByType(t reflect.Type)(*TypeMeta) {
	for _, meta := range ctx.typesMeta {
		if meta.Abstract == t {
			return meta
		}
	}
	return nil
}