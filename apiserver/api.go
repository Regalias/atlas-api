package apiserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/hlog"
)

// Read

func (s *server) handleGetLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse link id
		linkPath := httprouter.ParamsFromContext(r.Context()).ByName("linkpath")

		hlog.FromRequest(r).Debug().Msg("Requested link: " + linkPath)

		m, err := s.dataProvider.GetLinkDetails(linkPath)
		if err.Error() == "NotFound" {
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
		sendGenericResponse(w, r, "None", "ok", http.StatusNotImplemented)
	}
}

// Write

func (s *server) handleCreateLink() http.HandlerFunc {

	type requestResponseModel struct {
		LinkPath      string `json:"LinkPath" validate:"required,min=3,max=50,is-uri-path"`
		CanonicalName string `json:"CanonicalName" validate:"required,min=3,max=50,alphanumunicode"`
		TargetURL     string `json:"TargetURL" validate:"required,min=3,max=500,url"`
		Enabled       bool   `json:"Enabled" validate:"omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		var req requestResponseModel
		if err := s.getRequest(w, r, &req); err != nil {
			return
		}

		// Debug: remove
		fmt.Printf("\n%+v\n", req)

		// guid := xid.New()

		newLink := &LinkModel{
			// LinkID:         guid.String(),
			LinkPath:       req.LinkPath,
			CanonicalName:  req.CanonicalName,
			TargetURL:      req.TargetURL,
			CreatedTime:    time.Now().Unix(),
			LastModified:   time.Now().Unix(),
			LastModifiedBy: "some-user", // TODO
			Enabled:        req.Enabled,
		}
		// Debug
		// fmt.Printf("\n%+v\n", newLink)

		if err := s.dataProvider.CreateLink(newLink); err != nil {
			if err.Error() == "AlreadyExists" {
				sendGenericResponse(w, r, "ParameterError", "Specfied LinkPath is already in use", 400)
			} else {
				s.logger.Error().Str("Error", err.Error()).Msg("Could not insert new entry")
				// sendGenericResponse(w, r, http.StatusText(http.StatusInternalServerError), "None", http.StatusInternalServerError)
				throwISE(w, r)
			}
			return
		}

		// TODO: also push to redis cache server?

		resp := &requestResponseModel{
			// LinkID:        guid.String(),
			CanonicalName: req.CanonicalName,
			LinkPath:      req.LinkPath,
			TargetURL:     req.TargetURL,
			Enabled:       req.Enabled,
		}

		sendGenericResponse(w, r, "None", resp, http.StatusCreated)
	}
}

func (s *server) handleUpdateLink() http.HandlerFunc {
	type requestResponseModel struct {
		// LinkID        string `json:"LinkID" validate:"required,min=3,max=50"`
		CanonicalName string `json:"CanonicalName" validate:"required,min=3,max=50,alphanumunicode"`
		LinkPath      string `json:"LinkPath" validate:"required,min=3,max=50,is-uri-path"`
		TargetURL     string `json:"TargetURL" validate:"required,min=3,max=500,url"`
		Enabled       bool   `json:"Enabled"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		var req requestResponseModel
		if err := s.getRequest(w, r, &req); err != nil {
			return
		}

		newLink := &LinkModel{
			// LinkID:         req.LinkID,
			CanonicalName:  req.CanonicalName,
			LinkPath:       req.LinkPath,
			TargetURL:      req.TargetURL,
			LastModified:   time.Now().Unix(),
			LastModifiedBy: "some-user", // TODO
			Enabled:        req.Enabled,
			CreatedTime:    1337,
			// Don't modify creation timestamp!
		}

		if err := s.dataProvider.UpdateLink(newLink); err != nil {
			if err.Error() == "NotFound" {
				sendGenericResponse(w, r, "NotFound", http.StatusText(404), 404)
				return
			} else if err.Error() == "NoChange" {
				sendGenericResponse(w, r, "None", http.StatusText(http.StatusNotModified), http.StatusNotModified)
			} else if err != nil {
				throwISE(w, r)
				return
			}
		}
		// TODO: update redis cache!
		sendGenericResponse(w, r, "None", req, http.StatusOK)
	}
}

func (s *server) handleDeleteLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse link id
		linkPath := httprouter.ParamsFromContext(r.Context()).ByName("linkpath")

		hlog.FromRequest(r).Debug().Msg("Requested link: " + linkPath)

		err := s.dataProvider.DeleteLink(linkPath)
		if err != nil {
			if err.Error() == "NotFound" {
				sendGenericResponse(w, r, "NotFound", "Resource not found", 404)
				return
			} else if err != nil {
				throwISE(w, r)
				return
			}
		}

		// TODO: Purge the redis cache!
		sendGenericResponse(w, r, "None", http.StatusText(http.StatusOK), 200)
	}
}
