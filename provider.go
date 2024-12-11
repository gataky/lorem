package main

import (
	"math/rand"
)

type Provider func(*rand.Rand) any

const letters = "abcdefghijklmnopqrstuvwxyz"

func String(rand *rand.Rand) any {
	b := make([]byte, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func Bool(rand *rand.Rand) any {
	return rand.Intn(2) > 0
}

func Int(r *rand.Rand) any {
	return int(r.Int())
}

func Int8(r *rand.Rand) any {
	return int8(r.Intn(256))
}

func Int16(r *rand.Rand) any {
	return int16(r.Intn(65536))
}

func Int32(r *rand.Rand) any {
	return r.Int31()
}

func Int64(r *rand.Rand) any {
	return r.Int63()
}

func Uint(r *rand.Rand) any {
	return uint(r.Intn(256))
}

func Uint8(r *rand.Rand) any {
	return uint8(r.Intn(256))
}

func Uint16(r *rand.Rand) any {
	return uint16(r.Intn(65536))
}

func Uint32(r *rand.Rand) any {
	return r.Uint32()
}

func Uint64(r *rand.Rand) any {
	return r.Uint64()
}

func Float32(r *rand.Rand) any {
	return r.Float32()
}

func Float64(r *rand.Rand) any {
	return r.Float64()
}

func Complex64(r *rand.Rand) any {
	return complex(r.Float32(), r.Float32())
}

func Complex128(r *rand.Rand) any {
	return complex(r.Float64(), r.Float64())
}
