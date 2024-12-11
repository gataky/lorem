package main

import (
	"fmt"
	"math/rand"

	"github.com/davecgh/go-spew/spew"
)

type (
	Number  int
	Number2 Number
	Number3 Number2
)

type SubStruct struct {
	String string
}

type S struct {
	Pointer   *int
	Custom    Number
	String    string
	Int       int8
	Float     float32
	Map       map[Number3]int
	MapMap    map[string]map[string]string
	Slice     []string
	SubStruct SubStruct
	Func      func() `lorem:"providerFunc"`
}

// This is the fake function that will be set on the sturct
func fake() {
	fmt.Println("hi from the fake func")
}

// This is the provider that will "generate" the fake function above. All
// providers take a math.Rand struct to give us a deterministic random sample
// of values.
func providerFunc(rand *rand.Rand) any {
	return fake
}

func main() {
	// Some options to control the generated types
	o := Options{
		// Seed so we can have reproducible generations.
		Seed:     1,
		SliceLen: 10, // The len of slices
		MapLen:   3,  // the len of maps
	}
	lorem := NewLorem(o)

	// Register the providers for the tags specified on the struct.
	lorem.RegisterProvider("providerFunc", providerFunc)

	test := 10
	var pi *int = &test
	lorem.Fake(pi)
	spew.Dump(pi)
	fmt.Println("=============================")

	var i int
	lorem.Fake(&i)
	spew.Dump(i)
	fmt.Println("=============================")

	m := map[string]Number{}
	lorem.Fake(&m)
	spew.Dump(m)
	fmt.Println("=============================")

	sn := []Number3{}
	lorem.Fake(&sn)
	spew.Dump(sn)
	fmt.Println("=============================")

	a := [2]int{}
	lorem.Fake(&a)
	spew.Dump(a)
	fmt.Println("=============================")

	st := S{}
	lorem.Fake(&st)
	spew.Dump(st)
	fmt.Println("=============================")

	st.Func()
}
