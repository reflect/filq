package types

import (
	"reflect"

	"github.com/reflect/filq/context"
)

func DefineIn(ctx *context.Context) {
	ctx.DefineConverter(reflect.TypeOf([]interface{}{}), &ArrayConverter{})
	ctx.DefineConverter(reflect.TypeOf([]byte{}), &BytesConverter{})
	ctx.DefineConverter(reflect.TypeOf(map[string]interface{}{}), &ObjectConverter{})
	ctx.DefineConverter(reflect.TypeOf(""), &StrConverter{})
}
