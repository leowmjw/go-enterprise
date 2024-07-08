package main

import (
	"app/internal/authz"
	"fmt"
	"os"
)

func main() {

	fmt.Println("OpenFGA .. ")
	//authz.CheckPermission()

	apiURL := os.Getenv("FGA_API_URL")
	as := authz.NewAuthStore(apiURL)
	err := as.InitDemo("")
	if err != nil {
		fmt.Println("Error initializing auth store. ERR:", err.Error())
	}

	err = as.DemoDirectAccess()
	if err != nil {
		fmt.Println("Error showing Demo for Direct Access!! ERR:", err.Error())
	}
}
