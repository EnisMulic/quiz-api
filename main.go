package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/EnisMulic/quiz-api/db"

	"github.com/EnisMulic/quiz-api/handlers"
	"github.com/go-openapi/runtime/middleware"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// setup mongodb connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, connErr := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	defer func() {
		if connErr = client.Disconnect(ctx); connErr != nil {
			panic(connErr)
		}
	}()

	// create a logger
	logger := log.New(os.Stdout, "quiz-api", log.LstdFlags)

	// create repositories
	userRepo := db.NewUserRepository(client)
	quizRepo := db.NewQuizRepository(client)

	// create handlers
	userHandler := handlers.NewUser(logger, userRepo)
	quizHandler := handlers.NewQuiz(logger, quizRepo)
	authHandler := handlers.NewAuth(logger, userRepo)

	// create a new serve mux
	serverMux := mux.NewRouter()

	getRouter := serverMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/user", userHandler.GetUsers)
	getRouter.HandleFunc("/user/{id}", userHandler.GetUser)

	getRouter.HandleFunc("/quiz", quizHandler.GetQuizes)
	getRouter.HandleFunc("/quiz/{id}", quizHandler.GetQuiz)

	putRouter := serverMux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/user/{id}", userHandler.UpdateUser)
	// putRouter.Use(userHandler.MiddlewareValidateUser)

	putRouter.HandleFunc("/quiz/{id}", quizHandler.UpdateQuiz)

	postRouter := serverMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/user", userHandler.AddUser)
	// postRouter.Use(userHandler.MiddlewareValidateUser)

	postRouter.HandleFunc("/quiz", quizHandler.AddQuiz)

	postRouter.HandleFunc("/quiz/{id}/question", quizHandler.AddQuestion)

	deleteRouter := serverMux.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/user/{id}", userHandler.DeleteUser)

	deleteRouter.HandleFunc("/quiz/{id}", quizHandler.DeleteQuiz)

	authRouter := serverMux.Methods(http.MethodPost).Subrouter()
	authRouter.HandleFunc("/auth/register", authHandler.Register)
	authRouter.HandleFunc("/auth/login", authHandler.Login)

	// add auth middleware
	getRouter.Use(handlers.IsAuthorized)
	postRouter.Use(handlers.IsAuthorized)
	putRouter.Use(handlers.IsAuthorized)
	deleteRouter.Use(handlers.IsAuthorized)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	docsRouter := serverMux.Methods(http.MethodGet).Subrouter()
	docsRouter.Handle("/docs", sh)
	docsRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// CORS
	corsHandler := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	// create the server
	server := &http.Server{
		Addr:         ":8080",
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
