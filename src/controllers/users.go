package controllers

import (
	"api/src/db"
	"api/src/models"
	"api/src/repositories"
	"api/src/responses"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

//CreateUser insert user in database
func CreateUser(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.AppError(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User
	if err := json.Unmarshal(reqBody, &user); err != nil {
		responses.AppError(w, http.StatusBadRequest, err)
		return
	}

	if err = user.Prepare("create"); err != nil {
		responses.AppError(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.CreateConnection()
	if err != nil {
		responses.AppError(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewUserRepository(db)
	user.ID, err = repository.Create(user)
	if err != nil {
		responses.AppError(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, user)
}

//FindAllUsers find all user in database
func FindAllUsersFilteredByNameOrNick(w http.ResponseWriter, r *http.Request) {
	nameOrNick := strings.ToLower(r.URL.Query().Get("user"))

	db, err := db.CreateConnection()
	if err != nil {
		responses.AppError(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewUserRepository(db)
	users, err := repository.Find(nameOrNick)
	if err != nil {
		responses.AppError(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, users)
}

//FindUserById find one user in database
func FindUserById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userId, err := strconv.ParseUint(params["userId"], 10, 64)
	if err != nil {
		responses.AppError(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.CreateConnection()
	if err != nil {
		responses.AppError(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewUserRepository(db)
	user, err := repository.FindByID(userId)
	if err != nil {
		responses.AppError(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, user)
}

//UpdateUserById update one user in database
func UpdateUserById(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	userId, err := strconv.ParseUint(param["userId"], 10, 64)
	if err != nil {
		responses.AppError(w, http.StatusBadRequest, err)
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.AppError(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User
	if err = json.Unmarshal(reqBody, &user); err != nil {
		responses.AppError(w, http.StatusBadRequest, err)
		return
	}

	if err = user.Prepare("edit"); err != nil {
		responses.AppError(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.CreateConnection()
	if err != nil {
		responses.AppError(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewUserRepository(db)
	if err = repository.Update(userId, user); err != nil {
		responses.AppError(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

//DeleteUser remove user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId, err := strconv.ParseUint(params["userId"], 10, 64)
	if err != nil {
		responses.AppError(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.CreateConnection()
	if err != nil {
		responses.AppError(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewUserRepository(db)
	if err = repository.Delete(userId); err != nil {
		responses.AppError(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}
