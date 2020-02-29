package apiserver

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/hlog"
)

// Read

func (s *server) handleGetLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse link
		slinkParam := httprouter.ParamsFromContext(r.Context()).ByName("slink")

		// TODO: query db for link details
		hlog.FromRequest(r).Debug().Msg("Requested link: " + slinkParam)
		sendGenericResponse(w, r, "None", "ok", 200)
	}
}

func (s *server) handleListLinks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sendGenericResponse(w, r, "None", "ok", 200)
	}
}

// Write

func (s *server) handleCreateLink() http.HandlerFunc {

	type createLinkRequest struct {
		CanonicalName string `json:"canonicalName" validate:"required,min=3,max=50,alphanumunicode"`
		URI           string `json:"URI" validate:"required,min=3,max=50,is-uri"`
		TargetURL     string `json:"targetURL" validate:"required,min=3,max=500,url"`
	}

	type response struct {
		LinkID        string `json:"linkID"`
		CanonicalName string `json:"canonicalName"`
		URI           string `json:"URI"`
		TargetURL     string `json:"targetURL"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			io.WriteString(w, http.StatusText(400))
			return
		}

		var rq *createLinkRequest
		if err := json.Unmarshal(body, &rq); err != nil {
			sendGenericResponse(w, r, http.StatusText(400), "None", http.StatusBadRequest)
			return
		}

		errMsg, err := s.validateModel(rq)
		if err != nil {
			sendErrorResponse(w, r, "InvalidParameters", errMsg, http.StatusBadRequest)
			return
		}

		// Debug: remove
		fmt.Printf("\n%+v\n", rq)

		sendGenericResponse(w, r, "None", "ok", 200)
	}
}

func (s *server) handleUpdateLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sendGenericResponse(w, r, "None", "ok", 200)
	}
}

func (s *server) handleDeleteLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sendGenericResponse(w, r, "None", "ok", 200)
	}
}
