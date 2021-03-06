package dogo

import (
	"reflect"
	"sync"
	"strings"
	"context"
)

const (
	ValueTypeConst  = "const"
	ValueTypeConfig = "config"
	ValueTypeRef = "ref"
	ValueTypeAutoWired = "autowired"
)

type TypeMeta struct {
	Alias string
	Abstract reflect.Type
	Implement reflect.Type
}

type bluePrintField struct {
	Name string //field name
	ValueType string //value type
	Value interface{}
}

type blueprint struct {
	InitMethod string
	Type reflect.Type
	TypeAlias string
	Fields map[string]*bluePrintField
}

type Ioc struct {
	typeInstance map[reflect.Type]interface{} //type to instance map
	idInstance map[string]interface{} //id to instance map
	idBlueprint map[string]*blueprint //id to type definition map
	typeBlueprint map[reflect.Type]*blueprint //type to type definition map
	typesMeta []*TypeMeta
	mutex *sync.Mutex
}

func newIoc(ts []*TypeMeta) *Ioc {
	ctx := &Ioc{
		typeInstance:make(map[reflect.Type]interface{}),
		idInstance:make(map[string]interface{}),
		idBlueprint:make(map[string]*blueprint),
		typeBlueprint:make(map[reflect.Type]*blueprint),
		typesMeta : make([]*TypeMeta, 0 ,50),
		mutex:new(sync.Mutex),
	}
	ctx.regTypes(ts)
	return ctx
}

func(ioc *Ioc)GetInstanceWithId(id string) (interface{}) {
	var(
		ins interface{}
		exist bool
	)
	if ins, exist = ioc.idInstance[id];!exist {
		defer ioc.mutex.Unlock()
		ioc.mutex.Lock()
		ins = ioc.buildInstanceWithId(id, true)
	}
	return ins
}

func(ioc *Ioc)GetInstanceWithType(t reflect.Type)(interface{}) {
	var(
		ins interface{}
		exist bool
	)
	if ins, exist = ioc.typeInstance[t];!exist {
		defer ioc.mutex.Unlock()
		ioc.mutex.Lock()
		ins = ioc.buildInstanceWithType(t, true)
	}
	return ins
}

func(ioc *Ioc)NewInstanceWithId(id string)(interface{}) {
	defer ioc.mutex.Unlock()
	ioc.mutex.Lock()
	return ioc.buildInstanceWithId(id, false)
}

func(ioc *Ioc)NewInstanceWithType(t reflect.Type)(interface{}) {

	defer ioc.mutex.Unlock()
	ioc.mutex.Lock()

	return ioc.buildInstanceWithType(t, false)
}

func(ioc *Ioc)mergeBlueprintField(t reflect.Type, fields map[string]*bluePrintField) {
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
			fields[fieldStruct.Name] = &bluePrintField{
				Name : fieldStruct.Name,
				ValueType : ValueTypeRef,
				Value : tagValue,
			}
			continue
		}

		if tagValue :=tag.Get("Value");!strings.EqualFold(tagValue, "") {
			fields[fieldStruct.Name] = &bluePrintField{
				Name : fieldStruct.Name,
				ValueType : ValueTypeConst,
				Value : tagValue,
			}
			continue
		}

		if tagValue :=tag.Get("Config");!strings.EqualFold(tagValue, "") {
			fields[fieldStruct.Name] = &bluePrintField{
				Name : fieldStruct.Name,
				ValueType : ValueTypeConfig,
				Value : tagValue,
			}
			continue
		}

		if tagValue :=tag.Get("Autowired");!strings.EqualFold(tagValue, "") {
			fields[fieldStruct.Name] = &bluePrintField{
				Name : fieldStruct.Name,
				ValueType : ValueTypeAutoWired,
				Value : nil,
			}
			continue
		}
	}
}

func(ioc *Ioc)initField(fieldValue reflect.Value) {
	//Init map/slice
	switch fieldValue.Kind() {
	case reflect.Map:

		if !fieldValue.IsValid() || fieldValue.IsNil() {
			fieldValue.Set(reflect.MakeMap(fieldValue.Type()))
		}
		break;
	case reflect.Slice:
		if !fieldValue.IsValid() || fieldValue.IsNil() {
			fieldValue.Set(reflect.MakeSlice(fieldValue.Type(), 0, 0))
		}
		break;
	}
}

func(ioc *Ioc)injectField(fieldValue reflect.Value, bpField *bluePrintField) {
	switch fieldValue.Kind() {

	//1. 字面值情况,类型可能为所有普通类型
	//2. 不支持配置类型
	//3. 引用类型, 类型为字符串
	//4. 不支持autowired
	case reflect.Slice:
		for _, e := range bpField.Value.([]interface{}) {
			if bpField.ValueType == ValueTypeRef {
				fieldValue.Set(reflect.Append(fieldValue, reflect.ValueOf(ioc.buildInstanceWithId(e.(string), false))))
			} else {
				fieldValue.Set(reflect.Append(fieldValue, reflect.ValueOf(e)))
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
					reflect.ValueOf(ioc.buildInstanceWithId(v.(string), false)))
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
			fieldValue.Set(reflect.ValueOf(ioc.buildInstanceWithId(bpField.Value.(string), false)).Elem())
		} else if bpField.ValueType == ValueTypeAutoWired {
			fieldValue.Set(reflect.ValueOf(ioc.buildInstanceWithType(fieldValue.Type(), false)).Elem())
		}
		break
	case reflect.Ptr:
		if bpField.ValueType == ValueTypeRef {
			fieldValue.Set(reflect.ValueOf(ioc.buildInstanceWithId(bpField.Value.(string), false)))
		} else if bpField.ValueType == ValueTypeAutoWired {
			fieldValue.Set(reflect.ValueOf(ioc.buildInstanceWithType(fieldValue.Type().Elem(), false)))
		}
		break
	case reflect.Interface:
		if bpField.ValueType == ValueTypeRef {
			fieldValue.Set(reflect.ValueOf(ioc.buildInstanceWithId(bpField.Value.(string), false)))
		} else if bpField.ValueType == ValueTypeAutoWired {
			fieldValue.Set(reflect.ValueOf(ioc.buildInstanceWithType(fieldValue.Type(), false)))
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

func(ioc *Ioc)buildInstanceWithId(id string, save bool) (interface{}) {
	var(
		ins interface{}
		exist bool
		bp *blueprint
	)
	if ins, exist = ioc.idInstance[id]; exist {
		return ins
	}
	if bp, exist= ioc.idBlueprint[id]; !exist {
		return nil
	}
	ins = ioc.buildInstance(bp)
	if save {
		ioc.idInstance[id] = ins
		ioc.typeInstance[bp.Type] = ins
	}
	return ins
}

func(ioc *Ioc)buildInstanceWithType(t reflect.Type, save bool) (interface{}) {
	var(
		ins interface{}
		exist bool
		bp *blueprint
	)
	if ins, exist = ioc.typeInstance[t]; exist {
		return ins
	}
	if bp, exist= ioc.typeBlueprint[t]; !exist {
		typeMeta := ioc.searchTypeByType(t)
		if typeMeta == nil {
			return nil
		}
		bp = &blueprint{
			Type : typeMeta.Implement,
			TypeAlias : typeMeta.Alias,
			Fields:make(map[string]*bluePrintField),
		}
	}
	ins = ioc.buildInstance(bp)

	if save {
		ioc.typeInstance[bp.Type] = ins
	}
	return ins
}
func(ioc *Ioc)buildInstance(bp *blueprint) (interface{}) {
	t := bp.Type
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	ioc.mergeBlueprintField(t, bp.Fields)
	ins := reflect.New(t)
	//todo preinit callback
	for index := 0; index < t.NumField(); index++ {
		if bpField, ok := bp.Fields[t.Field(index).Name]; ok {
			fieldValue := ins.Elem().FieldByName(t.Field(index).Name)
			ioc.initField(fieldValue)
			ioc.injectField(fieldValue, bpField)
		}
	}
	return ins.Interface()
}

func(ioc *Ioc)parseBluePrint(confs []*Config) {
	for _, bpConf := range confs {
		bpConfId ,err := bpConf.String(CONF_CTX_ID)
		if err != nil {
			panic("context config error, id missing, config string " + bpConf.ToString())
		}
		bpConfId = strings.TrimSpace(bpConfId)
		if _,exist := ioc.idBlueprint[bpConfId]; exist {
			panic("context config error, id duplicate, config string" + bpConf.ToString())
		}
		bpConfAlias , err := bpConf.String(CONF_CTX_ALIAS)
		if err != nil {
			panic("context config error, alias missing, conf string:" + bpConf.ToString())
		}
		tpMeta := ioc.searchTypeByAlias(bpConfAlias)
		if tpMeta == nil {
			panic("context config error alias not exist , conf string:" + bpConf.ToString())
		}
		bpInitMethod, _ := bpConf.String(CONF_CTX_INIT_METHOD)
		bp := &blueprint{
			Type : tpMeta.Implement,
			InitMethod:bpInitMethod,
			TypeAlias : bpConfAlias,
			Fields:make(map[string]*bluePrintField),
		}
		bpFieldConfs, err := bpConf.ChildList(CONF_CTX_FIELDS)
		ioc.typeBlueprint[tpMeta.Abstract] = bp
		ioc.idBlueprint[bpConfId] = bp
		if err != nil {
			continue
		}
		for _, bpFieldConf := range bpFieldConfs {
			var err error
			bpField := &bluePrintField{}
			bpField.Name, err = bpFieldConf.String(CONF_CTX_FIELD_NAME)
			bpField.ValueType, err = bpFieldConf.String(CONF_CTX_FIELD_TYPE)
			if err != nil {
				continue
			}
			fieldSt, exist := tpMeta.Implement.FieldByName(bpField.Name)
			if !exist {
				continue
			}

			switch fieldSt.Type.Kind() {
			case reflect.Slice:
				var r []interface{}
				bpFieldConf.Get(CONF_CTX_FIELD_VALUE, &r)
				bpField.Value = r
			case reflect.Map:
				var r map[string]interface{}
				bpFieldConf.Get(CONF_CTX_FIELD_VALUE, &r)
				bpField.Value = r
			default:
				bpFieldConf.Get(CONF_CTX_FIELD_VALUE, &bpField.Value)
			}
			bp.Fields[bpField.Name] = bpField
		}
	}
}


func(ioc *Ioc)active(ctx context.Context) {
	param := []reflect.Value{
		reflect.ValueOf(ctx),
	}
	for id, bp := range ioc.idBlueprint{
		ins := ioc.GetInstanceWithId(id)
		if ins == nil {
			panic("error to get instance of " + id)
		}
		insValue := reflect.ValueOf(ins)
		if strings.EqualFold(bp.InitMethod, "") {
			continue
		}
		if initType, exist := insValue.Type().MethodByName(bp.InitMethod);exist && initType.Type.NumIn() == 2 {
			insValue.MethodByName(bp.InitMethod).Call(param)
		}
	}
}

func(ioc *Ioc)regTypes(metas []*TypeMeta) {
	ioc.typesMeta = append(ioc.typesMeta, metas...)
}

func(ioc *Ioc)searchTypeByAlias(alias string)(*TypeMeta) {
	for _, meta := range ioc.typesMeta {
		if strings.EqualFold(meta.Alias, alias) {
			return meta
		}
	}
	return nil
}

func(ioc *Ioc)searchTypeByType(t reflect.Type)(*TypeMeta) {
	for _, meta := range ioc.typesMeta {
		if meta.Abstract == t {
			return meta
		}
	}
	return nil
}