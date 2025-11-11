package handlers

import "net/http"

// ServeDocFile serves swagger documentation file
func ServeDocFile(w http.ResponseWriter, r *http.Request) {
	filePath := "docs/swagger.json"

	// Set headers if you want the browser to download the file
	w.Header().Set("Content-Disposition", "attachment; filename=swagger.json")
	w.Header().Set("Content-Type", "application/octet-stream")

	http.ServeFile(w, r, filePath)
}
