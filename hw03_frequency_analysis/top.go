package hw03frequencyanalysis

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

type wordsCount struct {
	word  string
	count int
}

type PairList []wordsCount

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].count < p[j].count }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func Top10(s string) []string {
	// s = strings.ToLower(s)
	freq := make(map[string]int)

	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c) && !unicode.IsPunct(c)
	}
	a := strings.FieldsFunc(s, f)
	for _, c := range a {
		freq[c]++
	}

	wcList := make(PairList, len(freq))
	i := 0
	for key, val := range freq {
		wcList[i] = wordsCount{key, val}
		i++
	}
	sort.Sort(sort.Reverse(wcList))

	// sort.Slice(wcList, func(i, j int) bool {
	// 	return wcList[i].count > wcList[j].count
	// })

	fmt.Println(wcList)
	result := make([]string, 0)
	for i, v := range wcList {
		if i < 10 {
			result = append(result, v.word)
		} else {
			break
		}
	}
	fmt.Println(result)
	return result
}
