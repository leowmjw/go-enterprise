package main

import (
	"app/internal/authz"
	"fmt"
)

func main() {

	fmt.Println("OpenFGA .. ")
	authz.CheckPermission()

}
