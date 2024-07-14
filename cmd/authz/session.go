package main

import (
	"fmt"
	"net/http"
)

// Handlers to handle session ..

// Login - show link for mleow, bob
// redirect back to demo ..
func loginHandler(w http.ResponseWriter, r *http.Request) {

	// Settign login .. overrides
	q := r.URL.Query()
	if q.Has("ID") {
		// Set Session cookie ...
		http.SetCookie(w, &http.Cookie{
			Name:   "ID",
			Path:   "/",
			Value:  q.Get("ID"),
			MaxAge: 60 * 60, // 3600s == 1h
		})
		http.Redirect(w, r, "/demo/", http.StatusFound)
		return
	}

	// If get this far .. then default page .. ignore any cookie ..
	// DEBUG
	//spew.Dump(r.Cookies())
	//_, err := r.Cookie("ID")
	//if err != nil {
	//	if errors.Is(err, http.ErrNoCookie) {
	//		// Output link to login here ..
	//		fmt.Fprintf(w, "Choose Login users ...")
	//		return
	//	}
	//}

	s := `
<html>
	<body>
	<demoHandler>Demo Login</demoHandler>
	<div>
		<a href="/demo/login/?ID=mleow">Login mleow</a><br/>
		<a href="/demo/login/?ID=bob">Login bob</a><br/>
	</div>
	</body>
</html>
`
	fmt.Fprintf(w, s)
	return
}

// Show login screen or if logged in .. then documents able to be accessed ..
// Ask for owner approval .. pending ..
// See public doc ..

// Logout - clear session ..
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Unset Session cookie ...
	http.SetCookie(w, &http.Cookie{
		Name:   "ID",
		Path:   "/",
		Value:  "",
		MaxAge: -1, // Delete cookie now ..
	})
	http.Redirect(w, r, "/demo/login/", http.StatusSeeOther)
	return
}
