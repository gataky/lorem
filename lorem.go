package lorem

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
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
	primitives map[LoremType]Provider
	categories map[string]Provider
	providers  map[string]Provider

	initializedFromStruct bool

	sliceLen int
	arrayLen int
	mapLen   int
}

func New(opts ...Options) *Lorem {
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

	// We want to know if we initially started off with a struct of something else.  This is used
	// by slice and map to decided how many elements a slice contains.  If from a struct
	// we'll use what's provided by the options otherwise we'll use the len specified when making.
	if LoremType(element.Kind()) == TypeStruct {
		l.initializedFromStruct = true
	}

	fakeValue, err := l.fakeIt(element)
	if err != nil {
		return err
	}

	element.Set(fakeValue.Convert(element.Type()))
	return nil
}

func (l *Lorem) fakeIt(element reflect.Value) (reflect.Value, error) {

	var kind LoremType
	switch element.Type().String() {
	case "time.Time":
		kind = TypeTime
	default:
		kind = LoremType(element.Kind())
	}

	switch kind {

	// primitive cases are the types that can't be broken down anymore.
	case
		TypeBool,
		TypeInt, TypeInt8, TypeInt16, TypeInt32, TypeInt64,
		TypeUint, TypeUint8, TypeUint16, TypeUint32, TypeUint64,
		TypeFloat32, TypeFloat64,
		TypeComplex64, TypeComplex128,
		TypeString,
		TypeTime:

		value := l.primitives[kind](l.rand)
		return reflect.ValueOf(value), nil

		// collection cases get recursed into until the reach the primitive case.
	case TypePointer:
		return l.handlePointer(element)

	case TypeMap:
		return l.handleMap(element)

	case TypeSlice:
		return l.handleSlice(element)

	case TypeArray:
		return l.handleArray(element)

	case TypeStruct:
		return l.handleStruct(element)

	default:
		return reflect.Value{}, fmt.Errorf("unsupported field kind: %d", kind)

	}
}

func (l *Lorem) handlePointer(element reflect.Value) (reflect.Value, error) {
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

func (l *Lorem) handleMap(element reflect.Value) (reflect.Value, error) {
	elementType := element.Type()
	newMap := reflect.MakeMap(elementType)

	length := l.mapLen

	for i := 0; i < length; i++ {

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

func (l *Lorem) handleSlice(element reflect.Value) (reflect.Value, error) {
	elementType := element.Type()

	length := l.sliceLen
	capacity := l.sliceLen
	if l.initializedFromStruct == false {
		length = element.Len()
		capacity = element.Cap()
	}

	newSlice := reflect.MakeSlice(elementType, length, capacity)
	itemType := newSlice.Index(0).Type()

	for i := 0; i < length; i++ {
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

func (l *Lorem) handleArray(element reflect.Value) (reflect.Value, error) {
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

func (l *Lorem) handleStruct(element reflect.Value) (reflect.Value, error) {
	elementType := element.Type()
	newStruct := reflect.New(elementType).Elem()

	for i := 0; i < newStruct.NumField(); i++ {
		field := newStruct.Field(i)
		if !field.CanSet() {
			continue
		}

		tag := elementType.Field(i).Tag.Get("lorem")
		var (
			err   error
			value reflect.Value
		)

		// ignore fields that have a "-" tag.  Useful for interfaces until
		// that gets sorted out.
		if tag == "-" {
			continue
			// provider order of precedence
			// 1. user defined providers
			// 2. categorie providers
			// 3. primitive providers
		} else if provider, ok := l.providers[tag]; ok {
			value = reflect.ValueOf(provider(l.rand))
		} else if provider, ok := l.categories[tag]; ok {
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
