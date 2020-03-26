package main

import (

)

func TestProbTable(t *testing.T) {
	examples := []struct{
		name string
		want map[string]float64
		goodMap map[string]int
		badMap map[string]int
		nGoodMail int
		nBadMail int
	}{
		{
			name: control,
			want: map[string]float64{
				"the": {0.3478260869565218},
				"quick": {0.01},
				"brown": {0.01},
				"fox": {0.01},
				"jumped": {0.5},
				"over": {0.01},
				"lazy": {0.5714285714285715},
				"white": {0.7499999999999999},
				"dog": {0.01},
				"slow": {0.99},
				"computer": {0.99},
			},
			goodMap: map[string]int{
				"the": {10},
				"quick": {3},
				"brown": {6},
				"fox": {8},
				"jumped": {2},
				"over": {4},
				"lazy": {5},
				"white": {2}
				"dog": {12},
			},
			badMap: map[string]int{
				"the": {8},
				"slow": {4},
				"white": {9},
				"computer": {10},
				"jumped": {3},
				"to": {3},
				"lazy": {10},
				"function": {4},
			},
			nGoodMail: 20,
			nBadMail: 15,
		},
	}
	for _, ex := range examples {
		t.Run(ex.name, func(t *testing.T) {
			got := probTable(ex.goodMap, ex.badMap, ex.nGoodMail, ex.nBadMail)
			if ex.want != got {
				t.Fatalf("got %s, want %s", got, ex.want)
			}
		})
	}
}