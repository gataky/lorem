package lorem

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

const (
	TAG = "lorem"
)

type Options struct {
	Seed     int64
	SliceLen int
	MapLen   int
}

var defaultOptions = Options{
	Seed:     time.Now().UnixNano(),
	SliceLen: 3,
	MapLen:   3,
}

type Lorem struct {
	seed       int64
	rand       *rand.Rand
	primitives map[reflect.Kind]Provider
	categories map[string]Provider
	providers  map[string]Provider

	sliceLen int
	arrayLen int
	mapLen   int
}

var primitives = map[reflect.Kind]Provider{
	reflect.Bool:       Bool,
	reflect.String:     String,
	reflect.Int:        Int,
	reflect.Int8:       Int8,
	reflect.Int16:      Int16,
	reflect.Int32:      Int32,
	reflect.Int64:      Int64,
	reflect.Uint:       Uint,
	reflect.Uint8:      Uint8,
	reflect.Uint16:     Uint16,
	reflect.Uint32:     Uint32,
	reflect.Uint64:     Uint64,
	reflect.Float32:    Float32,
	reflect.Float64:    Float64,
	reflect.Complex64:  Complex64,
	reflect.Complex128: Complex128,
}

var categories = map[string]Provider{
	"MaleName":   MaleName,
	"LastName":   LastName,
	"FemaleName": FemaleName,
}

func NewLorem(opts ...Options) *Lorem {
	var options Options
	if len(opts) > 0 {
		options = opts[0]
	}

	return &Lorem{
		seed:       options.Seed,
		rand:       rand.New(rand.NewSource(options.Seed)),
		primitives: primitives,
		categories: categories,
		providers:  make(map[string]Provider),
		sliceLen:   options.SliceLen,
		mapLen:     options.MapLen,
	}
}

// RegisterProvider will register a provider function to a tag on a struct.
// When the field of that tag is being processed it will use the registered
// function to generate its value.
func (l Lorem) RegisterProvider(tag string, provider Provider) {
	l.providers[tag] = provider
}

// ResetRandom will reset the random generator to the start of the seed sequence.
func (l *Lorem) ResetRandom() {
	l.rand = rand.New(rand.NewSource(l.seed))
}

// Fake will do the actual faking of the data. Source must be a pointer to the
// variable so Fake can set the value.
func (l Lorem) Fake(source any) error {
	// Get the reflect.Value of the pointer
	valueOfSource := reflect.ValueOf(source)

	// Ensure it's a pointer otherwise we won't be able to set the value.
	if valueOfSource.Kind() != reflect.Pointer {
		panic("Expected a pointer")
	}

	element := valueOfSource.Elem()

	fakeValue, err := l.fakeIt(element)
	if err != nil {
		return err
	}

	element.Set(fakeValue.Convert(element.Type()))
	return nil
}

func (l *Lorem) fakeIt(element reflect.Value) (reflect.Value, error) {
	switch kind := element.Kind(); kind {

	// primitive cases are the types that can't be broken down anymore.
	case reflect.Bool:
		fallthrough

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fallthrough

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fallthrough

	case reflect.Float32, reflect.Float64:
		fallthrough

	case reflect.Complex64, reflect.Complex128:
		fallthrough

	case reflect.String:
		value := l.primitives[kind](l.rand)
		return reflect.ValueOf(value), nil

		// collection cases get recursed into until the reach the primitive case.
	case reflect.Pointer:
		return l._pointer(element)

	case reflect.Map:
		return l._map(element)

	case reflect.Slice:
		return l._slice(element)

	case reflect.Array:
		return l._array(element)

	case reflect.Struct:
		return l._struct(element)

	default:
		return reflect.Value{}, fmt.Errorf("unsupported field kind: %s", kind)

	}
}

func (l *Lorem) _pointer(element reflect.Value) (reflect.Value, error) {
	elementType := element.Type()
	// subItem is the type the value points to i.e. if the pointer is a
	// *string then the subItem would be a string
	subItem := reflect.Zero(elementType.Elem())
	pointer := reflect.New(subItem.Type())

	value, err := l.fakeIt(subItem)
	if err != nil {
		return reflect.Value{}, err
	}

	pointer.Elem().Set(value.Convert(subItem.Type()))
	return pointer, nil
}

func (l *Lorem) _map(element reflect.Value) (reflect.Value, error) {
	elementType := element.Type()
	newMap := reflect.MakeMap(elementType)
	for i := 0; i < l.mapLen; i++ {

		key, err := l.fakeIt(reflect.New(elementType.Key()).Elem())
		if err != nil {
			return reflect.Value{}, err
		} else if reflect.ValueOf(key).IsZero() {
			continue
		}
		key = key.Convert(elementType.Key())

		value, err := l.fakeIt(reflect.New(elementType.Elem()).Elem())
		if err != nil {
			return reflect.Value{}, err
		} else if reflect.ValueOf(value).IsZero() {
			continue
		}
		value = value.Convert(elementType.Elem())

		newMap.SetMapIndex(key, value)
	}
	return newMap, nil
}

func (l *Lorem) _slice(element reflect.Value) (reflect.Value, error) {
	elementType := element.Type()
	newSlice := reflect.MakeSlice(elementType, l.sliceLen, l.sliceLen)
	itemType := newSlice.Index(0).Type()

	for i := 0; i < l.sliceLen; i++ {
		value, err := l.fakeIt(newSlice.Index(i))
		if err != nil {
			return reflect.Value{}, err
		} else if reflect.ValueOf(value).IsZero() {
			continue
		}
		newSlice.Index(i).Set(value.Convert(itemType))
	}
	return newSlice, nil
}

func (l *Lorem) _array(element reflect.Value) (reflect.Value, error) {
	elementType := element.Type()
	newArray := reflect.New(elementType).Elem()
	itemType := newArray.Index(0).Type()

	for i := 0; i < newArray.Len(); i++ {
		value, err := l.fakeIt(newArray.Index(i))
		if err != nil {
			return reflect.Value{}, err
		} else if reflect.ValueOf(value).IsZero() {
			continue
		}

		newArray.Index(i).Set(value.Convert(itemType))
	}
	return newArray, nil
}

func (l *Lorem) _struct(element reflect.Value) (reflect.Value, error) {
	elementType := element.Type()
	newStruct := reflect.New(elementType).Elem()

	for i := 0; i < newStruct.NumField(); i++ {
		field := newStruct.Field(i)
		if !field.CanSet() {
			continue
		}

		tag := elementType.Field(i).Tag.Get(TAG)
		var (
			err   error
			value reflect.Value
		)

		// ignore fields that have a "-" tag.  Useful for interfaces until
		// that gets sorted out.
		if tag == "-" {
			continue
		} else if provider, ok := l.providers[tag]; ok {
			value = reflect.ValueOf(provider(l.rand))
		} else if provider, ok := categories[tag]; ok {
			value = reflect.ValueOf(provider(l.rand))
		} else {
			value, err = l.fakeIt(field)
		}

		if err != nil {
			return reflect.Value{}, err
		} else if reflect.ValueOf(value).IsZero() {
			continue
		}

		field.Set(value.Convert(field.Type()))
	}
	return newStruct, nil
}

func inspect(v reflect.Value) {
	fmt.Println("===========================================================")
	fmt.Println("        v          ", v)
	fmt.Println("        v.interface", v.Interface())
	fmt.Println("        v.elem     ", v.Elem())
	fmt.Println("        v.kind     ", v.Kind())
	fmt.Println("        v.type     ", v.Type())
	fmt.Println("        v.type.elem", v.Type().Elem())
	fmt.Println("valueof v.type.elem", reflect.ValueOf(v.Type().Elem()))
	fmt.Println("        v.isvalid  ", v.IsValid())
	fmt.Println("        v.iszero   ", v.IsZero())
	fmt.Println("===========================================================")
}
