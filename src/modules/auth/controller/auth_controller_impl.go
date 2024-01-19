package controller

import (
	"encoding/json"
	"folkatech-customerIdentity/src/modules/auth/model"
	"folkatech-customerIdentity/src/modules/auth/service"
	"folkatech-customerIdentity/src/pkg/helpers"
	"net/http"
)

type AuthControllerImpl struct {
	Service service.AuthService
}

func NewAuthController(service service.AuthService) AuthController {
	return &AuthControllerImpl{Service: service}
}

func (c *AuthControllerImpl) Login(w http.ResponseWriter, r *http.Request) {
	datum := model.AuthLogin{}
	err := json.NewDecoder(r.Body).Decode(&datum)
	if err != nil {
		helpers.Logger("error", "In Server: Oopss server someting wrong"+err.Error())
		msg := helpers.MsgErr(http.StatusInternalServerError, "internal server error: ", err.Error())
		helpers.Response(w, http.StatusInternalServerError, msg)
		return
	}

	result, err := c.Service.Login(r.Context(), datum)
	if err != nil {
		helpers.Logger("error", "In Server: "+err.Error())
		msg := helpers.MsgErr(http.StatusBadRequest, "bad request", err.Error())
		helpers.Response(w, http.StatusBadRequest, msg)
		return
	}

	rMsg := helpers.MsgOk(200, "Successfully")
	rMsg["body"] = result
	logger, _ := json.Marshal(rMsg)
	helpers.Logger("info", "login: "+string(logger))
	helpers.Response(w, http.StatusOK, rMsg)
}

func (c *AuthControllerImpl) Register(w http.ResponseWriter, r *http.Request) {
	datum := model.AuthLogin{}
	err := json.NewDecoder(r.Body).Decode(&datum)
	if err != nil {
		helpers.Logger("error", "In Server: Oopss server someting wrong"+err.Error())
		msg := helpers.MsgErr(http.StatusInternalServerError, "internal server error: ", err.Error())
		helpers.Response(w, http.StatusInternalServerError, msg)
		return
	}

	result, err := c.Service.Register(r.Context(), datum)
	if err != nil {
		helpers.Logger("error", "In Server: "+err.Error())
		format := helpers.FormatError(err.Error())
		msg := helpers.MsgErr(http.StatusBadRequest, "bad request", format.Error())
		helpers.Response(w, http.StatusBadRequest, msg)
		return
	}

	rMsg := helpers.MsgOk(201, "Successfully")
	rMsg["body"] = result
	logger, _ := json.Marshal(rMsg)
	helpers.Logger("info", "created user: "+string(logger))
	helpers.Response(w, http.StatusCreated, rMsg)
}
