package cos418_hw1_1

import (
	"os"
	"regexp"
	"strings"
	"sort"
	"fmt"
)

// Find the top K most common words in a text document.
//
//	path: location of the document
//	numWords: number of words to return (i.e. k)
//	charThreshold: character threshold for whether a token qualifies as a word,
//		e.g. charThreshold = 5 means "apple" is a word but "pear" is not.
//
// Matching is case insensitive, e.g. "Orange" and "orange" is considered the same word.
// A word comprises alphanumeric characters only. All punctuation and other characters
// are removed, e.g. "don't" becomes "dont".
// You should use `checkError` to handle potential errors.
func topWords(path string, numWords int, charThreshold int) []WordCount {
	// TODO: implement me
	bytes, err := os.ReadFile(path)
	checkError(err)
	// HINT: You may find the `strings.Fields` and `strings.ToLower` functions helpful
	// HINT: To keep only alphanumeric characters, use the regex "[^0-9a-zA-Z]+"

	// strip non-alphanumeric characters, lowercase, split into tokens
	re := regexp.MustCompile("[^0-9a-zA-Z]+")
	cleaned := re.ReplaceAllString(string(bytes), " ")
	tokens := strings.Fields(strings.ToLower(cleaned))


	// count frequencies skip words below charThreshold
	counts := make(map[string]int)
	for _, word := range tokens {
		if len(word) >= charThreshold {
			counts[word]++
		}
	}

	// build slice of WordCount
	wcs := make([]WordCount, 0, len(counts))
	for word, count := range counts {
		wcs = append(wcs, WordCount{Word: word, Count: count})
	}

	// Sort and return top K
	sortWordCounts(wcs)
	if numWords > len(wcs) {
		numWords = len(wcs)
	}
	return wcs[:numWords]


}

// A struct that represents how many times a word is observed in a document
type WordCount struct {
	Word  string
	Count int
}

func (wc WordCount) String() string {
	return fmt.Sprintf("%v: %v", wc.Word, wc.Count)
}

// Helper function to sort a list of word counts in place.
// This sorts by the count in decreasing order, breaking ties using the word.
// DO NOT MODIFY THIS FUNCTION!
func sortWordCounts(wordCounts []WordCount) {
	sort.Slice(wordCounts, func(i, j int) bool {
		wc1 := wordCounts[i]
		wc2 := wordCounts[j]
		if wc1.Count == wc2.Count {
			return wc1.Word < wc2.Word
		}
		return wc1.Count > wc2.Count
	})
}
