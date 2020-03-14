package apiserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"

	// dataprovider "github.com/regalias/atlas-api/apiserver/providers"
	"github.com/rs/zerolog"
)

type server struct {
	router           *httprouter.Router
	validator        *validator.Validate
	logger           *zerolog.Logger
	http             *http.Server
	dataProvider     *DataProvider
	cacheTaskHandler *cacheTaskHandler
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

	r := httprouter.New()
	d, err := NewDataProvider(lgr, "atlas-table-main")
	if err != nil {
		lgr.Fatal().Str("Error", err.Error()).Msg("Could not initialize database provider")
	}
	// TODO: grab table name from config
	if err := d.ensureTable(); err != nil {
		lgr.Fatal().Str("Error", err.Error()).Msg("DDB table was not found and could not create required table")
	}

	// Create cache and async task handler
	lgr.Info().Msg("Starting cache worker...")
	c, err := newCache("127.0.0.1", 6379, nil)
	if err != nil {
		lgr.Fatal().Msg(err.Error())
	}
	tq := newTaskQueue(100, lgr, c)
	go tq.runWorker() // Start the worker
	lgr.Info().Msg("Cache worker started")

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
		dataProvider:     d,
		cacheTaskHandler: tq,
	}

	s.routes(lgr)

	lgr.Info().Msg("Atlas API server starting...")
	if err := s.http.ListenAndServe(); err != nil {
		lgr.Fatal().Err(err).Msg("API Startup failed")
	}

	return 0
}
