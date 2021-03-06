package apiserver

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/regalias/atlas-api/cache"
	"github.com/regalias/atlas-api/models"
	"github.com/regalias/atlas-api/util"
)

// Read

func (s *server) handleGetLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse link id
		linkPath := httprouter.ParamsFromContext(r.Context()).ByName("linkpath")

		// hlog.FromRequest(r).Debug().Msg("Requested link: " + linkPath)

		m, err := s.dataProvider.GetLinkDetails(linkPath)
		if err.Error() == "NotFound" {
			util.SendGenericResponse(w, r, "NotFound", http.StatusText(404), 404)
			return
		} else if err != nil {
			util.ThrowISE(w, r)
			return
		}
		util.SendGenericResponse(w, r, "None", m, 200)
	}
}

func (s *server) handleListLinks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		util.SendGenericResponse(w, r, "None", "ok", http.StatusNotImplemented)
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

		// guid := xid.New()

		newLink := &models.LinkModel{
			// LinkID:         guid.String(),
			LinkPath:       req.LinkPath,
			CanonicalName:  req.CanonicalName,
			TargetURL:      req.TargetURL,
			CreatedTime:    time.Now().Unix(),
			LastModified:   time.Now().Unix(),
			LastModifiedBy: "some-user", // TODO
			Enabled:        req.Enabled,
		}

		if err := s.dataProvider.CreateLink(newLink); err != nil {
			if err.Error() == "AlreadyExists" {
				util.SendGenericResponse(w, r, "ParameterError", "Specfied LinkPath is already in use", 400)
			} else {
				s.logger.Error().Str("Error", err.Error()).Msg("Could not insert new entry")
				util.ThrowISE(w, r)
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

		util.SendGenericResponse(w, r, "None", resp, http.StatusCreated)
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

		newLink := &models.LinkModel{
			// LinkID:         req.LinkID,
			CanonicalName:  req.CanonicalName,
			LinkPath:       req.LinkPath,
			TargetURL:      req.TargetURL,
			LastModified:   time.Now().Unix(),
			LastModifiedBy: "some-user", // TODO
			Enabled:        req.Enabled,
		}

		if err := s.dataProvider.UpdateLink(newLink); err != nil {
			if err.Error() == "NotFound" {
				util.SendGenericResponse(w, r, "NotFound", http.StatusText(404), 404)
				return
			} else if err.Error() == "NoChange" {
				util.SendGenericResponse(w, r, "None", http.StatusText(http.StatusNotModified), http.StatusNotModified)
			} else if err != nil {
				util.ThrowISE(w, r)
				return
			}
		}

		// TODO: update redis cache! (if exists?)
		if err := s.cacheTaskHandler.SubmitTask(&cache.Task{
			Operation: cache.SetLink,
			Linkpath:  req.LinkPath,
			Linkdest:  req.TargetURL,
		}); err != nil {
			s.logger.Error().Msg("Couldn't submit cache set task: " + err.Error())
			util.ThrowISE(w, r)
			return
		}

		util.SendGenericResponse(w, r, "None", req, http.StatusOK)
	}
}

func (s *server) handleDeleteLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse link id
		linkPath := httprouter.ParamsFromContext(r.Context()).ByName("linkpath")

		// hlog.FromRequest(r).Debug().Msg("Requested link: " + linkPath)

		err := s.dataProvider.DeleteLink(linkPath)
		if err != nil {
			if err.Error() == "NotFound" {
				util.SendGenericResponse(w, r, "NotFound", "Resource not found", 404)
				return
			} else if err != nil {
				util.ThrowISE(w, r)
				return
			}
		}

		// TODO: Purge the redis cache!
		if err := s.cacheTaskHandler.SubmitTask(&cache.Task{
			Operation: cache.RemoveLink,
			Linkpath:  linkPath,
		}); err != nil {
			s.logger.Error().Msg("Couldn't submit cache deletion task: " + err.Error())
			util.ThrowISE(w, r)
			return
		}

		util.SendGenericResponse(w, r, "None", http.StatusText(http.StatusOK), 200)
	}
}
