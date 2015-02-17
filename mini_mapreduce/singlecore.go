package main

import (
	"bufio"
	// "io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	// 1. read contents from file
	// 2. count words
	// 3. print results
	start := time.Now()

	words := readFromFile("input.txt")
	m := countWords(words)

	for k, v := range m {
		log.Printf("%s: %d\n", k, v)
	}

	elapsed := time.Since(start)
	log.Printf("Single core took %s", elapsed)
}

func readFromFile(filename string) []string {
	var result []string
	file, _ := os.Open(filename)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		result = append(result, strings.Fields(line)...)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return result
}

func countWords(words []string) map[string]int {
	m := make(map[string]int)

	for _, word := range words {
		if count, ok := m[word]; ok {
			count++
			m[word] = count
		} else {
			m[word] = 1
		}
	}
	return m
}
