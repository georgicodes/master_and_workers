package main

import (
	"fmt"
	"time"
)

func main() {
	c := boring("boring!")
	for i := 0; i < 5; i++ {
		fmt.Printf("You say %q\n", <-c)
	}

	fmt.Println("Im leaving")
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
