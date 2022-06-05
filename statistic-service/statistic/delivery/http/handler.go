package http

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"statistic-service/models"
	"statistic-service/statistic"
)

type Handler struct {
	useCase statistic.UseCase
}

func NewHandler(useCase statistic.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

type statisticOutput struct {
	Users   int64 `json:"users"`
	Posts   int64 `json:"posts"`
	Votes   int64 `json:"votes"`
	Threads int64 `json:"threads"`
	Forums  int64 `json:"forums"`
}

func statisticModelToOutput(statistic models.Status) statisticOutput {
	return statisticOutput{
		Users:   statistic.Users,
		Posts:   statistic.Posts,
		Votes:   statistic.Votes,
		Threads: statistic.Threads,
		Forums:  statistic.Forums,
	}
}

func (h *Handler) GetHandler(c echo.Context) error {
	if c.Request().Context().Value("user") == nil {
		return c.String(http.StatusForbidden, "only for admins")
	}

	status, err := h.useCase.GetStatistic(c.Request().Context())
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "gettting statistic error: %s", err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, statisticModelToOutput(*status))
}
