package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/EnisMulic/quiz-api/quiz-images/config"
	"github.com/EnisMulic/quiz-api/quiz-images/files"
	"github.com/EnisMulic/quiz-api/quiz-images/handlers"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	hclog "github.com/hashicorp/go-hclog"
)

func main() {

	logger := hclog.New(
		&hclog.LoggerOptions{
			Name:  "quiz-images",
			Level: hclog.LevelFromString("debug"),
		},
	)

	// create a logger
	serverLogger := logger.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true})

	// create the storage class, use local storage
	// max filesize 5MB
	stor, err := files.NewLocal("./imagestore", 1024*1000*5)
	if err != nil {
		logger.Error("Unable to create storage", "error", err)
		os.Exit(1)
	}

	// create the handlers
	fileHandler := handlers.NewFiles(stor, logger)

	// create a new serve mux
	serverMux := mux.NewRouter()

	// CORS
	corsHandler := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	// upload files
	postHandler := serverMux.Methods(http.MethodPost).Subrouter()
	postHandler.HandleFunc("/images/{id}/{filename:[a-zA-Z]+\\.[a-z]{3,}}", fileHandler.Upload)

	// get files
	getHandler := serverMux.Methods(http.MethodGet).Subrouter()
	getHandler.Handle(
		"/images/{id}/{filename:[a-zA-Z]+\\.[a-z]{3,}}",
		http.StripPrefix("/images/", http.FileServer(http.Dir("./imagestore"))),
	)
	getHandler.Use(handlers.GzipMiddleware)

	// create the server
	addr := config.GetEnvVariable("API_ADDRESS")
	server := &http.Server{
		Addr:         addr,
		Handler:      corsHandler(serverMux),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// start the server
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			serverLogger.Fatal(err)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	// block until a signal is received
	signal := <-signalChannel
	log.Println("Got signal:", signal)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, shutdownErr := context.WithTimeout(context.Background(), 30*time.Second)
	if shutdownErr != nil {
		log.Fatal(err)
	}

	server.Shutdown(ctx)
}
