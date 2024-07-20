package authz

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	. "github.com/openfga/go-sdk/client"
	"os"
)

type Document struct {
	ID      string
	Owner   string
	Content string
	// WF Lifecycle??
}

var as AuthStore

type AuthzDemo struct {
	as               AuthStore
	users            []string
	docs             []Document
	awaitingApproval []string // WorkflowID for Owner-Docs requested ..
}

// NewAuthzDemo to start workflow ..
func NewAuthzDemo(apiURL, policyPath string) (AuthzDemo, error) {
	// Start with reasonable defaults first ..
	if apiURL == "" {
		apiURL = os.Getenv("FGA_API_URL")
	}
	if policyPath == "" {
		policyPath = "openfga/models/direct-access.json"
	}
	// Load Store
	as := NewAuthStore(apiURL)
	// Load the model ..
	err := as.DemoPrepareModel(policyPath)
	if err != nil {
		fmt.Println("Failed to prepare model!! ERR:", err)
		return AuthzDemo{}, err
	}

	return AuthzDemo{as: as}, nil
}

func (ad AuthzDemo) setupTuples() error {
	// Start with a very naive impleemntation ..
	keys := make([]ClientTupleKey, 0)
	for _, doc := range ad.docs {
		fmt.Print("DocPath:", doc.ID, " Owner:", doc.Owner)
		// For each doc; set owner as viewer + editor
		// For docs with public; as viewer anyone?
	}
	// DEBUG
	spew.Dump(ClientWriteTuplesBody(keys))
	return nil
}
func (ad AuthzDemo) debugState() {
	spew.Dump(ad.docs)
	spew.Dump(ad.awaitingApproval)
	return
}

func (ad AuthzDemo) checkViewerAccess(user, document string) bool {
	//modelID, gerr := as.client.GetAuthorizationModelId()
	//if gerr != nil {
	//	fmt.Println("ERR:", gerr)
	//	return false
	//}
	//fmt.Println("Model ID: ", modelID)
	opts := ClientCheckOptions{}
	// Below not needed as it takes the latest modelID for store ..
	//opts := ClientCheckOptions{
	//	AuthorizationModelId: openfga.PtrString(modelID),
	//}

	data, cerr := ad.as.client.Check(context.Background()).Body(ClientCheckRequest{
		User:     "user:" + user,
		Relation: "viewer",
		Object:   "document:" + document,
		//Context:          nil,
		//ContextualTuples: []ClientTupleKey{}, // Like dynamic stuff .. MFA clicked ..
	}).Options(opts).Execute()

	if cerr != nil {
		fmt.Println("ERR:", cerr)
		return false
	}
	allowed, ok := data.GetAllowedOk()
	if ok && *allowed {
		// is OK
		return true
	}
	// Default deny ..
	return false
}

func (ad AuthzDemo) getDocumentContent(user, document string) (string, error) {
	// Naive implementation first ..
	for _, doc := range ad.docs {
		if doc.ID == document {
			if ad.checkViewerAccess(user, document) {
				return doc.Content, nil

			}
			return "", fmt.Errorf("user %s unauthorized viewer of %s", user, document)
		}
	}
	return "", fmt.Errorf("document not found")
}

func (d Document) checkEditorAccess(user, document string) bool {
	// Example check access
	//modelID, gerr := as.client.GetAuthorizationModelId()
	//if gerr != nil {
	//	return AuthzDemo{}, gerr
	//}
	//data, err := as.checkAccess(modelID)
	//allowed, ok := data.GetAllowedOk()
	//if ok {
	//	if *allowed {
	//		// is OK
	//	}
	//} else {
	//	// Default deny ..
	//}
	return false
}
