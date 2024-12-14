package lorem

import "reflect"

type LoremType uint8

const (
	// primitive types
	TypeBool       = LoremType(reflect.Bool)
	TypeString     = LoremType(reflect.String)
	TypeInt        = LoremType(reflect.Int)
	TypeInt8       = LoremType(reflect.Int8)
	TypeInt16      = LoremType(reflect.Int16)
	TypeInt32      = LoremType(reflect.Int32)
	TypeInt64      = LoremType(reflect.Int64)
	TypeUint       = LoremType(reflect.Uint)
	TypeUint8      = LoremType(reflect.Uint8)
	TypeUint16     = LoremType(reflect.Uint16)
	TypeUint32     = LoremType(reflect.Uint32)
	TypeUint64     = LoremType(reflect.Uint64)
	TypeFloat32    = LoremType(reflect.Float32)
	TypeFloat64    = LoremType(reflect.Float64)
	TypeComplex64  = LoremType(reflect.Complex64)
	TypeComplex128 = LoremType(reflect.Complex128)

	// collection types
	TypeArray   = LoremType(reflect.Array)
	TypeMap     = LoremType(reflect.Map)
	TypePointer = LoremType(reflect.Pointer)
	TypeSlice   = LoremType(reflect.Slice)
	TypeStruct  = LoremType(reflect.Struct)

	// custom types.  Time doesn't have its own reflect type so we're making our own to
	// distinguish between a time struct and other structs.
	TypeTime = LoremType(100)
)
