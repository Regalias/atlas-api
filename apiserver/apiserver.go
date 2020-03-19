package apiserver

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"

	// dataprovider "github.com/regalias/atlas-api/apiserver/providers"
	"github.com/rs/zerolog"

	"github.com/regalias/atlas-api/cache"
	"github.com/regalias/atlas-api/database"
	"github.com/regalias/atlas-api/logging"

	"github.com/regalias/atlas-api/util"
)

type server struct {
	router           *httprouter.Router
	validator        *validator.Validate
	logger           *zerolog.Logger
	http             *http.Server
	dataProvider     database.Provider
	cacheTaskHandler *cache.AsyncHandler
}

// Run does magic things
func Run(args []string) int {

	// TODO: parse run args

	// Create logger
	lgr, err := logging.New("debug", "atlas-api", true)
	if err != nil {
		fmt.Printf("Oh noes! Something went horribly wrong!")
		panic(err)
	}

	r := httprouter.New()
	d, err := database.NewDDB(lgr, "atlas-table-main")
	if err != nil {
		lgr.Fatal().Str("Error", err.Error()).Msg("Could not initialize database provider")
	}
	// TODO: grab table name from config
	if err := d.InitDatabase(); err != nil {
		lgr.Fatal().Str("Error", err.Error()).Msg("Database or table was not found and could not create required resources")
	}

	// Create cache and async task handler
	lgr.Info().Msg("Starting cache worker...")
	c, err := cache.NewRedisProvider("127.0.0.1", 6379, nil)
	if err != nil {
		lgr.Fatal().Msg(err.Error())
	}

	tq := cache.NewAsyncQueue(100, lgr, c)
	go tq.RunWorker() // Start the worker
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
			Addr:              ":8081",
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

// getRequest takes in an arbitrary struct, attempts to read the request, marshal the request into the struct, and perform validation
func (s *server) getRequest(w http.ResponseWriter, r *http.Request, model interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		io.WriteString(w, http.StatusText(400))
		return err
	}

	if err := json.Unmarshal(body, model); err != nil {
		util.SendGenericResponse(w, r, http.StatusText(400), "None", http.StatusBadRequest)
		return err
	}

	errMsg, err := s.validateModel(model)
	if err != nil {
		util.SendGenericResponse(w, r, "InvalidParameters", errMsg, http.StatusBadRequest)
		return err
	}
	return nil
}
