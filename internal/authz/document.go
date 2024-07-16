package authz

type Document struct {
	ID      string
	Owner   string
	Content string
	// WF Lifecycle??
}

var as AuthStore

type AuthzDemo struct {
	as AuthStore
}

// NewAuthzDemo to start workflow ..
func NewAuthzDemo(apiURL, policyPath string) AuthzDemo {
	// Load Policy model ..
	as := NewAuthStore(apiURL)
	// is it here??
	//err := as.DemoDirectAccess()
	//if err != nil {
	//	panic(err)
	//}
	// Connect Temporal?
	return AuthzDemo{as}
}

func init() {
	//apiURL := os.Getenv("FGA_API_URL")
	//as = NewAuthStore(apiURL)
	//err := as.InitDemo("")
	//if err != nil {
	//	fmt.Println("Error initializing auth store. ERR:", err.Error())
	//} else {
	//	fmt.Println("SUCCESS!! Init Authz Server!!")
	//}

	// Setup the Docs ... inside an Org Authz Workflow?? If not started yet ..

	// For each doc; owner will get the write access ..
	// All these within the workflow ..
}

func NewDocument(id, owner, content string) Document {
	return Document{}
}

func (d Document) StartWorkflow() {

}
