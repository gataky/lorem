package main

import (
	"math/rand"
	"reflect"
)

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

const letters = "abcdefghijklmnopqrstuvwxyz"

func providerString(rand *rand.Rand) reflect.Value {
	b := make([]byte, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return reflect.ValueOf(string(b))
}
