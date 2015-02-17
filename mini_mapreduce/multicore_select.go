package main

import (
	"bufio"
	"log"
	"os"
	"runtime"
	"strings"
	// "sync"
	"time"
)

type wordCount struct {
	word  string
	count int
}

func main() {
	// 1. read contents from file
	// 2. count words
	// 3. print results
	log.Println("Beginning Multi core select...")
	start := time.Now()
	runtime.GOMAXPROCS(2)
	// log.Printf("Running program with %d processes", maxProcs)

	out := readFromFile("input.txt")
	res1 := countWords(out)
	res2 := countWords(out)

	var done sync.WaitGroup
	done.Add(2)

	go func() {
		log.Println("am i in here")
		var m = make(map[string]int)
		for {
			select {
			case wc, closed := <-res1:
				if closed {
					log.Printf("something on res1")
				}
				if count, ok := m[wc.word]; ok {
					m[wc.word] = count + wc.count
				} else {
					m[wc.word] = wc.count
				}
			case wc := <-res2:
				if wc == nil {
					break
				}
				log.Printf("something on res2 %v", wc.word)
				if count, ok := m[wc.word]; ok {
					m[wc.word] = count + wc.count
				} else {
					m[wc.word] = wc.count
				}
			}
		}

		for k, v := range m {
			log.Printf("%s %d", k, v)
		}
	}()

	elapsed := time.Since(start)
	log.Printf("Multi core select took %s", elapsed)
	<-done
}

func countWords(in <-chan string) <-chan wordCount {
	log.Println("Beginning counting words...")
	out := make(chan wordCount)

	go func() {
		m := make(map[string]wordCount)

		for word := range in {
			// log.Printf("%s", word)
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
