package nft

import (
	"math/rand"
	"time"
)

type Choice[T any, W uint] struct {
	Item   T
	Weight W
}

// NewChoice creates a new Choice with specified item and weight.
func NewChoice[T any, W uint](item T, weight W) Choice[T, W] {
	return Choice[T, W]{Item: item, Weight: weight}
}

type RandomGenerator[T any, W uint] struct {
	Choices   []Choice[T, W]
	Cdf       []W
	WeightSum W
	ResultArr []T
	Result    *Choice[T, W]
	Decimal   int
}

func NewPicker[T any, W uint](t T, decimal int) *RandomGenerator[T, W] {
	return &RandomGenerator[T, W]{Decimal: decimal}
}

func (g *RandomGenerator[T, W]) calcCDF() {
	//rand.Seed(time.Now().UnixNano())
	weightSum := W(0)
	n := len(g.Choices)
	cdf := make([]W, n)
	for i, w := range g.Choices {
		if i > 0 {
			cdf[i] = cdf[i-1] + w.Weight
		} else {
			cdf[i] = w.Weight
		}
		weightSum += w.Weight
	}
	g.Cdf = cdf
	g.WeightSum = weightSum
}

func (g *RandomGenerator[T, W]) FromChoices(choices ...Choice[T, W]) *RandomGenerator[T, W] {
	g.Choices = choices
	g.calcCDF()
	return g
}
func (g *RandomGenerator[T, W]) AddChoice(choices ...Choice[T, W]) {
	g.Choices = append(g.Choices, choices...)
	g.calcCDF()
}
func (g *RandomGenerator[T, W]) PickOne(optional bool) *T {
	found := false
	for !found {
		found = g.WeightedGenerate()
		if optional {
			break
		}
	}
	if g.Result != nil {
		return &g.Result.Item
	} else {
		return nil
	}
}
func (g *RandomGenerator[T, W]) GenerateN(n int, unique bool, optional bool) {
	//g.ResultArr = make([]T, 0)
	choices := make([]Choice[T, W], len(g.Choices))
	copy(choices, g.Choices)
	count := 0
	for count < n && g.WeightSum > 0 {
		g.PickOne(optional)
		if g.Result != nil {
			if unique {
				g.WeightSum -= g.Result.Weight
				g.Result.Weight = 0
				g.calcCDF()
			}
			g.ResultArr = append(g.ResultArr, g.Result.Item)
		}
		count += 1
	}
	copy(g.Choices, choices)
	g.calcCDF()
}
func (g *RandomGenerator[T, W]) WeightedGenerate() bool {
	rand.Seed(time.Now().UnixNano())
	n := len(g.Choices)
	if g.WeightSum <= 0 {
		return false
	}
	//ranNum := rand.Intn(int(math.Pow10(g.Decimal)))
	ranNum := rand.Intn(int(g.WeightSum))
	r := uint32(ranNum)
	var l, h = 0, n - 1
	for l <= h {
		m := l + (h-l)/2
		if r <= uint32(g.Cdf[m]) {
			if m == 0 || (m > 0 && r > uint32(g.Cdf[m-1])) {
				g.Result = &g.Choices[m]
				return true
			}
			h = m - 1
		} else {
			l = m + 1
		}
	}
	g.Result = nil
	return false
}
