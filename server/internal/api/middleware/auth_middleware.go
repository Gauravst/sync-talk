package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gauravst/real-time-chat/internal/config"
	"github.com/gauravst/real-time-chat/internal/models"
	"github.com/gauravst/real-time-chat/internal/services"
	"github.com/gauravst/real-time-chat/internal/utils/jwtToken"
	"github.com/gauravst/real-time-chat/internal/utils/response"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserDataKey contextKey = "userData"

func Auth(cfg *config.Config, authService services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract the token from the request headers
			cookie, err := r.Cookie("accessToken")
			if err != nil {
				response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(err))
				return
			}
			token := cookie.Value

			// If the token is valid, call the next handler
			// var userData models.User
			userData, err := jwtToken.VerifyJwtAndGetData[models.AccessToken](token, cfg.JwtPrivateKey)
			if err != nil {
				if err.Error() == "token has expired" {
					// log.Print("Token has expired. Please reauthenticate.")
					// accessToken expired, give client next token here
					refreshToken, err := authService.GetRefreshToken(userData.UserId)
					if err != nil {
						//remove accessToken here
						jwtToken.RemoveAccessToken(w, false)
						response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(err))
						return
					}

					//1) get refresh token check it is valid or not
					_, err = jwtToken.VerifyJwtAndGetData[models.AccessToken](refreshToken, cfg.JwtPrivateKey)
					if err != nil {
						// remove accessToken here
						jwtToken.RemoveAccessToken(w, false)
						response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(err))
						return
					}

					//2) if refresh token valid genrete accessToken
					claims := jwt.MapClaims{
						"userId":   userData.UserId,
						"username": userData.Username,
						"role":     userData.Role,
						"exp":      time.Now().Add(2 * time.Minute).Unix(),
					}
					newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
					newAccessTokenString, err := newAccessToken.SignedString([]byte(cfg.JwtPrivateKey))
					if err != nil {
						jwtToken.RemoveAccessToken(w, false)
						response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(err))
						return
					}

					jwtToken.SetAccessToken(w, newAccessTokenString, false)
				} else {
					//remove accessToken here
					jwtToken.RemoveAccessToken(w, false)
					response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(err))
					return
				}
			}

			ctx := context.WithValue(r.Context(), UserDataKey, userData)
			// sending context data with request
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
