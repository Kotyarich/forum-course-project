package http

import (
	"dbProject/forum"
	"dbProject/models"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type ThreadHandler struct {
	useCase forum.UseCaseThread
}

func NewThreadHandler(useCase forum.UseCaseThread) *ThreadHandler {
	return &ThreadHandler{
		useCase: useCase,
	}
}

// Ветка обсуждений на форуме
// swagger:model
type Thread struct {
	// Пользователь, создавший ветку
	//
	// format: identity
	// example: j.sparrow
	Author     string    `json:"author"`

	// Человекопонятный URL
	//
	// format: identity
	// read only: true
	// example: jones-cache
	Slug       *string   `json:"slug"`

	// Кол-во голосов
	//
	// read only: true
	Votes      int       `json:"votes"`

	// Заголовок ветки
	//
	// example: Davy Jones cache
	Title      string    `json:"title"`

	// Дата создания ветки
	//
	// example: 2017-01-01T00:00:00.000Z
	Created    time.Time `json:"created"`

	// Форум ветки
	//
	// format: identity
	// read only: true
	// example: pirate-stories
	ForumName  string    `json:"forum"`

	// Идентификатор ветки обсуждения
	//
	// read only: true
	// example: 42
	Id         int       `json:"id"`

	// Описание ветки обсуждения
	//
	// example: An urgent need to reveal the hiding place of Davy Jones. Who is willing to help in this matter?
	Message    string    `json:"message"`

	// Количество сообщений в ветке
	//
	// read only: true
	// example: 100
	PostsCount int       `json:"posts"`
}

func threadToModel(t *Thread) *models.Thread {
	return &models.Thread{
		Author:     t.Author,
		Slug:       t.Slug,
		Title:      t.Title,
		Message:    t.Message,
		ForumName:  t.ForumName,
		Id:         t.Id,
		Created:    t.Created,
		Votes:      t.Votes,
		PostsCount: t.PostsCount,
	}
}

func modelToThread(t *models.Thread) *Thread {
	return &Thread{
		Author:     t.Author,
		Slug:       t.Slug,
		Title:      t.Title,
		Message:    t.Message,
		ForumName:  t.ForumName,
		Id:         t.Id,
		Created:    t.Created,
		Votes:      t.Votes,
		PostsCount: t.PostsCount,
	}
}

type ThreadResult struct {
	Author    string    `json:"author"`
	Title     string    `json:"title"`
	Created   time.Time `json:"created"`
	ForumName string    `json:"forum"`
	Id        int       `json:"id"`
	Message   string    `json:"message"`
}

// Сообщение для обновления ветки обсуждения на форуме. Пустые параметры остаются без изменений.
// swagger:model ThreadUpdate
type ThreadUpdate struct {
	// Описание ветки обсуждения
	//
	// example: An urgent need to reveal the hiding place of Davy Jones. Who is willing to help in this matter?
	Message string `json:"message"`

	// Заголовок ветки обсуждения
	//
	// example: Davy Jones cache
	Title   string `json:"title"`
}

// Информация о сообщении
// swagger:model
type Post struct {
	// Автор сообщения
	//
	// format: identity
	// example: j.sparrow
	Author    string    `json:"author"`

	// Дата создания
	//
	// read only: true
	Created   time.Time `json:"created"`

	// Идентификатор форума
	// read only: true
	// format: identity
	ForumName string    `json:"forum"`

	// Идентификатор сообщения
	//
	// read only: true
	Pk        int       `json:"id"`

	// Истина, если сообщение было изменено
	// read only: true
	IsEdited  bool      `json:"isEdited"`

	// Текст сообщения
	//
	// example: We should be afraid of the Kraken
	Message   string    `json:"message"`

	// Идентификатор родительского сообщения
	Parent    int       `json:"parent"`

	// Идентификатор ветви обсуждения данного сообщения
	//
	// read only: true
	Tid       int       `json:"thread"`
}

func toModelPost(p Post) *models.Post {
	return &models.Post{
		Author:    p.Author,
		Created:   p.Created,
		ForumName: p.ForumName,
		Id:        p.Pk,
		IsEdited:  p.IsEdited,
		Message:   p.Message,
		Parent:    p.Parent,
		Tid:       p.Tid,
	}
}

func ModelToPost(p *models.Post) *Post {
	return &Post{
		Author:    p.Author,
		Created:   p.Created,
		ForumName: p.ForumName,
		Pk:        p.Id,
		IsEdited:  p.IsEdited,
		Message:   p.Message,
		Parent:    p.Parent,
		Tid:       p.Tid,
	}
}

func modelsToPostsArray(p []*models.Post) []Post {
	if len(p) == 0 {
		return []Post{}
	}
	var posts []Post
	for i := 0; i < len(p); i++ {
		posts = append(posts, *ModelToPost(p[i]))
	}

	return posts
}

func (h *ThreadHandler) ThreadPostCreateHandler(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	defer c.Request().Body.Close()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var postsInput []Post
	err = json.Unmarshal(body, &postsInput)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	slug := c.Param("slug")
	var posts []*models.Post
	for i := 0; i < len(postsInput); i++ {
		posts = append(posts, toModelPost(postsInput[i]))
	}
	posts, err = h.useCase.CreateThreadPost(c.Request().Context(), slug, posts)

	if err == forum.ErrThreadNotFound {
		msg := map[string]string{"message": "Thread not found"}
		return c.JSON(http.StatusNotFound, msg)
	} else if err == forum.ErrUserNotFound {
		msg := map[string]string{"message": "User not found"}
		return c.JSON(http.StatusNotFound, msg)
	} else if err == forum.ErrWrongParentsThread {
		msg := map[string]string{"message": "Parent in another thread"}
		return c.JSON(http.StatusNotFound, msg)
	} else if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, modelsToPostsArray(posts))
}

func (h *ThreadHandler) GetThreadHandler(c echo.Context) error {
	slug := c.Param("slug")

	thread, err := h.useCase.GetThread(c.Request().Context(), slug)
	if err == forum.ErrThreadNotFound {
		msg := map[string]string{"message": "Thread not found"}
		return c.JSON(http.StatusNotFound, msg)
	} else if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, modelToThread(thread))
}

func (h *ThreadHandler) PostThreadHandler(c echo.Context) error {
	slug := c.Param("slug")

	var input ThreadUpdate
	if err := c.Bind(&input); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	thread, err := h.useCase.ChangeThread(c.Request().Context(), slug, input.Title, input.Message)
	if err == forum.ErrThreadNotFound {
		msg := map[string]string{"message": "Thread not found"}
		return c.JSON(http.StatusNotFound, msg)
	} else if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, modelToThread(thread))
}

func (h *ThreadHandler) GetThreadPosts(c echo.Context) error {
	slug := c.Param("slug")

	since, err := strconv.Atoi(c.QueryParam("since"))
	if err != nil {
		since = -1
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = -1
	}
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = 0
	}
	desc := c.QueryParam("desc") == "true"
	var sort models.PostSortType
	switch c.QueryParam("sort") {
	case "flat", "":
		sort = models.Flat
	case "tree":
		sort = models.Tree
	case "parent_tree":
		sort = models.ParentTree
	}

	posts, err := h.useCase.GetThreadPosts(c.Request().Context(), slug, limit, offset, since, desc, sort)
	if err == forum.ErrThreadNotFound {
		msg := map[string]string{"message": "Thread not found"}
		return c.JSON(http.StatusNotFound, msg)
	} else if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var result []*Post
	for i := 0; i < len(posts); i++ {
		post := ModelToPost(posts[i])
		result = append(result, post)
	}

	return c.JSON(http.StatusOK, result)
}


// Информация о голосовании пользователя
// swagger:model
type Vote struct {
	// Отданный голос
	//
	// enum:
	//   - -1
	//   - 1
	Voice    int    `json:"voice"`

	// Имя пользователя
	//
	// format: identity
	Nickname string `json:"nickname"`
}

func toModelVote(v *Vote) models.Vote {
	var voice models.Voice
	if v.Voice == 1 {
		voice = models.Up
	} else {
		voice = models.Down
	}

	return models.Vote{
		Voice:    voice,
		Nickname: v.Nickname,
	}
}

func (h *ThreadHandler) ThreadVoteHandler(c echo.Context) error {
	slug := c.Param("slug")

	var vote Vote
	if err := c.Bind(&vote); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	thread, err := h.useCase.VoteForThread(c.Request().Context(), slug, toModelVote(&vote))
	if err == forum.ErrThreadNotFound {
		msg := map[string]string{"message": "Thread not found"}
		return c.JSON(http.StatusNotFound, msg)
	} else if err == forum.ErrUserNotFound {
		msg := map[string]string{"message": "User not found"}
		return c.JSON(http.StatusNotFound, msg)
	} else if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, modelToThread(thread))
}
