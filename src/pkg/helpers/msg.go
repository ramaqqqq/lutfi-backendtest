package helpers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func MsgOk(status int, message string) map[string]interface{} {
	return map[string]interface{}{"code": status, "message": message}
}

func MsgErr(status int, message string, err string) map[string]interface{} {
	return map[string]interface{}{"code": status, "message": message, "error": err}
}

func FormatError(err string) error {
	if strings.Contains(err, "email") {
		Logger("error", "In Server: Email sudah di gunakan")
		return errors.New("Email sudah di gunakan")
	}
	return errors.New(err)
}

func Response(w http.ResponseWriter, status int, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
