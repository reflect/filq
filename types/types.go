package types

import (
	"reflect"

	"github.com/reflect/filq/context"
)

func DefineIn(ctx *context.Context) {
	ctx.DefineConverter(reflect.TypeOf([]interface{}{}), &ArrayConverter{})
	ctx.DefineConverter(reflect.TypeOf([]byte{}), &BytesConverter{})
	ctx.DefineConverter(reflect.TypeOf(float64(0)), &FloatFloat64Converter{})
	ctx.DefineConverter(reflect.TypeOf(float32(0)), &FloatFloat32Converter{})
	ctx.DefineConverter(reflect.TypeOf(int64(0)), &IntInt64Converter{})
	ctx.DefineConverter(reflect.TypeOf(int32(0)), &IntInt32Converter{})
	ctx.DefineConverter(reflect.TypeOf(int16(0)), &IntInt16Converter{})
	ctx.DefineConverter(reflect.TypeOf(int8(0)), &IntInt8Converter{})
	ctx.DefineConverter(reflect.TypeOf(int(0)), &IntIntConverter{})
	ctx.DefineConverter(reflect.TypeOf(map[string]interface{}{}), &ObjectConverter{})
	ctx.DefineConverter(reflect.TypeOf(""), &StrConverter{})
}
