package main

import (
	"app/internal/operations"
	"fmt"
)

func main() {
	fmt.Println("Welcome to ns1!!! Goodbye fool!!")
	fmt.Println("RES:", operations.AddGoo(1, 2))

	//http.ListenAndServe("localhost:8080", nil)
}
