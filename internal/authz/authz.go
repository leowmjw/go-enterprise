package authz

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	. "github.com/openfga/go-sdk/client"
	"os"
)

type AuthStore struct {
	client           *OpenFgaClient
	storeID, modelID string
}

// NewAuthStore loads store .. or if not create it ..
func NewAuthStore(apiURL string) AuthStore {
	fgaClient, err := NewSdkClient(&ClientConfiguration{
		ApiUrl: apiURL,
	})

	if err != nil {
		panic(err)
	}

	// Create store if needed ..
	gsresp, gserr := fgaClient.GetStore(context.Background()).Execute()
	if gserr != nil {
		panic(gserr)
	}

	id := ""
	ptrid, ok := gsresp.GetIdOk()
	if !ok {
		resp, cerr := fgaClient.CreateStore(context.Background()).Body(
			ClientCreateStoreRequest{
				Name: "Demo",
			}).Execute()

		if cerr != nil {
			panic(cerr)
		}
		id = resp.GetId()
	} else {
		fmt.Println("Existing STORE!!! ==> ", gsresp.GetName())
		id = *ptrid
	}

	return AuthStore{
		client:  fgaClient,
		storeID: id,
	}
}

func (a AuthStore) InitDemo(demoModelPath string) error {
	// CReate new Store .. store it for later ..
	return nil
}

func (a AuthStore) DemoDirectAccess(demoModelPath string) error {
	// CReate a new model ...
	// Store the model pointer ...

	return nil
}

func (a AuthStore) DemoUserGroups(demoModelPath string) error {
	return nil
}

func (a AuthStore) DemoRolesPermissions(demoModelPath string) error {
	return nil
}

func (a AuthStore) DemoUserParentChild(demoModelPath string) error {
	return nil
}

func (a AuthStore) DemoConditions(demoModelPath string) error {
	return nil
}

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
		panic(err)
	}

	// Check via query if it is allowed againsrt the docId
	resp, err := fgaClient.Check(context.Background()).Body(ClientCheckRequest{
		User:     "user:",
		Relation: "reader",
		Object:   "document:",
	}).Execute()
	if err != nil {
		panic(err)
	}
	b, ok := resp.GetAllowedOk()
	if !ok {
		fmt.Println("UNKNOWN OK .. check!! ====> ")
	} else {
		spew.Dump(b)
	}
}

// Add Permissions ...

func AddDocAccess(username, docId string) {
	apiURL := os.Getenv("FGA_API_URL")
	fgaClient, err := NewSdkClient(&ClientConfiguration{
		ApiUrl:  apiURL,
		StoreId: "01J1H266C1JPHAB8X1FV1M43D7",
	})

	if err != nil {
		panic(err)
	}

	// Example:
	//	user:mleow
	//	document:secret-topsecrety

	data, err := fgaClient.Write(context.Background()).Body(ClientWriteRequest{
		Writes: []ClientTupleKey{
			{
				User:     "user:" + username,
				Relation: "reader",
				Object:   "document:" + docId,
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

func RemoveDocAccess(username, docId string) {
	apiURL := os.Getenv("FGA_API_URL")
	fgaClient, err := NewSdkClient(&ClientConfiguration{
		ApiUrl:  apiURL,
		StoreId: "01J1H266C1JPHAB8X1FV1M43D7",
	})

	if err != nil {
		panic(err)
	}

	// Example:
	//	user:mleow
	//	document:secret-topsecrety

	data, err := fgaClient.Write(context.Background()).Body(ClientWriteRequest{
		Deletes: []ClientTupleKeyWithoutCondition{
			{
				"user:" + username,
				"reader",
				"document:" + docId,
			},
		},
	}).Execute()
	if err != nil {
		panic(err)
	}

	// DEBUG
	spew.Dump(data.Deletes)

}

// List down the Doc related to user

func ListDocAccess(username string) {
	apiURL := os.Getenv("FGA_API_URL")
	fgaClient, err := NewSdkClient(&ClientConfiguration{
		ApiUrl: apiURL,
	})
	if err != nil {
		panic(err)
	}
	data, err := fgaClient.ListObjects(context.Background()).Body(ClientListObjectsRequest{
		User:     "user:" + username,
		Relation: "reader",
		Type:     "document",
	}).Execute()

	if err != nil {
		panic(err)
	}

	s, ok := data.GetObjectsOk()
	if !ok {
		fmt.Println("UNKNOWN .. NOT SET!!!! =========>")
		spew.Dump(data)
	} else {
		spew.Dump(s)
	}

}
