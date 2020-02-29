package apiserver

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/hlog"
)

func isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func sendGenericResponse(w http.ResponseWriter, r *http.Request, errMsg string, details string, code int) {
	resp, err := json.Marshal(genericResponse{
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
