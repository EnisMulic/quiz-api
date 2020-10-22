package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/EnisMulic/quiz-api/handlers"
	"github.com/go-openapi/runtime/middleware"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// create a logger
	logger := log.New(os.Stdout, "quiz-api", log.LstdFlags)

	userHandler := handlers.NewUser(logger)
	quizHandler := handlers.NewQuiz(logger)

	// create a new serve mux
	serverMux := mux.NewRouter()
	//serverMux.Handle("/user", userHandler).Methods("GET")

	getRouter := serverMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/user", userHandler.GetUsers)
	getRouter.HandleFunc("/user/{id:[0-9]+}", userHandler.GetUser)

	getRouter.HandleFunc("/quiz", quizHandler.GetQuizes)
	getRouter.HandleFunc("/quiz/{id:[0-9]+}", quizHandler.GetQuiz)

	putRouter := serverMux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/user/{id:[0-9]+}", userHandler.UpdateUser)
	putRouter.Use(userHandler.MiddlewareValidateUser)

	putRouter.HandleFunc("/quiz/{id:[0-9]+}", quizHandler.UpdateQuiz)

	postRouter := serverMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/user", userHandler.AddUser)
	postRouter.Use(userHandler.MiddlewareValidateUser)

	postRouter.HandleFunc("/quiz", quizHandler.AddQuiz)

	deleteRouter := serverMux.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/user/{id:[0-9]+}", userHandler.DeleteUser)

	deleteRouter.HandleFunc("/quiz/{id:[0-9]+}", quizHandler.DeleteQuiz)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// CORS
	corsHandler := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	// create the server
	server := &http.Server{
		Addr:         ":9090",
		Handler:      corsHandler(serverMux),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// start the server
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
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
	ctx, err := context.WithTimeout(context.Background(), 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	server.Shutdown(ctx)
}
