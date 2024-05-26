package main

import (
	"context"
	"database-example/proto/follower"
	"database-example/repo"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8090"
	}

	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logger := log.New(os.Stdout, "[follower-api] ", log.LstdFlags)
	storeLogger := log.New(os.Stdout, "[follower-store] ", log.LstdFlags)

	store, err := repo.New(storeLogger)
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

	// lis, err := net.Listen("tcp", "localhost:8090")
	// if err != nil {
	// 	log.Fatalf("failed to listen: %v", err)
	// }

	// var opts []grpc.ServerOption
	// grpcServer := grpc.NewServer(opts...)

	// follower.RegisterFollowerServer(grpcServer, Server{FollowerRepo: store}) //da li store?
	// reflection.Register(grpcServer)
	// grpcServer.Serve(lis)

	lis, err := net.Listen("tcp", "followers:8090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	follower.RegisterFollowerServer(grpcServer, &Server{FollowerRepo: store})

	reflection.Register(grpcServer)
	log.Println("gRPC server starting on port 8090")
	grpcServer.Serve(lis)

}

type Server struct {
	follower.UnimplementedFollowerServer
	FollowerRepo *repo.FollowerRepository
}

func (s *Server) CreateNewFollowing(ctx context.Context, request *follower.UserFollowingDto) (*follower.NeoFollowerDto, error) {
	println("USAO JE NA FOLLOWERS")
	userID := request.UserId
	followerID := request.FollowerId
	userIDInt := int(userID)
	user, err := s.FollowerRepo.GetUserById(userIDInt)
	if err != nil {
		log.Println("Error getting user:", err)
		return nil, err
	}

	followerIDint := int(followerID)
	follower1, err := s.FollowerRepo.GetUserById(followerIDint)
	if err != nil {
		log.Println("Error getting follower:", err)
		return nil, err
	}

	err = s.FollowerRepo.CreateFollowers(user, follower1)
	if err != nil {
		log.Println("Error creating new following:", err)
		return nil, err
	}

	return &follower.NeoFollowerDto{
		UserId:            strconv.Itoa(user.ID),
		Username:          user.Username,
		FollowingUserId:   strconv.Itoa(follower1.ID),
		FollowingUsername: follower1.Username,
	}, nil
}

func (s *Server) GetUserRecommendations(ctx context.Context, request *follower.Id) (*follower.ListNeoUserDto, error) {
	userId := strconv.Itoa(int(request.Id))
	recommendations, err := s.FollowerRepo.GetRecommendations(userId)
	if err != nil {
		return nil, err
	}

	responseList := make([]*follower.NeoUserDto, len(recommendations))
	for i, recommendation := range recommendations {
		responseList[i] = &follower.NeoUserDto{
			Id:       int32(recommendation.ID),
			Username: recommendation.Username,
			// Ostali podaci, ako ih ima
		}
	}

	return &follower.ListNeoUserDto{
		ResponseList: responseList,
	}, nil
}

func (s *Server) GetFollowingsWithBlogs(ctx context.Context, request *follower.Id) (*follower.ListBlogPostDto, error) {
	// Vraćamo praznu listu BlogPostDto
	return &follower.ListBlogPostDto{
		ResponseList: []*follower.BlogPostDto{},
	}, nil
}

func (s *Server) FindUserFollowings(ctx context.Context, request *follower.Id) (*follower.ListNeoUserDto, error) {
	println("Usao je ovde FOLLOWERSSS")
	userID := strconv.Itoa(int(request.Id))

	followings, err := s.FollowerRepo.GetFollowings(userID)
	if err != nil {
		log.Println("Database exception: ", err)
		return nil, err
	}

	responseList := make([]*follower.NeoUserDto, len(followings))
	for i, following := range followings {
		responseList[i] = &follower.NeoUserDto{
			Id:       int32(following.ID),
			Username: following.Username,
		}
	}

	return &follower.ListNeoUserDto{
		ResponseList: responseList,
	}, nil
}
