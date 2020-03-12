package apiserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/xid"
	"github.com/rs/zerolog/hlog"
)

// Read

func (s *server) handleGetLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse link
		linkID := httprouter.ParamsFromContext(r.Context()).ByName("linkid")

		// TODO: query db for link details
		hlog.FromRequest(r).Debug().Msg("Requested link: " + linkID)

		m, err := s.dataProvider.GetLinkDetails(linkID)
		if err == redis.Nil {
			sendGenericResponse(w, r, "NotFound", http.StatusText(404), 404)
			return
		} else if err != nil {
			throwISE(w, r)
			return
		}
		sendGenericResponse(w, r, "None", m, 200)
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
		LinkPath      string `json:"linkPath" validate:"required,min=3,max=50,is-uri-path"`
		TargetURL     string `json:"targetURL" validate:"required,min=3,max=500,url"`
	}

	type responseModel struct {
		LinkID        string `json:"linkID"`
		CanonicalName string `json:"canonicalName"`
		LinkPath      string `json:"linkPath"`
		TargetURL     string `json:"targetURL"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		var req requestModel
		if err := s.getRequest(w, r, &req); err != nil {
			return
		}

		// Debug: remove
		fmt.Printf("\n%+v\n", req)

		guid := xid.New()

		newLink := &LinkModel{
			LinkID:         guid.String(),
			CanonicalName:  req.CanonicalName,
			LinkPath:       req.LinkPath,
			TargetURL:      req.TargetURL,
			Created:        time.Now().Unix(),
			LastModified:   time.Now().Unix(),
			LastModifiedBy: "some-user", // TODO
		}
		// Debug
		fmt.Printf("\n%+v\n", newLink)

		if err := s.dataProvider.CreateLink(newLink); err != nil {
			s.logger.Error().Str("Error", err.Error()).Msg("Could not insert new entry")
			// sendGenericResponse(w, r, http.StatusText(http.StatusInternalServerError), "None", http.StatusInternalServerError)
			throwISE(w, r)
			return
		}

		resp := &responseModel{
			LinkID:        guid.String(),
			CanonicalName: req.CanonicalName,
			LinkPath:      req.LinkPath,
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
