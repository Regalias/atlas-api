package apiserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

type server struct {
	router    *httprouter.Router
	validator *validator.Validate
	logger    *zerolog.Logger
	http      *http.Server
}

// Run does magic things
func Run(args []string) int {

	// TODO: parse run args

	// Create logger
	lgr, err := initLogger("debug")
	if err != nil {
		fmt.Printf("Oh noes! Something went horribly wrong!")
		panic(err)
	}
	// appLogger := *lgr

	r := httprouter.New()

	// Create server context struct
	s := server{
		router:    r,
		validator: newValidator(),
		logger:    lgr,
		http: &http.Server{
			ReadHeaderTimeout: 20 * time.Second,
			ReadTimeout:       1 * time.Minute,
			WriteTimeout:      2 * time.Minute,
			Addr:              ":8080",
			Handler:           r,
		},
	}

	s.routes(lgr)

	lgr.Info().Msg("Atlas API server starting...")
	if err := s.http.ListenAndServe(); err != nil {
		lgr.Fatal().Err(err).Msg("API Startup failed")
	}

	return 0
}
