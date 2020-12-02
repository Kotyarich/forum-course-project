package http

import (
	"dbProject/common"
	"dbProject/forum"
	"github.com/labstack/echo/v4"
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

func registerHTTPForumEndpoints(router *echo.Echo, uc forum.UseCaseForum) {
	handler := NewForumHandler(uc)

	// swagger:operation POST /api/v1/forum/create forums forumCreate
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
	router.POST("/api/v1/forum/create",
		common.CORSMiddlware(handler.ForumCreateHandler))

	// swagger:operation POST /api/v1/forum/:slug/create forums threadCreate
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
	router.POST("/api/v1/forum/:slug/create",
		common.CORSMiddlware(handler.ThreadCreateHandler))

	// swagger:operation GET /api/v1/forums forums forumGetAll
	// ---
	// description: "Получение информации о всех форумах."
	// summary: "Получить все форумы"
	// responses:
	//   '200':
	//     description: "Информация о форуме."
	//     schema:
	//       $ref: '#/definitions/Forum'
	router.GET("/api/v1/forums",
		common.CORSMiddlware(handler.ForumsHandler))

	// swagger:operation GET /api/v1/forum/:slug/details forums forumGetOne
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
	router.GET("/api/v1/forum/:slug/details",
		common.CORSMiddlware(handler.ForumDetailsHandler))

	// swagger:operation GET /api/v1/forum/:slug/threads forums forumGetThreads
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
	router.GET("/api/v1/forum/:slug/threads",
		common.CORSMiddlware(handler.ForumThreadsHandler))

	// swagger:operation GET /api/v1/forum/:slug/users forums forumGetUsers
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
	router.GET("/api/v1/forum/:slug/users",
		common.CORSMiddlware(handler.ForumUsersHandler))
}

func registerHTTPThreadEndpoints(router *echo.Echo, uc forum.UseCaseThread) {
	handler := NewThreadHandler(uc)

	// swagger:operation POST /api/v1/thread/:slug/create threads postsCreate
	// ---
	// description: "Добавление новых постов в ветку обсуждения на форум. Все посты, созданные в рамках одного вызова данного метода должны иметь одинаковую дату создания (Post.Created)."
	// summary: "Создание новых постов"
	// parameters:
	// - type: string
	//   format: identity
	//   description: "Идентификатор ветки обсуждения."
	//   name: slug_or_id
	//   in: path
	//   required: true
	// - description: Список создаваемых постов.
	//   name: posts
	//   in: body
	//   required: true
	//   schema:
	//     $ref: '#/definitions/Posts'
    // responses:
	//   '201':
	//     description: "Посты успешно созданы. Возвращает данные созданных постов в том же порядке, в котором их передали на вход метода."
	//     schema:
	//       $ref: '#/definitions/Posts'
    //   '404':
	//     description: "Ветка обсуждения отсутствует в базе данных."
	//     schema:
	//       $ref: '#/responses/Error'
    //   '409':
	//     description: "Хотя бы один родительский пост отсутсвует в текущей ветке обсуждения."
	//     schema:
	//       $ref: '#/responses/Error'
	router.POST("/api/v1/thread/:slug/create",
		common.CORSMiddlware(handler.ThreadPostCreateHandler))

	// swagger:operation GET /api/v1/thread/:slug/details threads threadGetOne
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
	router.GET("/api/v1/thread/:slug/details",
		common.CORSMiddlware(handler.GetThreadHandler))

	// swagger:operation POST /api/v1/thread/:slug/details threads threadUpdate
	// ---
	// description: "Обновление ветки обсуждения на форуме."
	// summary: "Обновление ветки"
	// parameters:
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
	router.POST("/api/v1/thread/:slug/details",
		common.CORSMiddlware(handler.PostThreadHandler))

	// swagger:operation GET /api/v1/thread/:slug/posts threads threadGetPosts
	// ---
	// description: "Получение списка сообщений в данной ветке форуме. Сообщения выводятся отсортированные по дате создания."
	// summary: "Сообщения данной ветви обсуждения"
	// parameters:
	// - type: string
	//   format: identity
	//   description: "Идентификатор ветки обсуждения."
	//   name: slug_or_id
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
	// - type: number
	//   format: int64
	//   description: "Идентификатор поста, после которого будут выводиться записи (пост с данным идентификатором в результат не попадает)."
	//   name: since
	//   in: query
	// - enum:
	//   - flat
	//   - tree
	//   - parent_tree
	//   type: string
	//   default: flat
	//   description: "Вид сортировки:\n* flat - по дате, комментарии выводятся простым списком в порядке создания;\n* tree - древовидный, комментарии выводятся отсортированные в дереве по N штук; \n* parent_tree - древовидные с пагинацией по родительским (parent_tree), на странице N родительских комментов и все комментарии прикрепленные к ним, в древвидном отображение."
	//   name: sort
	//   in: query
	// - type: boolean
	//   description: "Флаг сортировки по убыванию."
	//   name: desc
	//   in: query
	// responses:
	//   '200':
	//     description: "Информация о сообщениях форума."
	//     schema:
	//       $ref: '#/definitions/Posts'
    //   '404':
	//     description: "Ветка обсуждения отсутсвует в форуме."
	//     schema:
	//       $ref: '#/responses/Error'
	router.GET("/api/v1/thread/:slug/posts",
		common.CORSMiddlware(handler.GetThreadPosts))

	// swagger:operation POST /api/v1/thread/:slug/vote threads threadVote
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
	router.POST("/api/v1/thread/:slug/vote",
		common.CORSMiddlware(handler.ThreadVoteHandler))
}

func registerHTTPPostEndpoints(router *echo.Echo, uc forum.UseCasePost) {
	handler := NewPostHandler(uc)

	// swagger:operation GET /api/v1/post/:id/details posts postGetOne
	// ---
	// description: "Получение информации о ветке обсуждения по его имени."
	// summary: "Получение информации о ветке обсуждения"
	// parameters:
	// - type: number
	//   format: int64
	//   description: "Идентификатор сообщения."
	//   name: id
	//   in: path
	//   required: true
	// - type: array
	//   items:
	//     enum:
	//       - user
	//       - forum
	//       - thread
	//   type: string
	//   description: "Включение полной информации о соответвующем объекте сообщения. Если тип объекта не указан, то полная информация об этих объектах не передаётся."
	//   name: related
	//   in: query
	// responses:
	//   '200':
	//     description: "Информация о ветке обсуждения."
	//     schema:
	//       $ref: '#/definitions/PostFull'
    //   '404':
	//     description: "Ветка обсуждения отсутсвует в форуме."
	//     schema:
	//       $ref: '#/responses/Error'
	router.GET("/api/v1/post/:id/details",
		common.CORSMiddlware(handler.GetPostHandler))

	// swagger:operation POST /api/v1/post/:id/details posts postUpdate
	// ---
	// description: "Изменение сообщения на форуме. Если сообщение поменяло текст, то оно должно получить отметку `isEdited`."
	// summary: "Изменение сообщения"
	// parameters:
	// - type: number
	//   format: int64
	//   description: "Идентификатор сообщения."
	//   name: id
	//   in: path
	//   required: true
	// - description: Изменения сообщения.
	//   name: post
	//   in: body
	//   required: true
	//   schema:
	//     $ref: '#/definitions/PostUpdate'
    // responses:
	//   '200':
	//     description: "Информация о сообщении."
	//     schema:
	//       $ref: '#/definitions/Post'
    //   '404':
	//     description: "Сообщение отсутсвует в форуме."
	//     schema:
	//       $ref: '#/responses/Error'
	router.POST("/api/v1/post/:id/details",
		common.CORSMiddlware(handler.ChangePostHandler))
}

func registerHTTPServiceEndpoints(router *echo.Echo, uc forum.UseCaseService) {
	handler := NewServiceHandler(uc)

	// swagger:operation GET /api/v1/service/clear service clear
	// ---
	// description: "Безвозвратное удаление всей пользовательской информации из базы данных."
	// consumes:
	// - application/json
	// - application/octet-stream
	// summary: "Очистка всех данных в базе"
	// responses:
	//   '200':
	//     description: "Очистка базы успешно завершена"
	router.POST("/api/service/clear",
		common.CORSMiddlware(handler.ClearHandler))

	// swagger:operation GET /api/v1/service/status service status
	// ---
	// description: "Получение инфомарции о базе данных."
	// summary: "Получение инфомарции о базе данных"
	// responses:
	//   '200':
	//     description: "Кол-во записей в базе данных"
	//     schema:
	//       $ref: '#/definitions/Status'
	router.GET("/api/service/status",
		common.CORSMiddlware(handler.StatusHandler))
}

func RegisterHTTPEndpoints(router *echo.Echo, uc forum.UseCase) {
	registerHTTPForumEndpoints(router, uc.ForumUseCase)
	registerHTTPThreadEndpoints(router, uc.ThreadUseCase)
	registerHTTPPostEndpoints(router, uc.PostUseCase)
	registerHTTPServiceEndpoints(router, uc.ServiceUseCase)
}
