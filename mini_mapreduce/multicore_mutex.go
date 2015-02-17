package main

import (
	"bufio"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type wordCount struct {
	word  string
	count int
}

// This program reads words from an input file and provides a count for
// how many times each word appeared.
func main() {
	log.Println("Beginning Multi core mutex...")
	start := time.Now()
	// runtime.GOMAXPROCS(1)

	// 1. read in input from file
	out := readFromFile("input.txt")
	// 2. count words
	res1 := countWords(out)
	res2 := countWords(out)

	// 3. combine word counts to produce result and print
	for m := range merge(res1, res2) {
		for k, v := range m {
			log.Printf("%s %d", k, v)
		}
	}

	elapsed := time.Since(start)
	log.Printf("Multi core mutex took %s", elapsed)
}

// merge takes in a list of channels and create a master map of each word and their
// combined count. It then returns this map
func merge(cs ...<-chan wordCount) <-chan map[string]int {
	out := make(chan map[string]int)

	var wg sync.WaitGroup
	wg.Add(len(cs))

	// m contains a unique list of all words and their total count
	m := make(map[string]int)
	var mutex = &sync.Mutex{}

	mapIt := func(counted <-chan wordCount) {
		for wc := range counted {
			mutex.Lock()
			if count, ok := m[wc.word]; ok {
				m[wc.word] = count + wc.count
			} else {
				m[wc.word] = wc.count
			}
			mutex.Unlock()
		}
		wg.Done()
	}

	// for each channel, update the master map m, using mutexes to ensure
	// data integrity
	for _, counted := range cs {
		go mapIt(counted)
	}

	// when all channels have updated the master map, we can close the channel and return
	go func() {
		wg.Wait()
		out <- m
		close(out)
	}()
	return out
}

// countWords reads from a channel of strings and counts the number of
// times each word appears. It then puts each unique wordCount value onto a channel.
func countWords(in <-chan string) <-chan wordCount {
	log.Println("Beginning counting words... ")
	out := make(chan wordCount)

	go func() {
		m := make(map[string]wordCount)

		for word := range in {
			// log.Printf("word %s %v", word, r)
			var inc = 1
			if wc, ok := m[word]; ok {
				inc = wc.count + 1
			}
			val := wordCount{
				word:  word,
				count: inc,
			}
			m[word] = val
		}

		// put everything from map on out channel
		for _, v := range m {
			out <- v
		}
		close(out)
	}()
	return out
}

// readFromFile reads each line from the given filename and splits it out by whitespace.
// Each new word is added onto the out channel
func readFromFile(filename string) <-chan string {
	out := make(chan string)

	go func() {
		file, _ := os.Open(filename)
		defer file.Close()
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := scanner.Text()
			for _, word := range strings.Fields(line) {
				out <- word
			}
		}
		close(out)
	}()

	// if err := scanner.Err(); err != nil {
	// 	log.Fatal(err)
	// }

	return out
}
