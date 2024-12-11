package main

import (
	"math/rand"
)

func providerInt(rand *rand.Rand) any {
	return rand.Int()
}

func providerUInt(rand *rand.Rand) any {
	return rand.Uint64()
}

func providerFloat(rand *rand.Rand) any {
	return rand.Float64()
}

func providerComplex(rand *rand.Rand) any {
	return complex(rand.Float32(), rand.Float32())
}

func providerBool(rand *rand.Rand) any {
	return rand.Intn(2) > 0
}

const letters = "abcdefghijklmnopqrstuvwxyz"

func providerString(rand *rand.Rand) any {
	b := make([]byte, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
