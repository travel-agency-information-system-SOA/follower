package handler

import (
	"database-example/model"
	"database-example/service"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type FollowerHandler struct {
	FollowerService *service.FollowerService
}

func (handler *FollowerHandler) Get(writer http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	log.Printf("Follower sa id-em %s", id)
	// student, err := handler.StudentService.FindStudent(id)
	// writer.Header().Set("Content-Type", "application/json")
	// if err != nil {
	// 	writer.WriteHeader(http.StatusNotFound)
	// 	return
	// }
	writer.WriteHeader(http.StatusOK)
	// json.NewEncoder(writer).Encode(student)
}

func (handler *FollowerHandler) Create(writer http.ResponseWriter, req *http.Request) {
	//ResponseWriter - pisanje odgovora
	//Request - dolazni zahtev
	var follower model.Followers
	err := json.NewDecoder(req.Body).Decode(&follower) //dekodiranje json zahteva
	if err != nil {
		println("Error while parsing json")
		println("Greska:", err.Error())
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.FollowerService.Create(&follower)
	if err != nil {
		println("Error while creating a new tour")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(follower) // dodala sam
}
