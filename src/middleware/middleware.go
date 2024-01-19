package middleware

import (
	"context"
	"folkatech-customerIdentity/src/pkg/helpers"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func JwtAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenHeader := r.Header.Get("Authorization")
		middleUrl := os.Getenv("MIDDLE_URL")
		noAuth := []string{
			middleUrl + "/login",
			middleUrl + "/register",
		}

		requestPath := r.URL.Path

		for _, value := range noAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		if tokenHeader == "" {
			helpers.Logger("error", "In Server: No token bearer provided")
			msg := helpers.MsgErr(http.StatusUnauthorized, "Authorization Failed", "No token bearer provided")
			helpers.Response(w, http.StatusUnauthorized, msg)
			return
		}

		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			helpers.Logger("error", "In Server: No token bearer provided")
			msg := helpers.MsgErr(http.StatusUnauthorized, "Authorization Failed", "No token bearer provided")
			helpers.Response(w, http.StatusUnauthorized, msg)
			return
		}

		tokenPart := splitted[1]

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenPart, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			helpers.Logger("Failure", "In Server: "+err.Error())
			msg := helpers.MsgErr(http.StatusUnauthorized, "Authorization Failed", "Token in expired")
			helpers.Response(w, http.StatusUnauthorized, msg)
			return
		}

		if !token.Valid {
			helpers.Logger("Failure", "In Server: "+err.Error())
			msg := helpers.MsgErr(http.StatusUnauthorized, "Authorization Failed", "Token in expired")
			helpers.Response(w, http.StatusUnauthorized, msg)
			return
		}

		ctx := context.WithValue(r.Context(), "user", claims)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

	})
}

func CreateToken(uid, username, email string) (map[string]string, error) {

	expTime := time.Now().Add(time.Minute * 120).Unix()

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user_id"] = uid
	claims["username"] = username
	claims["email"] = email
	claims["exp"] = expTime

	access, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		helpers.Logger("error", "In Server: "+err.Error())
		return nil, err
	}

	refToken := jwt.New(jwt.SigningMethodHS256)
	refClaims := refToken.Claims.(jwt.MapClaims)
	refClaims["authorized"] = true
	refClaims["user_id"] = uid
	refClaims["username"] = username
	refClaims["email"] = email
	refClaims["exp"] = expTime

	refresh, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		helpers.Logger("error", "In Server: "+err.Error())
		return nil, err
	}

	return map[string]string{"accessToken": access, "refreshToken": refresh}, nil

}
