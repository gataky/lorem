package lorem

import (
	"math/rand"
	"time"
)

func Time(r *rand.Rand) any {

	minTime := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	maxTime := time.Date(2040, 1, 1, 0, 0, 0, 0, time.UTC)

	diffSeconds := int64(maxTime.Sub(minTime) / time.Second)

	randomSeconds := r.Int63n(diffSeconds)

	randomTime := minTime.Add(time.Duration(randomSeconds) * time.Second)
	return randomTime
}
