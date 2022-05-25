package controllers

import (
	"api/src/authentication"
	"api/src/db"
	"api/src/models"
	"api/src/repositories"
	"api/src/responses"
	"api/src/secure"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

//Login ensure athenticate user
func Login(w http.ResponseWriter, r *http.Request) {
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

	db, err := db.CreateConnection()
	if err != nil {
		responses.AppError(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewUserRepository(db)
	userExist, err := repository.FindByEmail(user.Email)
	if err != nil {
		responses.AppError(w, http.StatusInternalServerError, err)
	}

	if err = secure.VerifyPassword(userExist.Password, user.Password); err != nil {
		responses.AppError(w, http.StatusUnauthorized, err)
		return
	}

	token, err := authentication.CreateToken(userExist.ID)
	if err != nil {
		responses.AppError(w, http.StatusInternalServerError, err)
	}

	w.Write([]byte(token))

}
