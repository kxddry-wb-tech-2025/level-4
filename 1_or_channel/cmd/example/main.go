package main

import (
	"fmt"
	"time"

	"or_channel/pkg/or"
)

// Helper function to create a channel that closes after a duration
func sigAfter(name string, d time.Duration) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		time.Sleep(d)
		fmt.Printf("%s done\n", name)
	}()
	return ch
}

// Example usage of the or package
func main() {
	fmt.Println("Waiting for the first signal...")

	sig1 := sigAfter("A", 2*time.Second)
	sig2 := sigAfter("B", 1*time.Second)
	sig3 := sigAfter("C", 3*time.Second)

	start := time.Now()
	<-or.Or(sig1, sig2, sig3)
	fmt.Printf("First signal received after %v\n", time.Since(start))
}
