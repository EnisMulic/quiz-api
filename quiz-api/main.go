package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/EnisMulic/quiz-api/quiz-api/config"
	"github.com/EnisMulic/quiz-api/quiz-api/db"

	"github.com/EnisMulic/quiz-api/quiz-api/handlers"
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

	dbHost := config.GetEnvVariable("QUIZ_DB")
	client, connErr := mongo.Connect(ctx, options.Client().ApplyURI(dbHost))

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

	// user routers
	userGetRouter := serverMux.Methods(http.MethodGet).Subrouter()
	userGetRouter.HandleFunc("/user", userHandler.GetUsers)
	userGetRouter.HandleFunc("/user/{id}", userHandler.GetUser)
	userGetRouter.Use(handlers.IsAuthorized)

	userPutRouter := serverMux.Methods(http.MethodPut).Subrouter()
	userPutRouter.HandleFunc("/user/{id}", userHandler.UpdateUser)
	userPutRouter.Use(userHandler.MiddlewareValidateUser)
	userPutRouter.Use(handlers.IsAuthorized)

	userPostRouter := serverMux.Methods(http.MethodPost).Subrouter()
	userPostRouter.HandleFunc("/user", userHandler.AddUser)
	userPostRouter.Use(handlers.IsAuthorized)
	userPostRouter.Use(userHandler.MiddlewareValidateUser)

	userDeleteRouter := serverMux.Methods(http.MethodDelete).Subrouter()
	userDeleteRouter.HandleFunc("/user/{id}", userHandler.DeleteUser)
	userDeleteRouter.Use(handlers.IsAuthorized)

	// quiz routers
	quizGetRouter := serverMux.Methods(http.MethodGet).Subrouter()
	quizGetRouter.HandleFunc("/quiz", quizHandler.GetQuizes)
	quizGetRouter.HandleFunc("/quiz/{id}", quizHandler.GetQuiz)
	quizGetRouter.Use(handlers.IsAuthorized)

	quizPutRouter := serverMux.Methods(http.MethodPut).Subrouter()
	quizPutRouter.HandleFunc("/quiz/{id}", quizHandler.UpdateQuiz)
	quizPutRouter.Use(handlers.IsAuthorized)
	quizPutRouter.Use(quizHandler.MiddlewareValidateQuiz)

	quizPostRouter := serverMux.Methods(http.MethodPost).Subrouter()
	quizPostRouter.HandleFunc("/quiz", quizHandler.AddQuiz)
	quizPostRouter.Use(handlers.IsAuthorized)
	quizPostRouter.Use(quizHandler.MiddlewareValidateQuiz)

	quizDeleteRouter := serverMux.Methods(http.MethodDelete).Subrouter()
	quizDeleteRouter.HandleFunc("/quiz/{id}", quizHandler.DeleteQuiz)
	quizDeleteRouter.Use(handlers.IsAuthorized)

	// question router
	questionPostRouter := serverMux.Methods(http.MethodPost).Subrouter()
	questionPostRouter.HandleFunc("/quiz/{id}/question", quizHandler.AddQuestion)
	questionPostRouter.Use(handlers.IsAuthorized)
	questionPostRouter.Use(quizHandler.MiddlewareValidateQuestion)

	questionDeleteRouter := serverMux.Methods(http.MethodDelete).Subrouter()
	questionDeleteRouter.HandleFunc("/quiz/{id}/question/{question_id}", quizHandler.DeleteQuestion)
	questionDeleteRouter.Use(handlers.IsAuthorized)

	// authentication router
	authRouter := serverMux.Methods(http.MethodPost).Subrouter()
	authRouter.HandleFunc("/auth/register", authHandler.Register)
	authRouter.HandleFunc("/auth/login", authHandler.Login)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	docsRouter := serverMux.Methods(http.MethodGet).Subrouter()
	docsRouter.Handle("/docs", sh)
	docsRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// CORS
	corsHandler := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

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
