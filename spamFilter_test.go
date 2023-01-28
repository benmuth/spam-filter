package main

import (
	"testing"
	//"fmt"
)

func TestProbTable(t *testing.T) {
	examples := []struct {
		name      string
		want      map[string]float64
		goodMap   map[string]int
		badMap    map[string]int
		nGoodMail int
		nBadMail  int
	}{
		{
			name: "control",
			want: map[string]float64{
				"the":      0.3478260869565218,
				"quick":    0.01,
				"brown":    0.01,
				"fox":      0.01,
				"jumped":   0.5,
				"over":     0.01,
				"lazy":     0.5714285714285715,
				"white":    0.7499999999999999,
				"dog":      0.01,
				"computer": 0.99,
			},
			goodMap: map[string]int{
				"the":    10,
				"quick":  3,
				"brown":  6,
				"fox":    8,
				"jumped": 2,
				"over":   4,
				"lazy":   5,
				"white":  2,
				"dog":    12,
			},
			badMap: map[string]int{
				"the":      8,
				"slow":     4,
				"white":    9,
				"computer": 10,
				"jumped":   3,
				"to":       3,
				"lazy":     10,
				"function": 4,
			},
			nGoodMail: 20,
			nBadMail:  15,
		},
	}
	for _, ex := range examples {
		t.Run(ex.name, func(t *testing.T) {
			got := probTable(ex.goodMap, ex.badMap, ex.nGoodMail, ex.nBadMail)
			for gotWord, gotValue := range got {
				//fmt.Printf("(got word: %s. got value: %v)  ", gotWord, gotValue)
				_, gotWordIsInWant := ex.want[gotWord]
				if ex.want[gotWord] != gotValue {
					t.Fatalf("got %v, want %v, for %s. ", gotValue, ex.want[gotWord], gotWord)
				}
				if gotWordIsInWant != true {
					t.Fatalf("got %s but it isn't in want", gotWord)
				}
			}
		})
	}
}

func TestIsSpam(t *testing.T) {
	examples := []struct {
		name    string
		want    bool
		mail    string
		probMap map[string]float64
	}{
		{
			name: "is Spam",
			want: true,
			mail: "computer computer computer computer white white ",
			probMap: map[string]float64{
				"the":      0.3478260869565218,
				"quick":    0.01,
				"brown":    0.01,
				"fox":      0.01,
				"jumped":   0.5,
				"over":     0.01,
				"lazy":     0.5714285714285715,
				"white":    0.7499999999999999,
				"dog":      0.01,
				"computer": 0.99,
			},
		},
		{
			name: "is not Spam",
			want: false,
			mail: "dog over fox brown quick the lazy white computer ",
			probMap: map[string]float64{
				"the":      0.3478260869565218,
				"quick":    0.01,
				"brown":    0.01,
				"fox":      0.01,
				"jumped":   0.5,
				"over":     0.01,
				"lazy":     0.5714285714285715,
				"white":    0.7499999999999999,
				"dog":      0.01,
				"computer": 0.99,
			},
		},
	}
	for _, ex := range examples {
		t.Run(ex.name, func(t *testing.T) {
			got := isSpam(ex.mail, ex.probMap)
			if got != ex.want {
				t.Fatalf("got %v, want %v ", got, ex.want)
			}
		})
	}
}

/*
func TestQuickSort(t *testing.T) {
	examples := []struct{
		name string
		want []wordProb
		slice []wordProb
		start int
		end int
	}{
		{
			name: "reverse",
			want: []wordProb{
				{"the", 15.0, 0.0},
				{"quick", 3.0, 1.0},
				{"brown", 5.8, 2.0},
				{"fox", 4.0, 3.0},
				{"jumped", 3.7, 4.0},
				{"over", 0.04, 5.0},
				},
			slice: []wordProb{
				{"over", 0.04, 5.0},
				{"jumped", 3.7, 4.0},
				{"fox", 4.0, 3.0},
				{"brown", 5.8, 2.0},
				{"quick", 3.0, 1.0},
				{"the", 15.0, 0.0},
				},
			start: 0,
			end: 5,
		},
	}
	for _, ex := range examples {
		t.Run(ex.name, func(t *testing.T) {
			quickSort(ex.slice, ex.start, ex.end)
			fmt.Println(ex.slice)
			for i, _ := range ex.slice {
				if ex.slice[i] != ex.want[i] {
					t.Fatalf("got %v, want %v ", ex.slice[i], ex.want[i])
				}
			}
		})
	}
}
*/
