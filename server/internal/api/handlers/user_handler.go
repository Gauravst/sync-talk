package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gauravst/real-time-chat/internal/models"
	"github.com/gauravst/real-time-chat/internal/services"
	"github.com/gauravst/real-time-chat/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func CreateUser(userService services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var user models.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// Request validation
		err = validator.New().Struct(user)
		if err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		// call here services

		err = userService.CreateUser(&user)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		// return response
		response.WriteJson(w, http.StatusCreated, map[string]string{"success": "ok"})
	}
}

func GetUser(userService services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func UpdateUser(userService services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func DeleteUser(userService services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
