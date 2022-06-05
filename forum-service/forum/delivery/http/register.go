package http

import (
	"forum-service/common"
	"forum-service/forum"
	"github.com/labstack/echo/v4"
)


func registerHTTPForumEndpoints(router *echo.Echo, uc forum.UseCaseForum, as forum.AuthService) {
	handler := NewForumHandler(uc)

	// swagger:operation POST /forum/create forums forumCreate
	// ---
	// summary: "Создание форума."
	// description: "Создание нового форума."
	// parameters:
	// - name: forum
	//   in: body
	//   description: "Данные форума"
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/ForumInput"
	// responses:
	//   "201":
	//     description: "Форум успешно создан. Возвращает данные созданного форума."
	//     schema:
	//       "$ref": "#/definitions/Forum"
	//   "400":
	//     description: "Владелец форума не найден"
	//     schema:
	//       "$ref": "#/responses/Error"
	//   "404":
	//     description: "Форум уже существует. Возвращает данные ранее созданного форума."
	//     schema:
	//       "$ref": "#/definitions/Forum"
	router.POST("/forum/create",
		AuthMiddlware(as,
			common.CORSMiddlware(
				handler.ForumCreateHandler)))

	// swagger:operation GET /forums forums forumGetAll
	// ---
	// description: "Получение информации о всех форумах."
	// summary: "Получить все форумы"
	// responses:
	//   '200':
	//     description: "Информация о форуме."
	//     schema:
	//       $ref: '#/definitions/Forum'
	router.GET("/forums",
		common.CORSMiddlware(handler.ForumsHandler))

	// swagger:operation GET /forum/:slug/details forums forumGetOne
	// ---
	// description: "Получение информации о форуме по его идентификатору."
    // summary: "Получение информации о форуме"
	// parameters:
	// - type: string
	//   format: identity
	//   description: "Идентификатор форума."
	//   name: slug
	//   in: path
	//   required: true
	// responses:
	//   '200':
	//     description: "Информация о форуме."
	//     schema:
	//       $ref: '#/definitions/Forum'
	//   '404':
	//     description: "Форум отсутсвует в системе."
	//     schema:
	//       $ref: '#/responses/Error'
	router.GET("/forum/:slug/details",
		common.CORSMiddlware(handler.ForumDetailsHandler))

	// swagger:operation GET /forum/:slug/users forums forumGetUsers
	// ---
	// description: "Получение списка пользователей, у которых есть пост или ветка обсуждения в данном форуме. Пользователи выводятся отсортированные по nickname в порядке возрастания. Порядок сотрировки должен соответсвовать побайтовому сравнение в нижнем регистре."
	// summary: "Пользователи данного форума"
	// parameters:
	// - type: string
	//   format: identity
	//   description: "Идентификатор форума."
	//   name: slug
	//   in: path
	//   required: true
	// - maximum: 10000
	//   minimum: 1
	//   type: number
	//   format: int32
	//   default: 100
	//   description: "Максимальное кол-во возвращаемых записей."
	//   name: limit
	//   in: query
	// - type: string
	//   format: identity
	//   description: "Идентификатор пользователя, с которого будут выводиться пользоватли (пользователь с данным идентификатором в результат не попадает)."
	//   name: since
	//   in: query
	// - type: boolean
	//   description: "Флаг сортировки по убыванию."
	//   name: desc
	//   in: query
	// responses:
	//   '200':
	//     description: "Информация о пользователях форума."
	//     schema:
	//       $ref: '#/definitions/Users'
    //   '404':
	//     description: "Форум отсутсвует в системе."
	//     schema:
	//       $ref: '#/responses/Error'
	router.GET("/forum/:slug/users",
		common.CORSMiddlware(handler.ForumUsersHandler))
}

func RegisterHTTPEndpoints(router *echo.Echo, uc forum.UseCase, as forum.AuthService) {
	registerHTTPForumEndpoints(router, uc.ForumUseCase, as)
}
