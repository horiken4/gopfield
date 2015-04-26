package main

import (
	"fmt"
)

func main() {
	// Associative memory

	h := NewHopfield(5 * 5)

	// Train
	pats := [][]float32{
		[]float32{
			1, 1, 1, 1, 1,
			-1, 1, -1, 1, -1,
			1, 1, 1, 1, 1,
			1, -1, 1, -1, 1,
			1, 1, 1, 1, 1,
		},
		[]float32{
			-1, -1, 1, 1, 1,
			1, 1, 1, -1, -1,
			-1, -1, 1, 1, 1,
			1, 1, 1, -1, -1,
			-1, -1, 1, 1, 1,
		},
		[]float32{
			1, -1, 1, -1, 1,
			-1, 1, -1, 1, -1,
			1, -1, 1, -1, 1,
			-1, 1, -1, 1, -1,
			1, -1, 1, -1, 1,
		},
	}
	h.Train(pats)

	// Remember
	initPats := [][]float32{
		[]float32{
			1, 1, 1, 1, 1,
			1, -1, 1, -1, 1,
			1, 1, -1, 1, 1,
			1, -1, 1, 1, 1,
			1, 1, 1, -1, 1,
		},
		[]float32{
			-1, 1, 1, 1, 1,
			1, 1, 1, -1, 1,
			1, -1, 1, 1, 1,
			1, 1, 1, -1, -1,
			1, -1, 1, -1, 1,
		},
		[]float32{
			-1, 1, -1, 1, 1,
			1, -1, 1, -1, 1,
			1, -1, 1, 1, -1,
			1, -1, 1, -1, -1,
			1, -1, 1, -1, 1,
		},
	}
	for _, pat := range initPats {
		h.Feed(pat)
		h.Run(50)
		h.Print(5)
		fmt.Println("energy =", h.Energy())
	}

	/*
		// Solve 0-1 Knapsack probrem

		// A solution is 1 -1 1 1 -1 -1
		C := 14                          // Total capacity
		v := []int{10, 13, 10, 16, 2, 3} // Item value
		c := []int{3, 5, 4, 7, 2, 4}     // Item capacity
		var A float32 = 0.3

		h := NewHopfield(6)

		for i := 0; i < len(h.Neurons); i++ {
			for j := 0; j < len(h.Neurons); j++ {
				if i == j {
					continue
				}
				h.SetWeights(i, j, float32(-2*A*float32(c[i])*float32(c[j])))
			}
		}
		for i := 0; i < len(h.Neurons); i++ {
			h.SetThreshold(i, -A*float32(C)*float32(c[i])-float32(v[i]))
		}

		h.FeedRandomly()
		h.Run(100)
		h.Print(6)
		fmt.Println("energy =", h.Energy())
	*/
}
