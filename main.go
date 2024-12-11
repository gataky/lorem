package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/davecgh/go-spew/spew"
)

type Int int

type SS struct {
	String string
}

type S struct {
	Pointer *Int        `lorem:"Int,pointer"`
	String  string      `lorem:"custom"`
	Map     map[Int]Int `lorem:"map"`
	S       SS
}

func main() {
	lorem := NewLorem()
	lorem.RegisterProvider("custom", custom)

	test := 10
	var pi *int = &test
	lorem.Fake(pi)
	fmt.Println(*pi)

	var i int
	lorem.Fake(&i)
	fmt.Println(i)

	m := map[string]Int{}
	lorem.Fake(&m)
	fmt.Println(m)

	ss := []string{}
	lorem.Fake(&ss)
	fmt.Println(ss)

	sn := []Int{}
	lorem.Fake(&sn)
	fmt.Println(sn)

	a := [2]Int{}
	lorem.Fake(&a)
	fmt.Println(a)

	st := S{}
	lorem.Fake(&st)
	spew.Dump(st)
	fmt.Println(*st.Pointer)
}

func custom(rand *rand.Rand) any {
	return "foo"
}

const (
	TAG = "lorem"
)

type Options struct {
	Seed     int64
	SliceLen int
	ArrayLen int
	MapLen   int
}

type Lorem struct {
	seed       int64
	rand       *rand.Rand
	primitives map[reflect.Kind]Provider
	providers  map[string]Provider

	sliceLen int
	arrayLen int
	mapLen   int
}

type Provider func(*rand.Rand) any

func NewLorem() *Lorem {
	seed := time.Now().UnixNano()

	primitives := map[reflect.Kind]Provider{
		reflect.Bool:       providerBool,
		reflect.Int:        providerInt,
		reflect.Int8:       providerInt,
		reflect.Int16:      providerInt,
		reflect.Int32:      providerInt,
		reflect.Int64:      providerInt,
		reflect.Uint:       providerUInt,
		reflect.Uint8:      providerUInt,
		reflect.Uint16:     providerUInt,
		reflect.Uint32:     providerUInt,
		reflect.Uint64:     providerUInt,
		reflect.Float32:    providerFloat,
		reflect.Float64:    providerFloat,
		reflect.Complex64:  providerComplex,
		reflect.Complex128: providerComplex,
		reflect.String:     providerString,
	}

	providers := map[string]Provider{}

	return &Lorem{
		seed:       seed,
		rand:       rand.New(rand.NewSource(seed)),
		primitives: primitives,
		providers:  providers,

		sliceLen: 3,
		arrayLen: 3,
		mapLen:   3,
	}
}

func (l Lorem) RegisterProvider(tag string, provider Provider) {
	l.providers[tag] = provider
}

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

	elementType := element.Type()
	switch kind := element.Kind(); kind {
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

	case reflect.Pointer:

		subItem := reflect.Zero(element.Type().Elem())
		pointer := reflect.New(subItem.Type())

		value, err := l.fakeIt(subItem)
		if err != nil {
			return reflect.Value{}, err
		}

		pointer.Elem().Set(value.Convert(subItem.Type()))
		return pointer, nil

	case reflect.Map:
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

	case reflect.Slice:
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

	case reflect.Array:
		newArray := reflect.New(elementType).Elem()
		itemType := newArray.Index(0).Type()

		for i := 0; i < l.arrayLen; i++ {
			value, err := l.fakeIt(newArray.Index(i))
			if err != nil {
				return reflect.Value{}, err
			} else if reflect.ValueOf(value).IsZero() {
				continue
			}

			newArray.Index(i).Set(value.Convert(itemType))
		}
		return newArray, nil

	case reflect.Struct:
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

			if f, ok := l.providers[tag]; ok {
				value = reflect.ValueOf(f(l.rand))
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

	default:
		return reflect.Value{}, fmt.Errorf("unsupported kind: %s", kind)

	}
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
