package http

import (
	"github.com/labstack/echo/v4"
	"thread-service/common"
	"thread-service/forum"
)

// Error response payload
// swagger:response Error
type swaggerError struct {
	// in:body
	Body struct {
		// Текстовое описание ошибки.
		Message string `json:"message"`
	}
}

// swagger:model
type Posts struct {
	Posts []Post `json:"posts"`
}

// swagger:model
type Threads struct {
	Threads []Post `json:"threads"`
}

func registerHTTPThreadEndpoints(router *echo.Echo, uc forum.UseCaseThread, as forum.AuthService) {
	handler := NewThreadHandler(uc)

	// swagger:operation GET /forum/:slug/threads forums forumGetThreads
	// ---
	// description: "Получение списка ветвей обсужления данного форума. Ветви обсуждения выводятся отсортированные по дате создания."
	// summary: "Список ветвей обсужления форума"
	// parameters:
	// - type: string
	//   description: "Идентификатор форума."
	//   name: slug
	//   in: path
	//   required: true
	// - maximum: 10000
	//   minimum: 1
	//   type: number
	//   default: 100
	//   description: "Максимальное кол-во возвращаемых записей."
	//   name: limit
	//   in: query
	// - type: string
	//   description: "Дата создания ветви обсуждения, с которой будут выводиться записи (ветвь обсуждения с указанной датой попадает в результат выборки)."
	//   name: since
	//   in: query
	// - type: boolean
	//   description: "Флаг сортировки по убыванию."
	//   name: desc
	//   in: query
	// responses:
	//   '200':
	//     description: "Информация о ветках обсуждения на форуме."
	//     schema:
	//       $ref: '#/definitions/Threads'
	//   '404':
	//     description: "Форум отсутсвует в системе."
	//     schema:
	//       $ref: '#/responses/Error'
	router.GET("/forum/:slug/threads",
		common.CORSMiddlware(handler.ForumThreadsHandler))

	// swagger:operation GET /thread/:slug/details threads threadGetOne
	// ---
	// description: "Получение информации о ветке обсуждения по его имени."
	// summary: "Получение информации о ветке обсуждения"
	// parameters:
	// - type: string
	//   description: "Идентификатор ветки обсуждения."
	//   name: slug_or_id
	//   in: path
	//   required: true
	// responses:
	//   '200':
	//     description: "Информация о ветке обсуждения."
	//     schema:
	//       $ref: '#/definitions/Thread'
	//   '404':
	//     description: "Ветка обсуждения отсутсвует в форуме."
	//     schema:
	//       $ref: '#/responses/Error'
	router.GET("/thread/:slug/details",
		common.CORSMiddlware(handler.GetThreadHandler))

	// swagger:operation PATCH /thread/:slug/details threads threadUpdate
	// ---
	// description: "Обновление ветки обсуждения на форуме."
	// summary: "Обновление ветки"
	// parameters:forum
	// - type: string
	//   format: identity
	//   description: "Идентификатор ветки обсуждения."
	//   name: slug_or_id
	//   in: path
	//   required: true
	// - description: "Данные ветки обсуждения."
	//   name: thread
	//   in: body
	//   required: true
	//   schema:
	//     $ref: '#/definitions/ThreadUpdate'
	// responses:
	//   '200':
	//     description: "Информация о ветке обсуждения."
	//     schema:
	//       $ref: '#/definitions/Thread'
	//   '404':
	//     description: "Ветка обсуждения отсутсвует в форуме."
	//     schema:
	//       $ref: '#/responses/Error'
	router.PATCH("/thread/:slug/details",
		AuthMiddlware(as,
			common.CORSMiddlware(
				handler.PostThreadHandler)))

	// swagger:operation POST /thread/:slug/vote threads threadVote
	// ---
	// description: "Изменение голоса за ветвь обсуждения. Один пользователь учитывается только один раз и может изменить своё мнение."
	// summary: "Проголосовать за ветвь обсуждения"
	// parameters:
	// - type: string
	//   format: identity
	//   description: "Идентификатор ветки обсуждения."
	//   name: slug_or_id
	//   in: path
	//   required: true
	// - description: "Информация о голосовании пользователя."
	//   name: vote
	//   in: body
	//   required: true
	//   schema:
	//     $ref: '#/definitions/Vote'
	// responses:
	//   '200':
	//     description: "Информация о ветке обсуждения."
	//     schema:
	//       $ref: '#/definitions/Thread'
	//   '404':
	//     description: "Ветка обсуждения отсутсвует в форуме."
	//     schema:
	//       $ref: '#/responses/Error'
	router.POST("/thread/:slug/vote",
		AuthMiddlware(as,
			common.CORSMiddlware(
				handler.ThreadVoteHandler)))

	// swagger:operation POST /forum/:slug/create forums threadCreate
	// ---
	// description:	"Добавление новой ветки обсуждения на форум."
	// summary: "Создание ветки"
	// operationId: threadCreate
	// parameters:
	// - type: string
	//   format: identity
	//   description: "Идентификатор форума."
	//   name: slug
	//   in: path
	//   required: true
	// - description: Данные ветки обсуждения.
	//   name: thread
	//   in: body
	//   required: true
	//   schema:
	//     $ref: '#/definitions/Thread'
	// responses:
	//   '201':
	//     description: "Ветка обсуждения успешно создана. Возвращает данные созданной ветки обсуждения."
	//     schema:
	//       $ref: '#/definitions/Thread'
	//   '404':
	//     description: "Автор ветки или форум не найдены."
	//     schema:
	//       $ref: '#/responses/Error'
	//   '409':
	//     description: "Ветка обсуждения уже существует. Возвращает данные ранее созданной ветки обсуждения."
	//     schema:
	//       $ref: '#/definitions/Thread'
	router.POST("/forum/:slug/create",
		AuthMiddlware(as,
			common.CORSMiddlware(
				handler.ThreadCreateHandler)))

}

func RegisterHTTPEndpoints(router *echo.Echo, uc forum.UseCase, as forum.AuthService) {
	registerHTTPThreadEndpoints(router, uc.ThreadUseCase, as)
}
