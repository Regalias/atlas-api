package apiserver

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

type server struct {
	router    *httprouter.Router
	validator *validator.Validate
	logger    *zerolog.Logger
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
	appLogger := *lgr

	// Create server context object
	s := server{
		router:    httprouter.New(),
		validator: registerValidators(),
		logger:    &appLogger,
	}

	s.routes(&appLogger)

	appLogger.Info().Msg("Atlas API server starting...")
	if err := http.ListenAndServe(":8080", s.router); err != nil {
		appLogger.Fatal().Err(err).Msg("API Startup failed")
	}

	return 0
}
