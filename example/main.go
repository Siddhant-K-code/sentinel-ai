package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// unusedFunction is an example of dead code
func unusedFunction() {
	fmt.Println("This function is never called")
}

// vulnerableHandler has a potential XSS vulnerability
func vulnerableHandler(w http.ResponseWriter, r *http.Request) {
	userInput := r.URL.Query().Get("name")

	// This is vulnerable to XSS - user input is not escaped
	tmpl := `<h1>Hello {{.Name}}!</h1>`
	t, _ := template.New("test").Parse(tmpl)
	t.Execute(w, map[string]interface{}{
		"Name": userInput, // Vulnerable: no escaping
	})
}

// safeHandler shows the correct way to handle user input
func safeHandler(w http.ResponseWriter, r *http.Request) {
	userInput := r.URL.Query().Get("name")

	// This is safe - user input is properly escaped
	tmpl := `<h1>Hello {{.Name}}!</h1>`
	t, _ := template.New("test").Parse(tmpl)
	t.Execute(w, map[string]interface{}{
		"Name": template.HTMLEscapeString(userInput), // Safe: properly escaped
	})
}

func main() {
	http.HandleFunc("/vulnerable", vulnerableHandler)
	http.HandleFunc("/safe", safeHandler)

	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8080", nil)
}
