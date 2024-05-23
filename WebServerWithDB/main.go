package main

import (
	"context"
	"database-example/proto/follower"
	repository "database-example/repo"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
)

func main() {
	/*port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8090"
	}*/

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

	/*followerHandler := handlers.NewFollowerHandler(logger, store)
	router := mux.NewRouter()

	router.Use(followerHandler.MiddlewareContentTypeSet)

	// rutiranje ovde
	router.HandleFunc("/users", followerHandler.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/followers/{userId}/{followerId}", followerHandler.CreateFollowers).Methods(http.MethodPost)
	router.HandleFunc("/followers/all", followerHandler.GetAllFollowers).Methods(http.MethodGet)
	router.HandleFunc("/followers/recommendations/{userId}", followerHandler.GetRecommendations).Methods(http.MethodGet)
	router.HandleFunc("/followers/followings/{userId}", followerHandler.GetFollowings).Methods(http.MethodGet) //izmenila

	cors := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins([]string{"*"}))

	server := http.Server{
		Addr:         ":" + port, //ovde bilo 8080, sada je 5000, //8090
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
	logger.Println("Server stopped")*/

	lis, err := net.Listen("tcp", "localhost:8090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	follower.RegisterFollowerServer(grpcServer, Server{FollowerRepo: store}) //da li store?
	reflection.Register(grpcServer)
	grpcServer.Serve(lis)

}

type Server struct {
	follower.UnimplementedFollowerServer
	FollowerRepo *repository.FollowerRepository
}

func (s *Server) CreateNewFollowing(ctx context.Context, req *follower.UserFollowingDto) (*follower.NeoFollowerDto, error) {
	userIDInt := int(req.UserId)
	user, err := s.FollowerRepo.GetUserById(userIDInt)
	if err != nil {
		println("Failed to get user:", err)
		return nil, grpc.Errorf(codes.Internal, "Failed to get user: %v", err)
	}

	followerIDInt := int(req.FollowerId)
	follower, err := s.FollowerRepo.GetUserById(followerIDInt)
	if err != nil {
		println("Failed to get follower:", err)
		return nil, grpc.Errorf(codes.Internal, "Failed to get follower: %v", err)
	}

	err = s.FollowerRepo.CreateFollowers(user, follower)
	if err != nil {
		println("Database exception:", err)
		return nil, grpc.Errorf(codes.Internal, "Database exception: %v", err)
	}

	return &follower.NeoFollowerDto{
		UserId:            strconv.Itoa(userIDInt),
		Username:          user.Username,
		FollowingUserId:   strconv.Itoa(followerIDInt),
		FollowingUsername: follower.Username,
	}, nil
}
