package apiserver

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog/hlog"
)

// No longer used
// func isURL(str string) bool {
// 	u, err := url.Parse(str)
// 	return err == nil && u.Scheme != "" && u.Host != ""
// }

func sendGenericResponse(w http.ResponseWriter, r *http.Request, errMsg string, details interface{}, code int) {

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

// getRequest takes in an arbitrary struct, attempts to read the request, marshal the request into the struct, and perform validation
func (s *server) getRequest(w http.ResponseWriter, r *http.Request, model interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		io.WriteString(w, http.StatusText(400))
		return err
	}

	if err := json.Unmarshal(body, model); err != nil {
		sendGenericResponse(w, r, http.StatusText(400), "None", http.StatusBadRequest)
		return err
	}

	errMsg, err := s.validateModel(model)
	if err != nil {
		sendGenericResponse(w, r, "InvalidParameters", errMsg, http.StatusBadRequest)
		return err
	}
	return nil
}

func throwISE(w http.ResponseWriter, r *http.Request) {
	sendGenericResponse(w, r, http.StatusText(http.StatusInternalServerError), "None", http.StatusInternalServerError)
	return
}
