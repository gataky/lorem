package lorem

import "math/rand"

var (
	Female = []string{
		"Amy", "Turanga ",
	}

	Male = []string{
		"Bender", "Hermes", "Hubert", "John", "Philip ",
	}

	Last = []string{
		"Conrad", "Farnsworth", "Fry", "Leela", "Rodriguez", "Wong", "Zoidberg",
	}
)

func pick(r *rand.Rand, list []string) string {
	return list[r.Intn(len(list))]
}

func FemaleName(r *rand.Rand) any {
	return pick(r, Female)
}

func MaleName(r *rand.Rand) any {
	return pick(r, Male)
}

func LastName(r *rand.Rand) any {
	return pick(r, Last)
}
