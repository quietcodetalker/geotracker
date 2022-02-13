package util

import (
	"encoding/json"
	"net/http"
)

// DecodeBody decodes request body.
func DecodeBody(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

// EncodeBody encodes response body.
func EncodeBody(w http.ResponseWriter, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

// Respond set response header and encodes response body.
func Respond(w http.ResponseWriter, status int, data interface{}) error {
	w.WriteHeader(status)
	if data != nil {
		return EncodeBody(w, data)
	}
	return nil
}

// RespondErr builds response body with provided error and responds.
func RespondErr(w http.ResponseWriter, status int, err error) error {
	w.WriteHeader(status)
	return EncodeBody(
		w,
		map[string]interface{}{
			"error": err,
		},
	)
}
