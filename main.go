package main

import (
	"fmt"
	"math/rand"
	"reflect"

	"github.com/davecgh/go-spew/spew"
)

type Int int

type S struct {
	Pointer *Int
	String  string
	Map     map[Int]Int
}

func main() {
	lorem := NewLorem()
	// test := "foo"

	// var pi *int = &test
	// lorem.Fake(pi)
	// fmt.Println(*pi)

	// var i int
	// lorem.Fake(&i)
	// fmt.Println(i)
	//
	// m := map[string]Int{}
	// lorem.Fake(&m)
	// fmt.Println(m)
	//
	// ss := []string{}
	// lorem.Fake(&ss)
	// fmt.Println(ss)
	//
	// sn := []Int{}
	// lorem.Fake(&sn)
	// fmt.Println(sn)
	//
	// a := [2]Int{}
	// lorem.Fake(&a)
	// fmt.Println(a)

	st := S{}
	lorem.Fake(&st)
	spew.Dump(st)
	fmt.Println(*st.Pointer)
}

type Lorem struct {
	seed      int64
	rand      *rand.Rand
	providers map[reflect.Kind]Provider
	count     int
}

type Kind uint

type Provider func(*rand.Rand) reflect.Value

func NewLorem() *Lorem {
	// seed := time.Now().UnixNano()
	seed := int64(1)

	providers := map[reflect.Kind]Provider{
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

	return &Lorem{
		seed:      seed,
		rand:      rand.New(rand.NewSource(seed)),
		providers: providers,
		count:     0,
	}
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
		return l.providers[kind](l.rand), nil

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
		elementType := element.Type()
		newMap := reflect.MakeMap(elementType)
		// TODO: configure length
		for i := 0; i < 2; i++ {

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
		elementType := element.Type()
		newSlice := reflect.MakeSlice(elementType, 2, 2)
		itemType := newSlice.Index(0).Type()

		// TODO: configure length
		for i := 0; i < 2; i++ {
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
		elementType := element.Type()
		newArray := reflect.New(elementType).Elem()
		itemType := newArray.Index(0).Type()

		for i := 0; i < 2; i++ {
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
		elementType := element.Type()
		newStruct := reflect.New(elementType).Elem()

		for i := 0; i < newStruct.NumField(); i++ {
			field := newStruct.Field(i)
			if !field.CanSet() {
				continue
			}

			value, err := l.fakeIt(field)
			if err != nil {
				return reflect.Value{}, err
			} else if reflect.ValueOf(value).IsZero() {
				continue
			}

			fmt.Println("unpacked value", value)
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
