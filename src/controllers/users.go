package controllers

import (
	"api/src/authentication"
	"api/src/db"
	"api/src/models"
	"api/src/repositories"
	"api/src/responses"
	"api/src/secure"
	"encoding/json"
	"errors"
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

	userIdInToken, err := authentication.GetUserID(r)
	if err != nil {
		responses.AppError(w, http.StatusUnauthorized, err)
		return
	}

	if userId != userIdInToken {
		responses.AppError(w, http.StatusForbidden, errors.New("Não é possível atualizar um usuario que não seja o seu"))
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

	userIdInToken, err := authentication.GetUserID(r)
	if err != nil {
		responses.AppError(w, http.StatusUnauthorized, err)
	}

	if userId != userIdInToken {
		responses.AppError(w, http.StatusForbidden, errors.New("Não é possível remover um usuario que não seja o seu"))
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

// FollowerUser allows one user to follow another
func FollowerUser(w http.ResponseWriter, r *http.Request) {
	followerId, err := authentication.GetUserID(r)
	if err != nil {
		responses.AppError(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	userId, err := strconv.ParseUint(params["userId"], 10, 64)
	if err != nil {
		responses.AppError(w, http.StatusBadRequest, err)
		return
	}

	if followerId == userId {
		responses.AppError(w, http.StatusForbidden, errors.New("Não é possível seguir você mesmo"))
		return
	}

	db, err := db.CreateConnection()
	if err != nil {
		responses.AppError(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewUserRepository(db)
	if err = repository.Follower(userId, followerId); err != nil {
		responses.AppError(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func UnfollowUser(w http.ResponseWriter, r *http.Request) {
	followerId, err := authentication.GetUserID(r)
	if err != nil {
		responses.AppError(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	userId, err := strconv.ParseUint(params["userId"], 10, 64)
	if err != nil {
		responses.AppError(w, http.StatusBadRequest, err)
		return
	}

	if userId == followerId {
		responses.AppError(w, http.StatusForbidden, errors.New("Não é possível parar de seguir você mesmo"))
		return
	}

	db, err := db.CreateConnection()
	if err != nil {
		responses.AppError(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewUserRepository(db)
	if err = repository.Unfollow(userId, followerId); err != nil {
		responses.AppError(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

// FindFollowers find all followers from user
func FindFollowers(w http.ResponseWriter, r *http.Request) {
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
	followers, err := repository.FindFollowersByUserId(userId)
	if err != nil {
		responses.AppError(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, followers)
}

// FindFollowing find all users that user is following
func FindFollowing(w http.ResponseWriter, r *http.Request) {
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
	users, err := repository.FindFollowingByUserId(userId)
	if err != nil {
		responses.AppError(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, users)
}

//UpdatePassword update password from user
func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	userIdInToken, err := authentication.GetUserID(r)
	if err != nil {
		responses.AppError(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	userId, err := strconv.ParseUint(params["userId"], 10, 64)
	if err != nil {
		responses.AppError(w, http.StatusBadRequest, err)
		return
	}

	if userId != userIdInToken {
		responses.AppError(w, http.StatusForbidden, errors.New("Não é possível atualizar a senha de um usuario que não seja o seu"))
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)

	var password models.Password
	if err = json.Unmarshal(reqBody, &password); err != nil {
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
	existPassword, err := repository.FindPasswordById(userId)
	if err != nil {
		responses.AppError(w, http.StatusInternalServerError, err)
		return
	}

	if err = secure.VerifyPassword(existPassword, password.CurrentPassword); err != nil {
		responses.AppError(w, http.StatusBadRequest, errors.New("A senha atual não condiz com a senha existente"))
		return
	}

	passwordWithHash, err := secure.Hash(password.NewPassword)
	if err != nil {
		responses.AppError(w, http.StatusBadRequest, err)
		return
	}

	if err = repository.UpdateUserPassword(userId, string(passwordWithHash)); err != nil {
		responses.AppError(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}
