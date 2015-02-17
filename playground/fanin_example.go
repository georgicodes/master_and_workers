package main

import (
	"fmt"
	"time"
)

func main() {
	c := fanIn(boring("dean"), boring("georgi"))
	for i := 0; i < 10; i++ {
		fmt.Println(<-c)
	}

	fmt.Println("Im leaving")
}

func fanIn(a, b <-chan string) <-chan string {
	c := make(chan string)
	go func() {
		for {
			c <- <-a // receive value from a and send to chan c
		}
	}()
	go func() {
		for {
			c <- <-b
		}
	}()
	return c
}

func boring(msg string) <-chan string { // returns receive-only chan
	c := make(chan string)

	go func() { // launching go routine from inside the function
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(1))
		}
	}()
	return c // return channel to caller
}
