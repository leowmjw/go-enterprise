package authz

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"os"
	"time"

	. "github.com/openfga/go-sdk/client"
)

func CheckPermission() {

	apiURL := os.Getenv("FGA_API_URL")
	//println(apiURL)
	//
	//configuration, err := openfga.NewConfiguration(openfga.Configuration{
	//	//ApiUrl: "http://localhost", // required, e.g. https://api.fga.example
	//	ApiUrl: apiURL,
	//})
	//
	//spew.Dump(configuration)
	//
	//if err != nil {
	//	// .. Handle error
	//	panic(err)
	//}

	fgaClient, err := NewSdkClient(&ClientConfiguration{
		ApiUrl:  apiURL, // required, e.g. https://api.fga.example
		StoreId: "01J1H266C1JPHAB8X1FV1M43D7",
		//StoreId:              os.Getenv("FGA_STORE_ID"), // optional, not needed for \`CreateStore\` and \`ListStores\`, required before calling for all other methods
		//AuthorizationModelId: os.Getenv("FGA_MODEL_ID"), // Optional, can be overridden per request
	})

	if err != nil {
		// .. Handle error
	}

	// DEBUG ..
	//conf := fgaClient.GetConfig()
	// DEBUG
	//spew.Dump(conf)

	// Example of creating Store .. it does not check existing store name ...
	//resp, err := fgaClient.CreateStore(context.Background()).Body(ClientCreateStoreRequest{Name: "FGA Demo"}).Execute()
	//if err != nil {
	//	// .. Handle error
	//}
	//
	//spew.Dump(resp.Name)

	// Write tuple ..
	body := ClientWriteRequest{
		Writes: []ClientTupleKey{
			{
				User:     "user:mleow",
				Relation: "reader",
				Object:   "document:secretz",
			},
		},
	}

	data, err := fgaClient.Write(context.Background()).
		Body(body).
		//Options(options).
		Execute()

	if err != nil {
		// .. Handle error check validation
		//errors.Is(err, &openfga.ValidationErrorMessageResponse)
		fmt.Println(err.Error())
		//spew.Dump(err)
		//panic(err)
		return
	}

	// Small delay ..
	time.Sleep(time.Second)
	fmt.Println("I am slow ...")

	data, err = fgaClient.Write(context.Background()).Body(ClientWriteRequest{
		Deletes: []ClientTupleKeyWithoutCondition{
			{"user:mleow", "reader", "document:secrety"},
		},
	}).Execute()

	spew.Dump(data.Writes)
	spew.Dump(data.Deletes)

}

// Add Permissions ...
func AddSecretDocAccess(username string) {
	apiURL := os.Getenv("FGA_API_URL")
	fgaClient, err := NewSdkClient(&ClientConfiguration{
		ApiUrl:  apiURL,
		StoreId: "01J1H266C1JPHAB8X1FV1M43D7",
	})

	if err != nil {
		panic(err)
	}

	// Example: user:mleow
	data, err := fgaClient.Write(context.Background()).Body(ClientWriteRequest{
		Writes: []ClientTupleKey{
			{
				User:     "user:" + username,
				Relation: "reader",
				Object:   "document:secret-topsecrety",
			},
		},
	}).Execute()
	if err != nil {
		panic(err)
	}
	// DEBUG
	spew.Dump(data.Writes)
}

// Remove Permissions ...
func RemoveSecretDocAccess(username string) {
	apiURL := os.Getenv("FGA_API_URL")
	fgaClient, err := NewSdkClient(&ClientConfiguration{
		ApiUrl:  apiURL,
		StoreId: "01J1H266C1JPHAB8X1FV1M43D7",
	})

	if err != nil {
		panic(err)
	}

	// Example: user:mleow
	data, err := fgaClient.Write(context.Background()).Body(ClientWriteRequest{
		Deletes: []ClientTupleKeyWithoutCondition{
			{
				"user:" + username,
				"reader",
				"document:secret-topsecrety",
			},
		},
	}).Execute()
	if err != nil {
		panic(err)
	}

	spew.Dump(data.Writes)
	spew.Dump(data.Deletes)

}

// List down the
