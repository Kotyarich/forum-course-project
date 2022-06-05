package http

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"post-service/forum"
	"post-service/models"
	"strconv"
	"time"
)

type PostHandler struct {
	useCase forum.UseCasePost
}


func NewPostHandler(useCase forum.UseCasePost) *PostHandler {
	return &PostHandler{
		useCase: useCase,
	}
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


// Сообщение для обновления существующего
// swagger:model PostUpdate
type postInput struct {
	// Собственно сообщение
	//
	// example: We should be afraid of the Kraken
	Message string `json:"message"`
}

func (h *PostHandler) ChangePostHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "wrong ID")
	}

	var input postInput
	if err := c.Bind(&input); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	post, err := h.useCase.ChangePost(c.Request().Context(), id, input.Message)
	if err == forum.ErrPostNotFound {
		msg := map[string]string{"message": "Post not found"}
		return c.JSON(http.StatusNotFound, msg)
	} else if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, ModelToPost(post))
}


func (h *PostHandler) GetThreadPosts(c echo.Context) error {
	slug := c.Param("slug")
	tid, err := strconv.Atoi(slug)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}


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

	posts, err := h.useCase.GetThreadPosts(c.Request().Context(), tid, limit, offset, since, desc, sort)
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


func (h *PostHandler) ThreadPostCreateHandler(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	defer func() {_ = c.Request().Body.Close()}()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var postsInput []Post
	err = json.Unmarshal(body, &postsInput)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	slug := c.Param("slug")
	tid, err := strconv.Atoi(slug)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	var posts []*models.Post
	for i := 0; i < len(postsInput); i++ {
		posts = append(posts, toModelPost(postsInput[i]))
	}
	posts, err = h.useCase.CreateThreadPost(c.Request().Context(), tid, posts)

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
