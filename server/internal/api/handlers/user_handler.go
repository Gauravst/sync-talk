package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gauravst/real-time-chat/internal/models"
	"github.com/gauravst/real-time-chat/internal/services"
	"github.com/gauravst/real-time-chat/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

// type contextKey string
//
// const userDataKey contextKey = "userData"

func GetAllUsers(userService services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func GetUser(userService services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var userData models.User
		userData, ok := r.Context().Value(userDataKey).(models.User)
		if !ok {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("user data not found")))
			return
		}

		response.WriteJson(w, http.StatusOK, userData)
		return
	}
}

func GetUserById(userService services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == " " {
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(fmt.Errorf("name parms not found")))
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		var userData models.User
		userData, ok := r.Context().Value(userDataKey).(models.User)
		if !ok {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("user data not found")))
			return
		}

		if userData.Role != "ADMIN" && userData.Id != idInt {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		data, err := userService.GetUserByID(idInt)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, data)
		return
	}
}

func UpdateUser(userService services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == " " {
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(fmt.Errorf("name parms not found")))
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		var data models.UserRequest
		err = json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(err))
			return
		}

		var userData models.User
		userData, ok := r.Context().Value(userDataKey).(models.User)
		if !ok {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("user data not found")))
			return
		}

		if userData.Role != "ADMIN" && userData.Id != idInt {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		err = validator.New().Struct(data)
		if err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		err = userService.UpdateUser(&data)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, data)
		return
	}
}

func DeleteUser(userService services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == " " {
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(fmt.Errorf("name parms not found")))
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		var userData models.User
		userData, ok := r.Context().Value(userDataKey).(models.User)
		if !ok {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("user data not found")))
			return
		}

		if userData.Role != "ADMIN" && userData.Id != idInt {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		err = userService.DeleteUser(idInt)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, "User deleted")
		return
	}
}
