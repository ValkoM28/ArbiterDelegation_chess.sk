package handlers

import (
	"net/http"
)

// Serve the main HTML page
func FrontendHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}
