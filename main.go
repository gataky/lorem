package main

import (
	"fmt"
	"math/rand"
	"reflect"
)

func main() {

	lorem := NewLorem()

	var a int
	lorem.Fake(&a)
	fmt.Println(a)

	m := map[string]string{}
	lorem.Fake(&m)
	fmt.Println(m)

}

type Lorem struct {
	seed      int64
	rand      *rand.Rand
	providers map[reflect.Kind]Provider
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
		reflect.Array:      nil,
		reflect.Map:        nil,
		reflect.Slice:      nil,
		reflect.String:     nil,
		reflect.Struct:     nil,
	}

	return &Lorem{
		seed:      seed,
		rand:      rand.New(rand.NewSource(seed)),
		providers: providers,
	}
}

func (l Lorem) Fake(source any) error {
	// Get the reflect.Value of the pointer
	valueOfSource := reflect.ValueOf(source)

	// Ensure it's a pointer
	if valueOfSource.Kind() != reflect.Pointer {
		panic("Expected a pointer")
	}

	// Get the element type and element
	element := valueOfSource.Elem()
	// fmt.Println("element.Type", element.Type())
	// fmt.Println("element.Kind", element.Kind())

	fakeValue, err := l.fakeIt(element)
	if err != nil {
		return err
	}

	element.Set(fakeValue.Convert(element.Type()))
	return nil
}

func (l Lorem) fakeIt(element reflect.Value) (reflect.Value, error) {

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
		return l.providers[kind](l.rand), nil

	case reflect.Array:
	case reflect.Map:
		item := reflect.MakeMap(element.Type())
		for i := 0; i < 1; i++ {

			key, err := l.fakeIt(reflect.New(element.Type().Key()).Elem())
			if err != nil {
				return reflect.Value{}, err
			} else if reflect.ValueOf(key).IsZero() {
				continue
			}

			value, err := l.fakeIt(reflect.New(element.Type().Elem()).Elem())
			if err != nil {
				return reflect.Value{}, err
			} else if reflect.ValueOf(value).IsZero() {
				continue
			}

			item.SetMapIndex(key, value)
		}
		return item, nil

	case reflect.Slice:
	case reflect.String:
		return reflect.ValueOf("foo"), nil

	case reflect.Struct:
	default:
		return reflect.Value{}, fmt.Errorf("unsupported kind: %s", kind)

	}
	return reflect.Value{}, nil
}

func providerInt(rand *rand.Rand) reflect.Value {
	return reflect.ValueOf(rand.Int())
}

func providerUInt(rand *rand.Rand) reflect.Value {
	return reflect.ValueOf(rand.Uint64())
}

func providerFloat(rand *rand.Rand) reflect.Value {
	return reflect.ValueOf(rand.Float64())
}

func providerComplex(rand *rand.Rand) reflect.Value {
	return reflect.ValueOf(complex(rand.Float32(), rand.Float32()))
}

func providerBool(rand *rand.Rand) reflect.Value {
	value := rand.Intn(2) > 0
	return reflect.ValueOf(value)
}
