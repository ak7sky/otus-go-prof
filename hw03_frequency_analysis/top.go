package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var sep = regexp.MustCompile(`[.,"':;\-!?]*\s+[.,"':;\-!?]*\s*`)

var toTrim = regexp.MustCompile(`\A[.,"':;\-!?\s]+|[.,"':;\-!?\s]+\z`)

type wFrq struct {
	word string
	frq  int
}

func Top10(input string) []string {
	if input == "" {
		return []string{}
	}

	input = toTrim.ReplaceAllString(input, "")
	words := sep.Split(strings.ToLower(input), -1)
	wfMap := make(map[string]wFrq)

	for _, w := range words {
		if wi, exist := wfMap[w]; exist {
			wi.frq++
			wfMap[w] = wi
		} else {
			wfMap[w] = wFrq{word: w, frq: 1}
		}
	}

	wfSlc := make([]wFrq, 0, len(wfMap))
	for _, wf := range wfMap {
		wfSlc = append(wfSlc, wf)
	}

	sort.Slice(wfSlc, func(i, j int) bool {
		if wfSlc[i].frq != wfSlc[j].frq {
			return wfSlc[i].frq > wfSlc[j].frq
		}
		return wfSlc[i].word < wfSlc[j].word
	})

	wfSlc = wfSlc[:10]
	top10 := make([]string, 0, len(wfSlc))
	for _, wf := range wfSlc {
		top10 = append(top10, wf.word)
	}

	return top10
}
