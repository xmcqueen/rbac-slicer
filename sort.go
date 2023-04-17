package main

import (
 "sort"
)

// https://stackoverflow.com/questions/18695346/how-can-i-sort-a-mapstringint-by-its-values
type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func rankByWordCount(wordFrequencies map[string]int) (rv []string) {
	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	for _, v := range pl {
		rv = append(rv, v.Key)
	}
	return
}

// the above Pair thing is better
// make both of these use a Pair some day
func sortKeys(dat map[string]int) (rv []string) {
	for k := range dat {
		rv = append(rv, k)
	}
	sort.Strings(rv)
	return
}
