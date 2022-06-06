package controllers

import (
	"api/src/authentication"
	"api/src/db"
	"api/src/models"
	"api/src/repositories"
	"api/src/responses"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CreatePublication add new publication in database
func CreatePublication(w http.ResponseWriter, r *http.Request) {
	userId, err := authentication.GetUserID(r)
	if err != nil {
		responses.AppError(w, http.StatusUnauthorized, err)
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.AppError(w, http.StatusUnprocessableEntity, err)
		return
	}

	var publication models.Publication
	if err = json.Unmarshal(reqBody, &publication); err != nil {
		return
	}

	publication.AuthorId = userId

	if err = publication.Prepare(); err != nil {
		responses.AppError(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.CreateConnection()
	if err != nil {
		responses.AppError(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewPublicationRepository(db)
	publication.ID, err = repository.Create(publication)

	if err != nil {
		responses.AppError(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, publication)
}

// FindAllPublicationsByUser find publications from user
func FindAllPublicationsByUser(w http.ResponseWriter, r *http.Request) {
	userId, err := authentication.GetUserID(r)
	if err != nil {
		responses.AppError(w, http.StatusUnauthorized, err)
		return
	}

	db, err := db.CreateConnection()
	if err != nil {
		responses.AppError(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewPublicationRepository(db)
	publications, err := repository.Find(userId)
	if err != nil {
		responses.AppError(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, publications)
}

// FindPublicationByID find publications by publication id
func FindPublicationByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	publicationId, err := strconv.ParseUint(params["publicationId"], 10, 64)
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

	repository := repositories.NewPublicationRepository(db)
	publication, err := repository.FindById(publicationId)
	if err != nil {
		responses.AppError(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, publication)
}

// UpdatePublicationByID update publication by publication id
func UpdatePublicationByID(w http.ResponseWriter, r *http.Request) {
	userId, err := authentication.GetUserID(r)
	if err != nil {
		responses.AppError(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	publicationId, err := strconv.ParseUint(params["publicationId"], 10, 64)
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

	repository := repositories.NewPublicationRepository(db)
	existPublication, err := repository.FindById(publicationId)
	if err != nil {
		responses.AppError(w, http.StatusInternalServerError, err)
		return
	}

	if existPublication.AuthorId != userId {
		responses.AppError(w, http.StatusForbidden, errors.New("Não é possível atualizar uma publicação que não seja sua"))
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.AppError(w, http.StatusUnprocessableEntity, err)
		return
	}

	var publication models.Publication
	if err = json.Unmarshal(reqBody, &publication); err != nil {
		responses.AppError(w, http.StatusBadRequest, err)
		return
	}

	if err = publication.Prepare(); err != nil {
		responses.AppError(w, http.StatusBadRequest, err)
		return
	}

	if err = repository.Update(publicationId, publication); err != nil {
		responses.AppError(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

// DeletePublicationByID remove publication by publication id
func DeletePublicationByID(w http.ResponseWriter, r *http.Request) {
	userId, err := authentication.GetUserID(r)
	if err != nil {
		responses.AppError(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	publicationId, err := strconv.ParseUint(params["publicationId"], 10, 64)
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

	repository := repositories.NewPublicationRepository(db)
	existPublication, err := repository.FindById(publicationId)
	if err != nil {
		responses.AppError(w, http.StatusInternalServerError, err)
		return
	}

	if existPublication.AuthorId != userId {
		responses.AppError(w, http.StatusForbidden, errors.New("Não é possível remover uma publicação que não seja sua"))
		return
	}

	if err = repository.Delete(publicationId); err != nil {
		responses.AppError(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}
