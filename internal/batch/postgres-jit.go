package batch

import "fmt"

func AddRole(username, role string) error {
	fmt.Println("AddRole: ", role, " for user: ", username)
	
	return nil
}
