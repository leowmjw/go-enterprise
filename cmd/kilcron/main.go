package main

import "fmt"

func main() {
	fmt.Println("kilcron")
	Run()
}

func Run() {
	// Start the workers listening to the taskqueue ..
	// with policy of 2s timeout and retry 1 time ..
	// Concurrent call to 10 calls between 1-6
}
