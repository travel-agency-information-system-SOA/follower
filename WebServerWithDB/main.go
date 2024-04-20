package main

import (
	"context"
	handlers "database-example/handler"
	"database-example/model"
	repository "database-example/repo"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logger := log.New(os.Stdout, "[follower-api] ", log.LstdFlags)
	storeLogger := log.New(os.Stdout, "[follower-store] ", log.LstdFlags)

	store, err := repository.New(storeLogger)
	if err != nil {
		logger.Fatal(err)
	}
	defer store.CloseDriverConnection(timeoutContext)
	store.CheckConnection()

	user := &model.User{
		ID:       1,
		Username: "milica",
	}
	user2 := &model.User{
		ID:       2,
		Username: "kristina",
	}
	user3 := &model.User{
		ID:       3,
		Username: "ana",
	}
	err = store.CreateUser(user3)
	if err != nil {
		logger.Fatal("Error creating User:", err)
		return
	}
	logger.Println("Hardcoded user created successfully")

	err = store.CreateFollowers(user, user2)
	if err != nil {
		logger.Fatal("Error creating Follower:", err)
		return
	}
	logger.Println("Hardcoded follower created successfully")

	// PROVERA ZA PREPORUKE DA LI RADE ?
	recommendations, err := store.GetRecommendations("2")
	if err != nil {
		logger.Fatal("Error getting recommendations:", err)
		return
	}

	// Ispisi rezultat u konzoli
	logger.Println("Recommendations for userID 2:")
	for _, user := range recommendations {
		logger.Printf("User ID: %d, Username: %s\n", user.ID, user.Username)
	}

	followerHandler := handlers.NewFollowerHandler(logger, store)
	router := mux.NewRouter()

	router.Use(followerHandler.MiddlewareContentTypeSet)

	// rutiranje ovde
	router.HandleFunc("/users", followerHandler.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/followers", followerHandler.CreateFollowers).Methods(http.MethodPost)
	router.HandleFunc("/followers/all", followerHandler.GetAllFollowers).Methods(http.MethodGet)

	cors := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins([]string{"*"}))

	server := http.Server{
		Addr:         ":8090",
		Handler:      cors(router),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	logger.Println("Server listening on port", port)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	sig := <-sigCh
	logger.Println("Received terminate, graceful shutdown", sig)

	if server.Shutdown(timeoutContext) != nil {
		logger.Fatal("Cannot gracefully shutdown...")
	}
	logger.Println("Server stopped")
}
