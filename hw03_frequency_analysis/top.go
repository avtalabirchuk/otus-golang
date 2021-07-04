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

	wcList := make([]wordsCount, len(freq))
	i := 0
	for key, val := range freq {
		wcList[i] = wordsCount{key, val}
		i++
	}

	sort.Slice(wcList, func(i, j int) bool {
		return wcList[i].count > wcList[j].count
	})

	fmt.Println(wcList)
	result := make([]string, 0)
	for i, v := range wcList {
		result = append(result, v.word)
		if i > 10 {
			break
		}
	}
	fmt.Println(result)
	return result
}
