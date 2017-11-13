package dogo

import "reflect"

type typeDef struct {
	t reflect.Type

}

type Ctx struct {
	typeInstance map[reflect.Type]interface{}
	idInstance map[string]interface{}
}

func(ctx *Ctx)GetInstanceWithId(id string) (interface{}, error) {
	return nil, nil
}

func(ctx *Ctx)GetInstanceWithType(t reflect.Type)(interface{}, error) {
	return nil, nil
}