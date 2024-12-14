package lorem_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gataky/lorem"
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
	Time    time.Time
	Pointer *int
	Custom  Number
	String  string `lorem:"LastName"`
	Int     int8
	// To ignore a field add a "-" tag
	IgnoreField int8 `lorem:"-"`
	Float       float32
	Map         map[Number3]int
	MapMap      map[string]map[string]string
	Slice       []any `lorem:"myFakeSlice"`
	SubStruct   SubStruct
	Func        func() `lorem:"providerFunc"`
}

// Example of using a custom provider to generate data for you. Useful if you
// have an interface value or slices to interfaces.
func myFakeSliceProvider(rand *rand.Rand) any {
	// you don't have to use lorem's providers here but they're available
	// to you is you want.
	a := lorem.String(rand).(string)
	b := lorem.String(rand).(string)
	c := lorem.Int8(rand).(int8)

	return []any{a, b, c}
}

// This is the provider that will "generate" the fake function above. All
// providers take a math.Rand struct to give us a deterministic random sample
// of values.
func funcProvider(rand *rand.Rand) any {
	return fake
}

// This is the fake function that will be set on the sturct
func fake() {
	fmt.Println("hi from the fake func")
}

func Test_Example(t *testing.T) {
	// Some options to control the generated types
	o := lorem.Options{
		// Seed so we can have reproducible generations.
		SliceLen: 10, // The len of slices
		MapLen:   3,  // the len of maps
	}
	l := lorem.New(o)

	// Register the providers for the tags specified on the struct.
	l.RegisterProvider("providerFunc", funcProvider)
	l.RegisterProvider("myFakeSlice", myFakeSliceProvider)

	st := S{}
	l.Fake(&st)
	spew.Dump(st)

	// The func will return the function that we used from our custom method.
	st.Func() // -> "hi from the fake func"

	t.Fail()
}
