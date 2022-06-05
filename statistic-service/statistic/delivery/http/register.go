package http

import (
	"github.com/labstack/echo/v4"
	"statistic-service/common"
	"statistic-service/statistic"
)

func RegisterHTTPEndpoints(router *echo.Echo, uc statistic.UseCase) {
	handler := NewHandler(uc)

	router.GET("/statistic",
		common.CORSMiddlware(
			AuthMiddleware(handler.GetHandler, uc)))
}
