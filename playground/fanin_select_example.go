package main

import (
	"fmt"
	"math/rand"
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
			// when either a or b has a value, the case statment will be enacted
			// and the value put onto c
			select {
			case val := <-a:
				c <- val
			case val := <-b:
				c <- val
			}
		}
	}()
	return c
}

func boring(msg string) <-chan string { // returns receive-only chan
	c := make(chan string)
	random := rand.Intn(6400) + 1

	go func() { // launching go routine from inside the function
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(random) * time.Millisecond)
		}
	}()
	return c // return channel to caller
}
