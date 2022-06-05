package http

import (
	"github.com/labstack/echo/v4"
	"user-service/common"
	"user-service/auth"
)

// swagger:model
type Users struct {
	Users []UserOutput `json:"users"`
}

func RegisterHTTPEndpoints(router *echo.Echo, uc auth.UseCase) {
	handler := NewHandler(uc)

	// swagger:operation GET /auth/signout auth signOut
	// ---
	// description: "Разлогирование пользователя"
	// summary: "Разлогирование пользователя"
	// responses:
	//   "200":
	//     description: "Cookie Auth проставляется просроченной"
	//   "400":
	//     description: "Пользователь не был авторизован"
	router.GET("/user/signout",
		common.CORSMiddlware(
			AuthMiddleware(handler.SignOutHandler, uc)))

	// swagger:operation GET /auth/check auth userCheck
	// ---
	// description: "Проверка авторизации пользователя, возвращает информацию о пользователе, если он авторизован"
	// summary: "Проверка авторизации пользователя"
	// responses:
	//   "200":
	//     description: "Пользователь авторизован"
	//     schema:
	//       "$ref": '#/definitions/User'
	//   "403":
	//     description: "Пользователь не был авторизован"
	router.GET("/user/check",
		common.CORSMiddlware(
			AuthMiddleware(handler.UserCheckAuthHandler, uc)))

	// swagger:operation POST /auth/:nickname/create auth userCreate
	// ---
	// description: "Создание нового пользователя в базе данных."
	// summary: "Создание нового пользователя"
	// parameters:
	// - type: string
	//   description: "Идентификатор пользователя."
	//   name: nickname
	//   in: path
	//   required: true
	// - description: "Данные пользовательского профиля."
	//   name: profile
	//   in: body
	//   required: true
	//   schema:
	//     $ref: '#/definitions/User'
    // responses:
	//   '201':
	//     description: "Пользователь успешно создан. Возвращает данные созданного пользователя."
	//     schema:
	//       $ref: '#/definitions/User'
    //   '409':
	//     description: "Пользователь уже присутсвует в базе данных. Возвращает данные ранее созданных пользователей с тем же nickname-ом иои email-ом."
	//     schema:
	//       $ref: '#/definitions/Users'
	router.POST("/user/:nickname/create",
		common.CORSMiddlware(handler.UserCreateHandler))

	// swagger:operation POST /auth/auth auth userAuth
	// ---
	// description: "Авторизая пользователя, возвращает cookie c токеном."
	// summary: "Авторизация пользователя"
	// parameters:
	// - description: "Данные пользователя."
	//   name: credentials
	//   in: body
	//   required: true
	//   schema:
	//     $ref: '#/definitions/UserLogIn'
	// responses:
	//   "200":
	//     description: "Пользователь успешно авторизован. Устанавливается HttpOnly cookie Auth с токеном"
	//   "403":
	//     description: "Неверные данные пользователя"
	router.POST("/user/auth", common.CORSMiddlware(handler.UserAuthHandler))
}
