package handler

import (
	"database-example/model"
	repository "database-example/repo"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type KeyProduct struct{}

type FollowerHandler struct {
	logger *log.Logger
	repo   *repository.FollowerRepository
}

func NewFollowerHandler(l *log.Logger, r *repository.FollowerRepository) *FollowerHandler {
	return &FollowerHandler{l, r}
}

func (m *FollowerHandler) CreateUser(rw http.ResponseWriter, h *http.Request) {
	user := h.Context().Value(KeyProduct{}).(*model.User)
	err := m.repo.CreateUser(user)
	if err != nil {
		m.logger.Print("Database exception: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

// func (m *FollowerHandler) CreateFollowers(rw http.ResponseWriter, h *http.Request) {
// 	// Dekodiranje JSON tela zahteva
// 	var follower model.Followers
// 	body, err := ioutil.ReadAll(h.Body)
// 	if err != nil {
// 		m.logger.Println("Failed to read request body:", err)
// 		rw.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
// 	defer h.Body.Close()

// 	if err := json.Unmarshal(body, &follower); err != nil {
// 		m.logger.Println("Failed to decode request body:", err)
// 		rw.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	// Konvertovanje stringa u int za UserID
// 	userID, err := strconv.Atoi(follower.UserId)
// 	if err != nil {
// 		m.logger.Println("Failed to convert UserID to int:", err)
// 		rw.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	// Konvertovanje stringa u int za FollowingUserID
// 	followingUserID, err := strconv.Atoi(follower.FollowingUserId)
// 	if err != nil {
// 		m.logger.Println("Failed to convert FollowingUserID to int:", err)
// 		rw.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	// Dobavljanje korisnika iz baze podataka
// 	user, err := m.repo.GetUserById(userID)
// 	if err != nil {
// 		m.logger.Println("Failed to get user:", err)
// 		rw.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	// Dobavljanje korisnika koji prati
// 	followingUser, err := m.repo.GetUserById(followingUserID)
// 	if err != nil {
// 		m.logger.Println("Failed to get following user:", err)
// 		rw.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	// Kreiranje novog followera
// 	err = m.repo.CreateFollowers(&user, &followingUser)
// 	if err != nil {
// 		m.logger.Println("Database exception:", err)
// 		rw.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	rw.WriteHeader(http.StatusCreated)
// }

func (m *FollowerHandler) CreateFollowers(rw http.ResponseWriter, h *http.Request) {
	// Dobavljanje ID-jeva korisnika i korisnika koji ga prati iz putanje
	params := mux.Vars(h)

	userID := params["userId"]
	if userID == "" {
		m.logger.Println("User ID is missing")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	followerID := params["followerId"]
	if followerID == "" {
		m.logger.Println("Invalid follower ID")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Pretvaranje userID u int
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		m.logger.Println("Invalid user ID format:", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Dobavljanje korisnika i korisnika koji ga prati preko servisa
	user, err := m.repo.GetUserById(userIDInt)
	if err != nil {
		m.logger.Println("Failed to get user:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Ispisivanje informacija o pronađenom korisniku
	m.logger.Printf("User found - ID: %d, Username: %s\n", user.ID, user.Username)

	// Pretvaranje followerID u int
	followerIDInt, err := strconv.Atoi(followerID)
	if err != nil {
		m.logger.Println("Invalid follower ID format:", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	follower1, err := m.repo.GetUserById(followerIDInt)
	if err != nil {
		m.logger.Println("Failed to get follower:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Ispisivanje informacija o pronađenom korisniku
	m.logger.Printf("User found - ID: %d, Username: %s\n", follower1.ID, follower1.Username)

	// Ponovno dodeljivanje vrednosti promenljivoj err
	err = m.repo.CreateFollowers(user, follower1)
	if err != nil {
		m.logger.Print("Database exception: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func (m *FollowerHandler) GetAllFollowers(rw http.ResponseWriter, h *http.Request) {
	followers, err := m.repo.GetAllFollowers()
	if err != nil {
		m.logger.Print("Database exception: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	followersJSON, err := json.Marshal(followers)
	if err != nil {
		m.logger.Print("Error marshaling followers to JSON: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(followersJSON)
}

func (m *FollowerHandler) MiddlewareContentTypeSet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		m.logger.Println("Method [", h.Method, "] - Hit path :", h.URL.Path)

		rw.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(rw, h)
	})
}

func (m *FollowerHandler) GetFollowings(rw http.ResponseWriter, h *http.Request) {
	m.logger.Println("Usao u handler za get followings")
	// Dobavljanje ID-ja korisnika iz URL putanje
	params := mux.Vars(h)
	userIDStr := params["userId"]

	// Poziv funkcije za dobavljanje korisnika koje prati dati korisnik
	followings, err := m.repo.GetFollowings(userIDStr)
	if err != nil {
		m.logger.Print("Database exception: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Konverzija u JSON i slanje odgovora
	followingsJSON, err := json.Marshal(followings)
	if err != nil {
		m.logger.Print("Error marshaling followings to JSON: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(followingsJSON)
}

func (m *FollowerHandler) GetRecommendations(rw http.ResponseWriter, h *http.Request) {
	m.logger.Println("Usao u handler za get followings")
	// Dobavljanje ID-ja korisnika iz URL putanje
	params := mux.Vars(h)
	userIDStr := params["userId"]

	// Poziv funkcije za dobavljanje korisnika koje prati dati korisnik
	followings, err := m.repo.GetRecommendations(userIDStr)
	if err != nil {
		m.logger.Print("Database exception: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Konverzija u JSON i slanje odgovora
	followingsJSON, err := json.Marshal(followings)
	if err != nil {
		m.logger.Print("Error marshaling followings to JSON: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(followingsJSON)
}
