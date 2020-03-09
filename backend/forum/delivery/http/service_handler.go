package http

import (
	"dbProject/common"
	"dbProject/forum"
	"dbProject/models"
	"encoding/json"
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


func (h *ServiceHandler) ClearHandler(writer http.ResponseWriter, request *http.Request, ps map[string]string) {
	err := h.useCase.Clear(request.Context())
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

type Status struct {
	Forums  int `json:"forum"`
	Posts   int `json:"post"`
	Threads int `json:"thread"`
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

func (h *ServiceHandler) StatusHandler(writer http.ResponseWriter, request *http.Request, ps map[string]string) {
	stats, err := h.useCase.Status(request.Context())
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(modelToStatus(stats))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	common.WriteData(writer, http.StatusOK, data)
}
