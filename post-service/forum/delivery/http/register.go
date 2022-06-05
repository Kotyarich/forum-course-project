package http

import (
	"github.com/labstack/echo/v4"
	"post-service/common"
	"post-service/forum"
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

func registerHTTPPostEndpoints(router *echo.Echo, uc forum.UseCasePost, as forum.AuthService) {
	handler := NewPostHandler(uc)

	// swagger:operation POST /thread/:slug/create threads postsCreate
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
	router.POST("/thread/:slug/create",
		AuthMiddlware(as,
			common.CORSMiddlware(
				handler.ThreadPostCreateHandler)))

	// swagger:operation GET /thread/:slug/posts threads threadGetPosts
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
	router.GET("/thread/:slug/posts",
		common.CORSMiddlware(handler.GetThreadPosts))


	// swagger:operation PUT /post/:id/details posts postUpdate
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
	router.PUT("/post/:id/details",
		AuthMiddlware(as,
				common.CORSMiddlware(
					handler.ChangePostHandler)))
}

func RegisterHTTPEndpoints(router *echo.Echo, uc forum.UseCase, as forum.AuthService) {
	registerHTTPPostEndpoints(router, uc.PostUseCase, as)
}
