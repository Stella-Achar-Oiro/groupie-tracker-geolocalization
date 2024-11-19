package handlers

import (
	"net/http"
	"html/template"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	// Parse the HTML template from a file
	tmplPath := "templates/error.html"
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Define data to pass to the template
	data := struct {
		StatusCode  int
		Message     string
		Description string
	}{
		StatusCode:  statusCode,
		Message:     message,
		Description: http.StatusText(statusCode),
	}
	// Set the status code and render the template
	w.WriteHeader(statusCode)
	tmpl.Execute(w, data)
}
