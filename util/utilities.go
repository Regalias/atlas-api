package util

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/rs/zerolog/hlog"
)

// No longer used
// func isURL(str string) bool {
// 	u, err := url.Parse(str)
// 	return err == nil && u.Scheme != "" && u.Host != ""
// }

// SendGenericResponse sends a HTTP response with the specified code, and result encoded in a defined JSON format
func SendGenericResponse(w http.ResponseWriter, r *http.Request, errMsg string, details interface{}, code int) {

	type errorResponse struct {
		Error   string      `json:"error"`
		Details interface{} `json:"details"`
	}

	resp, err := json.Marshal(errorResponse{
		Error:   errMsg,
		Details: details,
	})
	if err != nil {
		// We really shouldn't get here... throw a 500 ISE
		hlog.FromRequest(r).Error().Err(err).Msg("")
		io.WriteString(w, http.StatusText(500))
	} else {
		w.WriteHeader(code)
		w.Write(resp)
	}
	return
}

// ThrowISE is a helper function that returns a generic 500 ISE response
func ThrowISE(w http.ResponseWriter, r *http.Request) {
	SendGenericResponse(w, r, http.StatusText(http.StatusInternalServerError), "None", http.StatusInternalServerError)
	return
}
