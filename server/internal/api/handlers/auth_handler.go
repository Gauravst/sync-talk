package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gauravst/real-time-chat/internal/config"
	"github.com/gauravst/real-time-chat/internal/models"
	"github.com/gauravst/real-time-chat/internal/services"
	"github.com/gauravst/real-time-chat/internal/utils/jwtToken"
	"github.com/gauravst/real-time-chat/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func LoginUser(authService services.AuthService, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var user models.LoginRequest

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

		token, err := authService.LoginUser(&user, cfg)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		// seting new access token
		jwtToken.SetAccessToken(w, token, false)

		// return response
		response.WriteJson(w, http.StatusCreated, map[string]string{"success": "ok"})
	}
}
