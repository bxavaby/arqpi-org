package search

import (
	"sort"
	"strings"
	"unicode"

	"github.com/bxavaby/arqpi-org/internal/models"
)

type SearchIndex struct {
	fragments []models.Fragment
	tokenMap  map[string][]int
}

func NewSearchIndex(fragments []models.Fragment) *SearchIndex {
	idx := &SearchIndex{
		fragments: fragments,
		tokenMap:  make(map[string][]int),
	}
	idx.buildIndex()
	return idx
}

func (idx *SearchIndex) buildIndex() {
	for _, fragment := range idx.fragments {
		titleTokens := tokenize(fragment.Title)
		for _, token := range titleTokens {
			for range 3 {
				idx.tokenMap[token] = append(idx.tokenMap[token], fragment.ID)
			}
		}

		textTokens := tokenize(fragment.Text)
		for _, token := range textTokens {
			idx.tokenMap[token] = append(idx.tokenMap[token], fragment.ID)
		}
	}
}

type SearchResult struct {
	Fragment models.Fragment
	Score    int
}

func (idx *SearchIndex) Search(query string, limit int) []models.Fragment {
	if limit <= 0 {
		limit = 10
	}

	queryTokens := tokenize(query)
	if len(queryTokens) == 0 {
		return []models.Fragment{}
	}

	matchCounts := make(map[int]int)
	for _, token := range queryTokens {
		for _, fragID := range idx.tokenMap[token] {
			matchCounts[fragID]++
		}
	}

	results := make([]SearchResult, 0, len(matchCounts))
	for fragID, score := range matchCounts {
		for _, frag := range idx.fragments {
			if frag.ID == fragID {
				results = append(results, SearchResult{
					Fragment: frag,
					Score:    score,
				})
				break
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Score != results[j].Score {
			return results[i].Score > results[j].Score
		}

		return results[i].Fragment.Length < results[j].Fragment.Length
	})

	fragments := make([]models.Fragment, 0, limit)
	for i := 0; i < len(results) && i < limit; i++ {
		fragments = append(fragments, results[i].Fragment)
	}

	return fragments
}

func tokenize(text string) []string {
	text = strings.ToLower(text)

	words := strings.FieldsFunc(text, func(r rune) bool {
		return !isWordChar(r)
	})

	var result []string
	for _, word := range words {
		word = strings.TrimSpace(word)

		if len(word) <= 2 || isStopWord(word) {
			continue
		}

		result = append(result, word)
	}

	return result
}

func isWordChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) ||
		r == 'ç' || r == 'á' || r == 'à' || r == 'ã' || r == 'â' ||
		r == 'é' || r == 'ê' || r == 'í' || r == 'ó' || r == 'ô' ||
		r == 'õ' || r == 'ú' || r == 'ü'
}

func isStopWord(word string) bool {
	ptStopwords := map[string]bool{
		"a": true, "o": true, "e": true, "de": true, "da": true, "do": true,
		"em": true, "um": true, "uma": true, "que": true, "com": true, "no": true,
		"na": true, "por": true, "para": true, "os": true, "as": true, "dos": true,
		"das": true, "ao": true, "não": true, "mas": true, "se": true,
	}

	enStopwords := map[string]bool{
		"a": true, "an": true, "the": true, "and": true, "or": true, "but": true,
		"of": true, "in": true, "on": true, "at": true, "to": true, "for": true,
		"by": true, "with": true, "as": true, "is": true, "are": true, "was": true,
		"were": true, "be": true, "this": true, "that": true, "it": true,
	}

	return ptStopwords[word] || enStopwords[word]
}
