package apiserver

import (
	"net/http"
	"time"

	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

// appHeaders is middleware that adds application headers
func appHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Server", "atlas-api")
		w.Header().Set("content-type", "application/json")
		h.ServeHTTP(w, r)
	})
}

func (s *server) routes(appLogger *zerolog.Logger) {

	// Setup middleware chain
	// Build middleware chains from base logger
	c := alice.New().Append(hlog.NewHandler(*appLogger))
	c = c.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	}))
	c = c.Append(hlog.RemoteAddrHandler("ip"))
	c = c.Append(hlog.UserAgentHandler("user_agent"))
	c = c.Append(hlog.RefererHandler("referer"))
	// c = c.Append(hlog.RequestIDHandler("req_id", "Request-Id"))
	c = c.Append(appHeaders)

	// API Routes
	s.router.Handler("GET", "/api/v1/link", c.ThenFunc(s.handleListLinks()))
	s.router.Handler("GET", "/api/v1/link/:linkpath", c.ThenFunc(s.handleGetLink()))
	s.router.Handler("PUT", "/api/v1/link", c.ThenFunc(s.handleUpdateLink()))
	s.router.Handler("POST", "/api/v1/link", c.ThenFunc(s.handleCreateLink()))
	s.router.Handler("DELETE", "/api/v1/link/:linkpath", c.ThenFunc(s.handleDeleteLink()))

}
