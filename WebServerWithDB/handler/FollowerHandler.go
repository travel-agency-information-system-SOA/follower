package handler

import (
	"database-example/model"
	repository "database-example/repo"
	"encoding/json"
	"log"
	"net/http"
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

func (m *FollowerHandler) CreateFollowers(rw http.ResponseWriter, h *http.Request) {
	user := h.Context().Value(KeyProduct{}).(*model.User)
	followingUser := h.Context().Value(KeyProduct{}).(*model.User)

	err := m.repo.CreateFollowers(user, followingUser)
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

func (m *FollowerHandler) GetRecommendationsHandler(rw http.ResponseWriter, h *http.Request) {
	userID := h.Context().Value(KeyProduct{}).(string)
	recommendations, err := m.repo.GetRecommendations(userID)
	if err != nil {
		m.logger.Print("Database exception: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	recommendationsJSON, err := json.Marshal(recommendations)
	if err != nil {
		m.logger.Print("Error marshaling recommendations to JSON: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(recommendationsJSON)
}

func (m *FollowerHandler) MiddlewareContentTypeSet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		m.logger.Println("Method [", h.Method, "] - Hit path :", h.URL.Path)

		rw.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(rw, h)
	})
}
