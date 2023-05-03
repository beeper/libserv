package response

import (
	"net/http"
)

func AddCORSHeaders(w http.ResponseWriter) {
	// Recommended CORS headers can be found in https://spec.matrix.org/v1.3/client-server-api/#web-browser-clients
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization")
}
