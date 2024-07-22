package authz

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	openfga "github.com/openfga/go-sdk"
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

	id := ""
	// Create store if needed ..
	gsresp, gserr := fgaClient.ListStores(context.Background()).Execute()
	if gserr != nil {
		//spew.Dump(gserr)
		fmt.Println("ERR: ", gserr.Error())
		panic(gserr)
	}
	if len(gsresp.GetStores()) == 0 {
		resp, cerr := fgaClient.CreateStore(context.Background()).Body(
			ClientCreateStoreRequest{
				Name: "Demo",
			}).Execute()

		if cerr != nil {
			panic(cerr)
		}
		id = resp.GetId()

	} else {
		store := gsresp.GetStores()[0]
		id = store.GetId()
		fmt.Println("Existing STORE!!! ==> ", store.GetName())
	}

	// Update the client with the StoreID; no need storeID in struct i think ..
	serr := fgaClient.SetStoreId(id)
	if serr != nil {
		panic(serr)
	}

	return AuthStore{
		client:  fgaClient,
		storeID: id,
	}
}

func (a AuthStore) addTuple(body ClientWriteTuplesBody) error {
	wopts := ClientWriteOptions{}
	wresp, werr := a.client.WriteTuples(context.Background()).
		Body(body).Options(wopts).Execute()
	if werr != nil {
		// can ignore existnig ...
		fmt.Println("ERR: ", werr.Error())
		//return werr
	}
	// DEBUG
	fmt.Println("WRITE: ", len(wresp.Writes))
	//spew.Dump(wresp.Writes)
	//spew.Dump(wresp.Deletes)
	return nil
}

func (a AuthStore) removeTuple(body ClientDeleteTuplesBody) error {
	dopts := ClientWriteOptions{}
	dresp, derr := a.client.DeleteTuples(context.Background()).
		Body(body).Options(dopts).Execute()
	if derr != nil {
		fmt.Println("ERR: ", derr.Error())
		return derr
	}
	// DEBUG
	fmt.Println("DELETED: ", len(dresp.Deletes))
	//spew.Dump(dresp.Writes)
	//spew.Dump(dresp.Deletes)
	return nil
}

func (a AuthStore) hasAccess(user, relation, document string) (bool, error) {
	// Opts empty; uses the latest model ..
	opts := ClientCheckOptions{}
	data, cerr := a.client.Check(context.Background()).Body(ClientCheckRequest{
		User:     "user:" + user,
		Relation: relation,
		Object:   "document:" + document,
		//Context:          nil,
		//ContextualTuples: []ClientTupleKey{}, // Like dynamic stuff .. MFA clicked ..
	}).Options(opts).Execute()
	// Any unexpected view ..
	if cerr != nil {
		fmt.Println("ERR: ", cerr.Error())
		return false, cerr
	}
	// Chck if allowed and not a nil ..
	allowed, ok := data.GetAllowedOk()
	if ok {
		if *allowed {
			fmt.Println("User: ", user, " allowed to view Doc:", document)
			return true, nil
		}
	}
	// Default no access
	return false, nil
}

func (a AuthStore) CanViewDocument(user, document string) (bool, error) {
	return a.hasAccess(user, "viewer", document)
}

func (a AuthStore) CanEditDocument(user, document string) (bool, error) {
	return a.hasAccess(user, "editor", document)
}

func (a AuthStore) AddViewRelationship(user, document string) error {
	// TODO: What further valdiations??
	// This can add conditions ..
	t := []ClientTupleKey{
		{
			User:     "user:" + user,
			Relation: "viewer",
			Object:   "document:" + document,
		},
	}
	return a.addTuple(t)
}

func (a AuthStore) RemoveViewRelationship(user, document string) error {
	// TODO: What further valdiations??
	t := []ClientTupleKeyWithoutCondition{
		{
			User:     "user:" + user,
			Relation: "viewer",
			Object:   "document:" + document,
		},
	}
	return a.removeTuple(t)
}

func (a AuthStore) InitDemo(demoModelPath string) error {
	// CReate new Store .. store it for later ..
	// dEBUzg
	//spew.Dump(a.storeID)

	// This gets all the models
	gamresp, gerr := a.client.ReadAuthorizationModels(context.Background()).Execute()
	if gerr != nil {
		panic(gerr)
	}
	for i, m := range gamresp.GetAuthorizationModels() {
		fmt.Println("ID: ", i, " MODEL: ", m.GetId())
	}
	// For future .. when need to cintune ..
	//fmt.Println("TOKEN: ", gamresp.GetContinuationToken())

	// Probably need to set the client options ..
	//opts := ClientReadAuthorizationModelOptions{
	//	AuthorizationModelId: nil,
	//}
	// without above; it pulls the latest only ..

	//a.client.ReadAuthorizationModel(context.Background()).
	//	Body(ClientReadAuthorizationModelRequest{}).
	//	Options(opts).
	//	Execute()

	return nil
}

// ===================================================================================
//
//	<<<<<<<<<<<<<<<<<<  OLD CODE BELOW   >>>>>>>>>>>>>>>>>>>
//
// ===================================================================================
func (a AuthStore) DemoPrepareModel(demoModelPath string) error {
	// Read the special JSON model .. it will crap out if not proper JSON!!
	b, err := os.ReadFile(demoModelPath)
	if err != nil {
		return err
	}
	var body ClientWriteAuthorizationModelRequest
	uerr := json.Unmarshal(b, &body)
	if uerr != nil {
		return uerr
	}
	data, werr := a.client.WriteAuthorizationModel(context.Background()).Body(body).Execute()

	if werr != nil {
		return werr
	}
	// Set the client to the laets one ..
	modelID := data.GetAuthorizationModelId()
	serr := a.client.SetAuthorizationModelId(modelID)
	if serr != nil {
		return serr
	}
	// Store the modelID if needed ...
	a.modelID = modelID
	// All OK ..
	return nil

}

func (a AuthStore) checkAccess(modelID string) (*ClientCheckResponse, error) {
	fmt.Println("Model ID: ", modelID)

	opts := ClientCheckOptions{
		AuthorizationModelId: openfga.PtrString(modelID),
	}
	data, cerr := a.client.Check(context.Background()).Body(ClientCheckRequest{
		User:     "user:mleow",
		Relation: "viewer",
		Object:   "document:public/welcome.doc",
		//Context:          nil,
		//ContextualTuples: []ClientTupleKey{}, // Like dynamic stuff .. MFA clicked ..
	}).Options(opts).Execute()

	if cerr != nil {
		return nil, cerr
	}
	return data, nil
}

func (a AuthStore) DemoDirectAccess() error {
	// Store the model pointer ...
	err := a.DemoPrepareModel("openfga/models/direct-access.json")
	if err != nil {
		return err
	}
	modelID, gerr := a.client.GetAuthorizationModelId()
	if gerr != nil {
		return gerr
	}
	data, err := a.checkAccess(modelID)

	allowed, ok := data.GetAllowedOk()
	if ok {
		if !*allowed {
			fmt.Println("As expected .. not allowed .. yet!!")
		} else {
			return fmt.Errorf("DemoDirectAccess: UNEXPECTED ALLOWED!!! ")
		}
	} else {
		fmt.Println("Expected no OK")
	}

	wopts := ClientWriteOptions{
		AuthorizationModelId: openfga.PtrString(modelID),
	}
	wresp, werr := a.client.WriteTuples(context.Background()).
		Body(ClientWriteTuplesBody{
			{
				User:     "user:mleow",
				Relation: "viewer",
				Object:   "document:public/welcome.doc",
			},
			{
				User:     "user:bob",
				Relation: "editor",
				Object:   "document:public/welcome.doc",
			},
		}).
		Options(wopts).Execute()
	if werr != nil {
		return werr
	}
	// DEBUG
	fmt.Println("No of writes", len(wresp.Writes))
	fmt.Println("No of deletes", len(wresp.Deletes))
	//spew.Dump(wresp.Writes)

	// Now chevk again is fine ..
	data, err = a.checkAccess(modelID)

	allowed, ok = data.GetAllowedOk()
	if ok {
		if *allowed {
			fmt.Println("As expected .. now is ok!!")
		} else {
			return fmt.Errorf("DemoDirectAccess: UNEXPECTED NOT ALLOWED!!! ")
		}
	} else {
		return fmt.Errorf("Expected OK; but NOT!!")
	}

	return nil
}

func (a AuthStore) DemoAccess() error {

	// Now with contextual data ..
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
