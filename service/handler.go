package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"scootin/db"
	"scootin/logger"
	"scootin/models"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Hello, welcome to the Scootin")
}

func BookScooter(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	scooterID := ps.ByName("scooter_id")
	userID := r.Header.Get("user-id")
	if err := db.BookScooter(r.Context(), scooterID, userID); err != nil {
		logger.Errorf("couldn't book scooter %s for user %s: %s", scooterID, userID, err)
		http.Error(w, err.Error(), 500)
	}
}

func ListAvailableScooter(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sc, err := db.ListAvailableScooter(r.Context())
	if err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), 500)
	}

	if err = json.NewEncoder(w).Encode(sc); err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), 500)
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		body []byte
		err  error
		user *models.User
	)
	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), 400)
	}
	if err = json.Unmarshal(body, &user); err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), 400)
	}

	// assign user uuid
	user.ID = uuid.New().String()

	if err := db.CreateUser(r.Context(), user); err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), 500)
	}
}

func CreateScooter(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		body []byte
		err  error
		sc   *models.ScooterInfo
	)
	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), 400)
	}
	if err = json.Unmarshal(body, &sc); err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), 400)
	}

	if err := db.CreateScooter(r.Context(), uuid.New().String()); err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), 500)
	}
}
