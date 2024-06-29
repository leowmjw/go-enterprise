package authz

import (
	"github.com/davecgh/go-spew/spew"
	openfga "github.com/openfga/go-sdk"
	"os"
)

func CheckPermission() {

	apiURL := os.Getenv("FGA_API_URL")
	println(apiURL)

	configuration, err := openfga.NewConfiguration(openfga.Configuration{
		//ApiUrl: "http://localhost", // required, e.g. https://api.fga.example
		ApiUrl: apiURL,
	})

	spew.Dump(configuration)

	if err != nil {
		// .. Handle error
		panic(err)
	}
	
}

// Add Permissions ...

// Remove Permissions ...

// List down the
