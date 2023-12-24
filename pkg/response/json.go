package response

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

var EmptyJSONBody = struct{}{}

func JSON(w http.ResponseWriter, httpStatusCode int, jsonData any) {
	AddCORSHeaders(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)

	if err := json.NewEncoder(w).Encode(jsonData); err != nil {
		log.Err(err).Type("json_data_type", jsonData).Msg("Failed to marshal output JSON")
	}
}

func JSONData(w http.ResponseWriter, httpStatusCode int, data []byte) {
	AddCORSHeaders(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	_, _ = w.Write(data)
}
