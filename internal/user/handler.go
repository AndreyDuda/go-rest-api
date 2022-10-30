package user

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"rest_api/internal/apperror"
	"rest_api/internal/handlers"
	"rest_api/pkg/logging"
)

const (
	usersURL = "/users"
	userURL  = "/users/:uuid"
)

type handler struct {
	logger *logging.Logger
}

func NewHandler(logger *logging.Logger) handlers.Handler {
	return &handler{
		logger: logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, usersURL, apperror.Middleware(h.GetList))
	router.HandlerFunc(http.MethodGet, userURL, apperror.Middleware(h.GetUserByUUID))
	router.HandlerFunc(http.MethodPost, usersURL, apperror.Middleware(h.CreateUser))
	router.HandlerFunc(http.MethodPut, userURL, apperror.Middleware(h.UpdateUser))
	router.HandlerFunc(http.MethodPatch, userURL, apperror.Middleware(h.PartUpdateUser))
	router.HandlerFunc(http.MethodDelete, userURL, apperror.Middleware(h.DeleteUser))
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("this is list of users"))

	return nil
}

func (h *handler) GetUserByUUID(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("this is get user by uuid"))

	return nil
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("this is create user"))

	return nil
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusNoContent)
	_, _ = w.Write([]byte("this is update user"))

	return nil
}

func (h *handler) PartUpdateUser(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusNoContent)
	_, _ = w.Write([]byte("this is part update user"))

	return nil
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusNoContent)
	_, _ = w.Write([]byte("this is delete user"))

	return nil
}
