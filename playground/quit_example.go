package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	quit := make(chan bool)
	c := boring("dean", quit)
	for i := 0; i < 10; i++ {
		fmt.Println(<-c)
	}

	fmt.Println("Im quitting")
	quit <- true
}

func boring(msg string, quit chan bool) <-chan string { // returns receive-only chan
	c := make(chan string)
	random := rand.Intn(6400) + 1

	// go func() { // launching go routine from inside the function
	// 	for i := 0; ; i++ {
	// 		c <- fmt.Sprintf("%s %d", msg, i)
	// 		time.Sleep(time.Duration(random) * time.Millisecond)
	// 	}
	// }()

	select {
	case c <- fmt.Sprintf("%s", msg):
	case <-quit:
		return nil
	}
	return c // return channel to caller
}
