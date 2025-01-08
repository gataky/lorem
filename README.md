# Lorem

Lorem will generate fake data for your primitive variables (int, float, string, ...) and collections (map, slice, array, structs).  You can also fake custom types and pointers.  Lorem has a very simple tag system allowing for custom provider registration.

This project was inspired by (faker)[https://github.com/go-faker/faker] which was very useful to me but I really wanted a way of having reproducible random generations and there wasn't an easy way of achieving that with faker from what I could see.

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

You can specify your own options by passing in a new `lorem.Options` struct with the values you want to `lorem.New`

```go
myOptions := lorem.Options{
    Seed: 123,
    SliceLen: 10,
    MapLen: 5,
}
```

## Lorem

To use lorem, you'll need to create a new lorem.  You can use lorem without any options which will use the default options or pass in your custom options.

```go
// Default options
l := lorem.New()

// Custom options
l := lorem.New(myOptions)
```

## Generating Fake Data

To generate fake data you must pass a pointer of the element to `Fake`.

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

## Custom Type

Custom types like `var MyInt int`, will recursively dive down to the primitive type. So in the end an int will be assigned.

## Struct

This is perhaps the most useful case, recursively generating a struct and fields within it.

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

Which will produce the following struct

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

A simple tag system exists to control the behavior of generating fields.

To tell the generator to ignore a field on a struct you can use `lorem:"-"` for that field

### Stock providers

Lorem comes with some predefined providers which can be specified through tags to control the values a field will generate.  For example,

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

The stock providers are currently limited and will expand over time.  The easiest way to see what's available is to look in the `mappings.go` file. The `categories` file will have the list of available providers.

### Custom providers

To use your own provider you can create a function that accepts a `*rand.Rand` argument and returns `any` and register that provider with lorem with a name that will be used in a tag

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

This is useful when you and specific types a field should be or when you have a field that's of type `any` or `interface{}` because lorem will ignore any field that's unknown.  If you don't know what it is how can lorem? So you have to control it.

### Future Tags

Tags will probably expand to handle more complex cases.

# Example

To see a full working example checkout `example_test.go` There you'll find all the features described here.
