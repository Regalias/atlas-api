package apiserver

import (
	"fmt"
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

	type requestModel struct {
		CanonicalName string `json:"canonicalName" validate:"required,min=3,max=50,alphanumunicode"`
		URI           string `json:"URI" validate:"required,min=3,max=50,is-uri"`
		TargetURL     string `json:"targetURL" validate:"required,min=3,max=500,url"`
	}

	type responseModel struct {
		LinkID        string `json:"linkID"`
		CanonicalName string `json:"canonicalName"`
		URI           string `json:"URI"`
		TargetURL     string `json:"targetURL"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		var req requestModel
		err := s.getRequest(w, r, &req)
		if err != nil {
			return
		}

		// Debug: remove
		fmt.Printf("\n%+v\n", req)

		resp := &responseModel{
			LinkID:        "abcd",
			CanonicalName: req.CanonicalName,
			URI:           req.URI,
			TargetURL:     req.TargetURL,
		}

		sendGenericResponse(w, r, "None", resp, 200)
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
