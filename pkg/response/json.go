package response

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

var EmptyJSONBody = struct{}{}

func JSON(w http.ResponseWriter, httpStatusCode int, jsonData interface{}) {
	AddCORSHeaders(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)

	if err := json.NewEncoder(w).Encode(jsonData); err != nil {
		log.Err(err).Msgf("Failed to marshal output JSON from %T", jsonData)
	}
}

func JSONData(w http.ResponseWriter, httpStatusCode int, data []byte) {
	AddCORSHeaders(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	_, _ = w.Write(data)
}
