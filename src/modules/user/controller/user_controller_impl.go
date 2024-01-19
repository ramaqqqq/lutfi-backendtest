package controller

import (
	"encoding/json"
	"folkatech-customerIdentity/src/modules/user/model"
	"folkatech-customerIdentity/src/modules/user/service"
	"folkatech-customerIdentity/src/pkg/helpers"
	"folkatech-customerIdentity/src/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UserControllerImpl struct {
	Service service.UserServ
}

func NewUserController(service service.UserServ) UserController {
	return &UserControllerImpl{Service: service}
}

func (c *UserControllerImpl) CreateUser(w http.ResponseWriter, r *http.Request) {
	data := model.User{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		helpers.Logger("error", "In Server: Oopss server someting wrong"+err.Error())
		msg := helpers.MsgErr(http.StatusInternalServerError, "internal server error: ", err.Error())
		helpers.Response(w, http.StatusInternalServerError, msg)
		return
	}

	err = c.Service.CreateUser(r.Context(), data)
	if err != nil {
		helpers.Logger("error", "In Server: "+err.Error())
		msg := helpers.MsgErr(http.StatusBadRequest, "bad request: ", err.Error())
		helpers.Response(w, http.StatusBadRequest, msg)
		return
	}

	rMsg := helpers.MsgOk(201, "Successfully")
	rMsg["body"] = "Create user success"
	helpers.Logger("info", "Create user success")
	helpers.Response(w, http.StatusCreated, rMsg)
}

func (s *UserControllerImpl) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	idINT, _ := strconv.Atoi(id)
	data := model.User{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		helpers.Logger("error", "In Server: Oopss server someting wrong"+err.Error())
		msg := helpers.MsgErr(http.StatusInternalServerError, "internal server error: ", err.Error())
		helpers.Response(w, http.StatusInternalServerError, msg)
		return
	}

	err = s.Service.UpdateUser(r.Context(), data, int64(idINT))
	if err != nil {
		helpers.Logger("error", "In Server: "+err.Error())
		msg := helpers.MsgErr(http.StatusBadRequest, "bad request", err.Error())
		helpers.Response(w, http.StatusBadRequest, msg)
		return
	}

	rMsg := helpers.MsgOk(200, "Successfully")
	rMsg["body"] = "Updated user success"
	helpers.Logger("info", "Udpated user success")
	helpers.Response(w, http.StatusOK, rMsg)
}

func (s *UserControllerImpl) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	idINT, _ := strconv.Atoi(id)
	err := s.Service.DeleteUser(r.Context(), int64(idINT))
	if err != nil {
		helpers.Logger("error", "In Server: "+err.Error())
		msg := helpers.MsgErr(http.StatusBadRequest, "bad request", err.Error())
		helpers.Response(w, http.StatusBadRequest, msg)
		return
	}

	rMsg := helpers.MsgOk(200, "Successfully")
	rMsg["body"] = "deleted user success"
	helpers.Logger("info", "Deleted user success")
	helpers.Response(w, http.StatusOK, rMsg)
}

func (s *UserControllerImpl) GetList(w http.ResponseWriter, r *http.Request) {
	pg, err := utils.GetPaginateQueryOffset(r)
	if err != nil {
		helpers.Logger("error", "In Server: "+err.Error())
		msg := helpers.MsgErr(http.StatusBadRequest, "bad requesr", err.Error())
		helpers.Response(w, http.StatusBadRequest, msg)
		return
	}

	filter := model.FilterUser{
		AccountNumber:  r.URL.Query().Get("account_number"),
		IdentityNumber: r.URL.Query().Get("identity_number"),
		Search:         r.URL.Query().Get("search"),
	}

	result, err := s.Service.GetList(r.Context(), filter, pg)
	if err != nil {
		helpers.Logger("error", "In Server: "+err.Error())
		msg := helpers.MsgErr(http.StatusBadRequest, "bad requesr", err.Error())
		helpers.Response(w, http.StatusBadRequest, msg)
		return
	}

	rMsg := helpers.MsgOk(200, "Successfully")
	rMsg["body"] = result
	helpers.Logger("info", "view all data user")
	helpers.Response(w, http.StatusOK, rMsg)
}

func (s *UserControllerImpl) GetDetail(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	idINT, _ := strconv.Atoi(id)
	result, err := s.Service.GetDetail(r.Context(), int64(idINT))
	if err != nil {
		helpers.Logger("error", "In Server: "+err.Error())
		msg := helpers.MsgErr(http.StatusBadRequest, "bad requesr", err.Error())
		helpers.Response(w, http.StatusBadRequest, msg)
		return
	}

	rMsg := helpers.MsgOk(200, "Successfully")
	rMsg["body"] = result
	logger, _ := json.Marshal(rMsg)
	helpers.Logger("info", "view detail data: "+string(logger))
	helpers.Response(w, http.StatusOK, rMsg)
}
