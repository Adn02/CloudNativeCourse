// Find the top K most common words in a text document.
// Input path: location of the document, K top words
// Output: Slice of top K words
// For this excercise, word is defined as characters separated by a whitespace

// Note: You should use `checkError` to handle potential errors.

package textproc

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func topWords(path string, K int) []WordCount {
	file, errors := os.Open(path)

	// By convention, errors are the last return value and have type error, a built-in interface
	checkError(errors)

	// Defer is used to ensure that a function call is performed later
	// in a program’s execution, usually for purposes of cleanup. defer
	// is often used where e.g. ensure and finally would be used in other languages.
	defer file.Close()

	// Create hashtable with word as the key and number of occurences as the value
	wordCount := make(map[string]int)

	// Bufio package to implement buffered reader for reeading files
	scanner := bufio.NewScanner(file)

	var wordList []string
	for scanner.Scan() {
		newLine := scanner.Text()        // Returns the entire line of text as a single string, needs to be broken up by word
		words := strings.Fields(newLine) // Splits newLine into individual words based on whitespace

		wordList = append(wordList, words...)

		if err := scanner.Err(); err != nil { // A nil value in the error position indicates that there was no error.
			checkError(err)
		}
	}

	for _, word := range wordList { // Iterates through each word in words returning the index and the word hence the _, first
		wordCount[word]++
	}

	var mostWords []WordCount
	for word, count := range wordCount {
		mostWords = append(mostWords, WordCount{Word: word, Count: count}) // Convert hash table back to slice
	}

	sortWordCounts(mostWords) // Utilize helper function to sort slice in descending order

	if K <= len(mostWords) { // Print the top K most common words in a text document
		return mostWords[:K]
	}
	return mostWords
}

//--------------- DO NOT MODIFY----------------!

// A struct that represents how many times a word is observed in a document
type WordCount struct {
	Word  string
	Count int
}

// Method to convert struct to string format
func (wc WordCount) String() string {
	return fmt.Sprintf("%v: %v", wc.Word, wc.Count)
}

// Helper function to sort a list of word counts in place.
// This sorts by the count in decreasing order, breaking ties using the word.

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

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
