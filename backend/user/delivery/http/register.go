package http

import (
	"dbProject/common"
	"dbProject/user"
	"github.com/labstack/echo/v4"
)

// swagger:model
type Users struct {
	Users []UserOutput `json:"users"`
}

func RegisterHTTPEndpoints(router *echo.Echo, uc user.UseCase) {
	handler := NewHandler(uc)

	// swagger:operation GET /api/v1/user/signout user signOut
	// ---
	// description: "Разлогирование пользователя"
	// summary: "Разлогирование пользователя"
	// responses:
	//   "200":
	//     description: "Cookie Auth проставляется просроченной"
	//   "400":
	//     description: "Пользователь не был авторизован"
	router.GET("/api/v1/user/signout",
		common.CORSMiddlware(
			AuthMiddleware(handler.SignOutHandler, uc)))

	// swagger:operation GET /api/v1/user/check user userCheck
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
	router.GET("/api/v1/user/check",
		common.CORSMiddlware(
			AuthMiddleware(handler.UserCheckAuthHandler, uc)))

	// swagger:operation POST /api/v1/user/:nickname/create user userCreate
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
	router.POST("/api/v1/user/:nickname/create",
		common.CORSMiddlware(handler.UserCreateHandler))

	// swagger:operation GET /api/v1/user/:nickname/profile user userGetOne
	// ---
	// description: "Получение информации о пользователе форума по его имени."
	// summary: Получение информации о пользователе
	// parameters:
	// - type: string
	//   description: "Идентификатор пользователя."
	//   name: nickname
	//   in: path
	//   required: true
	// responses:
	//   '200':
	//     description: "Информация о пользователе."
	//     schema:
	//       $ref: '#/definitions/User'
    //   '404':
	//     description: "Пользователь отсутсвует в системе."
	//     schema:
	//       $ref: '#/responses/Error'
	router.GET("/api/v1/user/:nickname/profile",
		common.CORSMiddlware(handler.UserGetHandler))

	// swagger:operation PATCH /api/v1/user/:nickname/profile user userUpdate
	// ---
	// description: "Изменение информации в профиле пользователя."
	// summary: "Изменение данных о пользователе"
	// parameters:
	// - type: string
	//   description: "Идентификатор пользователя."
	//   name: nickname
	//   in: path
	//   required: true
	// - description: "Изменения профиля пользователя."
	//   name: profile
	//   in: body
	//   required: true
	//   schema:
	//     $ref: '#/definitions/User'
    // responses:
	//   '200':
	//     description: "Актуальная информация о пользователе после изменения профиля."
	//     schema:
	//       $ref: '#/definitions/User'
    //   '404':
	//     description: "Пользователь отсутсвует в системе."
	//     schema:
	//       $ref: '#/responses/Error'
    //   '409':
	//     description: "Новые данные профиля пользователя конфликтуют с имеющимися пользователями."
	//     schema:
	//       $ref: '#/responses/Error'
	router.PATCH("/api/v1/user/:nickname/profile",
		common.CORSMiddlware(handler.UserPostHandler))

	// swagger:operation POST /api/v1/user/auth user userAuth
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
	router.POST("/api/v1/user/auth", common.CORSMiddlware(handler.UserAuthHandler))
}
