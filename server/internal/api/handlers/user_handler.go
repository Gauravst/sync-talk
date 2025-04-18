package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gauravst/real-time-chat/internal/api/middleware"
	"github.com/gauravst/real-time-chat/internal/models"
	"github.com/gauravst/real-time-chat/internal/services"
	"github.com/gauravst/real-time-chat/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func GetAllUsers(userService services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userDataRaw := r.Context().Value(middleware.UserDataKey)
		if userDataRaw == nil {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("Unauthorized")))
			return
		}

		// Correct the type assertion to *models.AccessToken
		userData, ok := userDataRaw.(*models.AccessToken)
		if !ok {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("Unauthorized")))
			return
		}

		if userData.Role != "ADMIN" {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("Unauthorized user")))
			return
		}

		data, err := userService.GetAllUsers()
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, data)
		return
	}
}

func GetUser(userService services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get value from context
		userDataRaw := r.Context().Value(middleware.UserDataKey)
		if userDataRaw == nil {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("Unauthorized")))
			return
		}

		// Correct the type assertion to *models.AccessToken
		userData, ok := userDataRaw.(*models.AccessToken)
		if !ok {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("Unauthorized")))
			return
		}

		data := &models.User{
			Id:         userData.UserId,
			Username:   userData.Username,
			Role:       userData.Role,
			ProfilePic: userData.ProfilePic,
		}

		response.WriteJson(w, http.StatusOK, data)
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

		// geting middleware data
		userDataRaw := r.Context().Value(middleware.UserDataKey)
		if userDataRaw == nil {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("Unauthorized")))
			return
		}

		// Correct the type assertion to *models.AccessToken
		userData, ok := userDataRaw.(*models.AccessToken)
		if !ok {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("Unauthorized")))
			return
		}

		if userData.Role != "ADMIN" && userData.UserId != idInt {
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

		// geting middleware data
		userDataRaw := r.Context().Value(middleware.UserDataKey)
		if userDataRaw == nil {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("Unauthorized")))
			return
		}

		// Correct the type assertion to *models.AccessToken
		userData, ok := userDataRaw.(*models.AccessToken)
		if !ok {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("Unauthorized")))
			return
		}

		if userData.Role != "ADMIN" && userData.UserId != idInt {
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

		// geting middleware data
		userDataRaw := r.Context().Value(middleware.UserDataKey)
		if userDataRaw == nil {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("Unauthorized")))
			return
		}

		// Correct the type assertion to *models.AccessToken
		userData, ok := userDataRaw.(*models.AccessToken)
		if !ok {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("Unauthorized")))
			return
		}

		if userData.Role != "ADMIN" && userData.UserId != idInt {
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
