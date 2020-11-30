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

type Thread struct {
	Author     string    `json:"author"`
	Slug       *string   `json:"slug"`
	Votes      int       `json:"votes"`
	Title      string    `json:"title"`
	Created    time.Time `json:"created"`
	ForumName  string    `json:"forum"`
	Id         int       `json:"id"`
	Message    string    `json:"message"`
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

type ThreadUpdate struct {
	Message string `json:"message"`
	Title   string `json:"title"`
}

type Post struct {
	Author    string    `json:"author"`
	Created   time.Time `json:"created"`
	ForumName string    `json:"forum"`
	Pk        int       `json:"id"`
	IsEdited  bool      `json:"isEdited"`
	Message   string    `json:"message"`
	Parent    int       `json:"parent"`
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

type Vote struct {
	Voice    int    `json:"voice"`
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
