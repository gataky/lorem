package lorem

var primitives = map[LoremType]Provider{
	TypeBool:       Bool,
	TypeString:     String,
	TypeInt:        Int,
	TypeInt8:       Int8,
	TypeInt16:      Int16,
	TypeInt32:      Int32,
	TypeInt64:      Int64,
	TypeUint:       Uint,
	TypeUint8:      Uint8,
	TypeUint16:     Uint16,
	TypeUint32:     Uint32,
	TypeUint64:     Uint64,
	TypeFloat32:    Float32,
	TypeFloat64:    Float64,
	TypeComplex64:  Complex64,
	TypeComplex128: Complex128,
	TypeTime:       Time,
}

var categories = map[string]Provider{
	"MaleName":   MaleName,
	"LastName":   LastName,
	"FemaleName": FemaleName,
}
