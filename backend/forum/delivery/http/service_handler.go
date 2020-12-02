package http

import (
	"dbProject/forum"
	"dbProject/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ServiceHandler struct {
	useCase forum.UseCaseService
}


func NewServiceHandler(useCase forum.UseCaseService) *ServiceHandler {
	return &ServiceHandler{
		useCase: useCase,
	}
}


func (h *ServiceHandler) ClearHandler(c echo.Context) error {
	err := h.useCase.Clear(c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "")
}

//swagger:model
type Status struct {
	// Кол-во разделов в базе данных
	//
	// example: 100
	Forums  int `json:"forum"`

	// Кол-во сообщений в базе данных
	//
	// example: 1000000
	Posts   int `json:"post"`

	// Кол-во веток обсуждений в базе данных
	//
	// example: 1000
	Threads int `json:"thread"`

	// Кол-во пользователей в базе данных
	//
	// example: 1000
	Users   int `json:"user"`
}

func modelToStatus(status *models.Status) Status {
	return Status{
		Forums:  status.Forums,
		Posts:   status.Posts,
		Threads: status.Threads,
		Users:   status.Users,
	}
}

func (h *ServiceHandler) StatusHandler(c echo.Context) error {
	stats, err := h.useCase.Status(c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, modelToStatus(stats))
}
