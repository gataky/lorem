# Lorem

Lorem generates fake data for primitive variables (e.g., int, float, string) and collections (e.g., map, slice, array, structs). It also supports faking custom types and pointers. Lorem features a simple tagging system that allows for custom provider registration.

This project was inspired by (faker)[https://github.com/go-faker/faker], which I found very useful. However, I wanted a way to produce reproducible random generations, and achieving that with faker wasn't straightforward based on my experience.

# Usage

## Initial Defaults

Lorem is defaulted with the following options

```go
var defaultOptions = Options{
	Seed:     time.Now().UnixNano(),
	SliceLen: 3,
	MapLen:   3,
}
```

* Seed: is the seed for the random generator to use.
* SliceLen: specifies the number of elements a slice will have.
* MapLen: specifies the number of elements a map will have.

You can override these defaults by passing a custom `lorem.Options` struct to `lorem.New`:

```go
myOptions := lorem.Options{
    Seed: 123,
    SliceLen: 10,
    MapLen: 5,
}
```

## Lorem

To use Lorem, create a new instance. You can use the default options or provide your custom options.

```go
// Default options
l := lorem.New()

// Custom options
l := lorem.New(myOptions)
```

## Generating Fake Data

To generate fake data, pass a pointer to the element you want to fake to the `Fake` method:

```go
test := 0
l.Fake(&test)
```

## Slice

To generate a slice created with make, the length specified with make will be used over the options length. For example, if we have a slice with 5 elements and we pass to `Fake`.

```go
mySlice = make([]string, 5, 6)
l := lorem.New()
l.Fake(&mySlice)
```

The slice will have 5 random values and a capacity of 6.  However, slices specified in a struct use the options `SliceLen`.

```go
type S struct {
    slice []string
}
myS := &S{}
l := lorem.New()
l.Fake(&mySlice)
```

This will create a slice with 3 element.

## Custom Types

Custom types, such as `var MyInt int`, will be processed down to their primitive type. For example, an `int` will be assigned.

## Struct

Lorem can recursively generate fake data for structs and their fields:

```go
type Number int

type S struct {
	Time    time.Time
	Pointer *int
	Custom  Number
	String  string
	Int     int8
	Float   float32
	Map     map[Number]string
	Slice   []string
}


o := lorem.Options{
    Seed:     123,
    SliceLen: 10,
    MapLen:   3,
}
l := lorem.New(o)

s := &S{}
l.Fake(s)
spew.Dump(s)
```

This will produce the following output

```go
(lorem_test.S) {
 Time: (time.Time) 2012-07-17 08:38:09 +0000 UTC,
 Pointer: (*int)(0x14000104230)(241876450138978746),
 Custom: (lorem_test.Number) 2305561650894865143,
 String: (string) (len=6) "ldddpx",
 Int: (int8) 71,
 Float: (float32) 0.51044065,
 Map: (map[lorem_test.Number]string) (len=3) {
  (lorem_test.Number) 4640137937568901621: (string) (len=6) "kjafbn",
  (lorem_test.Number) 7910589188243225516: (string) (len=6) "wksibn",
  (lorem_test.Number) 5055213988076362636: (string) (len=6) "doakse"
 },
 Slice: ([]string) (len=10 cap=10) {
  (string) (len=6) "qtxcuz",
  (string) (len=6) "ldccrl",
  (string) (len=6) "zssynj",
  (string) (len=6) "jhcile",
  (string) (len=6) "wklrls",
  (string) (len=6) "tegngm",
  (string) (len=6) "nnguun",
  (string) (len=6) "jbbxvl",
  (string) (len=6) "aqtdxw",
  (string) (len=6) "idphti"
 }
}

```

## Tags

A simple tag system controls how fields are generated. To ignore a field, use the lorem:"-" tag.

### Predefined Providers

Lorem includes several predefined providers, which can be specified via tags to control the values generated. For example:

```go
type S struct {
	String  string `lorem:"LastName"`
}
```

produces

```go
lorem_test.S) {
 String: (string) (len=9) "Rodriguez"
}
```

Check the `mappings.go` file for a list of available providers.

### Custom providers

You can register custom providers by creating a function that accepts a `*rand.Rand` argument and returns `any`. Register this function with Lorem using a unique name:

```go
type S struct {
	Slice []any `lorem:"myFakeSlice"`
}

func myFakeSliceProvider(rand *rand.Rand) any {
	// you don't have to use lorem's providers here but they're available
	// to you is you want.
	a := lorem.String(rand).(string)
	b := lorem.String(rand).(string)
	c := lorem.Int8(rand).(int8)

	return []any{a, b, c}
}

l := lorem.New()
l.RegisterProvider("myFakeSlice", myFakeSliceProvider)
```

This is particularly useful for fields of type `any` or `interface{}`, which Lorem cannot process without additional information.

### Future Tags

Tags may expand to handle more complex cases in future versions.

# Example

For a full working example, check out `example_test.go`. This file demonstrates all the features described here.
