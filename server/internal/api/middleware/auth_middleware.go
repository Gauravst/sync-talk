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

const userDataKey contextKey = "userData"

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
						"exp":      time.Now().Add(24 * 30 * time.Hour).Unix(),
					}
					newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
					newAccessTokenString, err := newAccessToken.SignedString([]byte(cfg.JwtPrivateKey))
					if err != nil {
						jwtToken.RemoveAccessToken(w, false)
						response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(err))
						return
					}

					jwtToken.SetAccessToken(w, newAccessTokenString, false)
				}
				//remove accessToken here
				jwtToken.RemoveAccessToken(w, false)
				response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(err))
				return
			}

			ctx := context.WithValue(r.Context(), userDataKey, userData)
			r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
