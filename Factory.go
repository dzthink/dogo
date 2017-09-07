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
	data interface{}
}

type instanceDefinition struct {
	scope string
	proto reflect.Type
	dependsOn map[string]*dependDefinition
}

func NewInstanceDefinition(scope string, proto reflect.Type) *instanceDefinition {
	if !strings.EqualFold(scope, SCOPE_STATELESS) && !strings.EqualFold(scope, SCOPE_STATEFUL) {
		panic("scope should be either 'stateful' or 'stateless'")
	}
	return &instanceDefinition{
		scope : scope,
		proto : proto,
		dependsOn : make(map[string]*dependDefinition),
	}
}


type Factory struct {
	typeMap map[string]reflect.Type
	dependencies map[string]*instanceDefinition
	instances map[string]interface{}
	mutex *sync.Mutex
}

func NewFactory() *Factory {
	return &Factory{
		typeMap : make(map[string]reflect.Type),
		dependencies : make(map[string]*instanceDefinition),
		instances : make(map[string]interface{}),
		mutex : new(sync.Mutex),
	}
}


func(c *Factory)Make(name string) (interface{}, error) {
	var (
		curDef *instanceDefinition
		ins interface{}
		ok bool
	)

	if ins, ok = c.instances[name];!ok {
		if curDef, ok = c.dependencies[name];!ok {
			return nil, errors.New(name + " definition not found");
		}
		insValue := reflect.New(curDef.proto)
		for k,v := range curDef.dependsOn {
			fieldValue := insValue.Elem().FieldByName(k)
			if c.injectField(fieldValue, v) != nil {
				return nil, errors.New("can't get instance for dependency:" + name);
			}
		}

		ins = insValue.Interface()
		if strings.EqualFold(curDef.scope, SCOPE_STATELESS) {
			ins = c.saveInstance(name, ins)
		}
	}

	return ins, nil;
}

func(c *Factory)saveInstance(name string, ins interface{}) (interface{}) {
	if instance, ok := c.instances[name];!ok {
		defer c.mutex.Unlock()
		c.mutex.Lock()
		c.instances[name] = ins
		return ins
	} else {
		return instance
	}
}


func(c *Factory)injectField(fieldValue reflect.Value, depend *dependDefinition) error {
	if !depend.isRef {
		//todo 这里的值应该来自properties配置读取
		fieldValue.Set(reflect.ValueOf(depend.data).Convert(fieldValue.Type()))
		return nil
	}

	dataValue := reflect.ValueOf(depend.data)
	switch dataValue.Kind() {
	case reflect.Slice:
		dependData := depend.data.([]string)
		insFieldData := reflect.MakeSlice(fieldValue.Type(), 0, len(dependData))
		for _, dependId := range depend.data.([]string) {
			dependIns, err := c.Make(dependId)
			if err != nil {
				return err
			}
			insFieldData = reflect.Append(insFieldData, reflect.ValueOf(dependIns))
		}
		fieldValue.Set(insFieldData)
		break;
	case reflect.Map:
		fieldValue.Set(reflect.MakeMap(fieldValue.Type()))
		for k, v := range depend.data.(map[string]string) {
			dependIns, err := c.Make(v)
			if err != nil {
				return err
			}
			fieldValue.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(dependIns))
		}
		break;
	case reflect.String:
		dependIns, err := c.Make(depend.data.(string))
		if err != nil {
			return err
		}
		if fieldValue.Kind() == reflect.Ptr {
			fieldValue.Set(reflect.ValueOf(dependIns))
		} else {
			fieldValue.Set(reflect.ValueOf(dependIns).Elem())
		}

		break;
	default:
		panic("type of data Field of A ref dependDefiniton's shoud be one of (slice, map, string)")
	}
	return nil
}

func(f *Factory) RegisType(t interface{}, name, scope string)   {
	tType := reflect.TypeOf(t)
	if tType.Kind() != reflect.Struct {
		panic("Only struct can be registerd")
	}
	if strings.EqualFold(name, "") {
		panic("type name should not be empty")
	}
	f.typeMap[name] = tType

	insDef := NewInstanceDefinition(scope, tType)
	f.parseStructTag(tType, insDef)
	f.RegisDefinition(name, insDef)
}

func(f *Factory) RegisDefinition(id string, def *instanceDefinition) {
	f.dependencies[id] = def
}

func(f *Factory) parseStructTag(t reflect.Type, insDef *instanceDefinition) {
	for index := 0; index < t.NumField(); index++ {
		tag := t.Field(index).Tag
		if tagName, ok := tag.Lookup("Autowired"); ok {
			insDef.dependsOn[t.Field(index).Name] = &dependDefinition{
				isRef : true,
				data : tagName,
			}
			continue
		}

		if tagName, ok := tag.Lookup("Value"); ok {
			insDef.dependsOn[t.Field(index).Name] = &dependDefinition{
				isRef : false,
				data : tagName,
			}
		}
	}
}
