package apiserver

import (
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
	return func(w http.ResponseWriter, r *http.Request) {
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
