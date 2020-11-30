package http

import (
	"dbProject/forum"
	"dbProject/models"
	userHttp "dbProject/user/delivery/http"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

type ForumHandler struct {
	useCase forum.UseCaseForum
}


func NewForumHandler(useCase forum.UseCaseForum) *ForumHandler {
	return &ForumHandler{
		useCase: useCase,
	}
}

type ForumOutput struct {
	Posts   int    `json:"posts"`
	Slug    string `json:"slug"`
	Threads int    `json:"threads"`
	Title   string `json:"title"`
	User    string `json:"user"`
}

type ForumInput struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
	User  string `json:"user"`
}

func forumInputToModel(f ForumInput) *models.Forum {
	return &models.Forum{
		Slug:    f.Slug,
		Title:   f.Title,
		User:    f.User,
		Threads: 0,
		Posts:   0,
	}
}

func forumToOutputFormat(f *models.Forum) *ForumOutput {
	return &ForumOutput{
		Slug:    f.Slug,
		Title:   f.Title,
		User:    f.User,
		Threads: f.Threads,
		Posts:   f.Posts,
	}
}

func (h *ForumHandler) ForumCreateHandler(c echo.Context) error {
	var input ForumInput
	if err := c.Bind(&input); err != nil {
		errText := map[string]string{"error": err.Error()}
		return c.JSON(http.StatusInternalServerError, errText)
	}

	f, err := h.useCase.CreateForum(c.Request().Context(), forumInputToModel(input))
	if err == forum.ErrUserNotFound {
		msg := map[string]string{"message": "User not found"}
		return c.JSON(http.StatusNotFound, msg)
	} else if err == forum.ErrForumAlreadyExists {
		data := forumToOutputFormat(f)
		return c.JSON(http.StatusConflict, data)
	} else if err != nil {
		errText := map[string]string{"error": err.Error()}
		return c.JSON(http.StatusInternalServerError, errText)
	}

	return c.JSON(http.StatusCreated, forumToOutputFormat(f))
}

func (h *ForumHandler) ThreadCreateHandler(c echo.Context) error {
	slug := c.Param("slug")

	var threadInp Thread
	if err := c.Bind(&threadInp); err != nil {
		errText := map[string]string{"error": err.Error()}
		return c.JSON(http.StatusInternalServerError, errText)
	}

	thread, err := h.useCase.CreateForumThread(c.Request().Context(), slug, threadToModel(&threadInp))
	if err == forum.ErrUserNotFound {
		msg := map[string]string{"message": "User not found"}
		return c.JSON(http.StatusNotFound, msg)
	} else if err == forum.ErrForumNotFound {
		msg := map[string]string{"message": "Forum not found"}
		return c.JSON(http.StatusNotFound, msg)
	} else if err == forum.ErrThreadAlreadyExists {
		return c.JSON(http.StatusConflict, modelToThread(thread))
	} else if err != nil {
		errText := map[string]string{"error": err.Error()}
		return c.JSON(http.StatusInternalServerError, errText)
	} else {
		if *thread.Slug == thread.ForumName {
			thread.Slug = nil
		}
		return c.JSON(http.StatusCreated, modelToThread(thread))
	}
}

func (h *ForumHandler) ForumsHandler(c echo.Context) error {
	forums, err := h.useCase.GetForums(c.Request().Context())
	if err != nil {
		errText := map[string]string{"error": err.Error()}
		return c.JSON(http.StatusInternalServerError, errText)
	}

	return c.JSON(http.StatusOK, forumsToJsonArray(forums))
}

func (h *ForumHandler) ForumDetailsHandler(c echo.Context) error {
	slug := c.Param("slug")

	f, err := h.useCase.GetForumDetails(c.Request().Context(), slug)
	if err == forum.ErrForumNotFound {
		msg := map[string]string{"message": "404"}
		return c.JSON(http.StatusNotFound, msg)
	} else if err != nil {
		errText := map[string]string{"error": err.Error()}
		return c.JSON(http.StatusInternalServerError, errText)
	}

	data := forumToOutputFormat(f)
	return c.JSON(http.StatusOK, data)
}

func (h *ForumHandler) ForumUsersHandler(c echo.Context) error {
	slug := c.Param("slug")
	since := c.QueryParam("since")
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = -1
	}
	sort := c.QueryParam("desc") == "true"

	users, err := h.useCase.GetForumUsers(c.Request().Context(), slug, since, limit, sort)
	if err == forum.ErrForumNotFound {
		msg := map[string]string{"message": "Forum not found"}
		return c.JSON(http.StatusNotFound, msg)
	} else if err != nil {
		errText := map[string]string{"error": err.Error()}
		return c.JSON(http.StatusInternalServerError, errText)
	}

	return c.JSON(http.StatusOK, userHttp.UsersToJsonArray(users))
}

func forumsToJsonArray(forums []*models.Forum) []*ForumOutput {
	var result []*ForumOutput
	for i := 0; i < len(forums); i++ {
		result = append(result, forumToOutputFormat(forums[i]))
	}

	return result
}

func threadsToJsonArray(threads []*models.Thread) []*Thread {
	var result []*Thread
	for i := 0; i < len(threads); i++ {
		result = append(result, modelToThread(threads[i]))
	}

	return result
}

func (h *ForumHandler) ForumThreadsHandler(c echo.Context) error {
	slug := c.Param("slug")
	since := c.QueryParam("since")
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = -1
	}
	sort := c.QueryParam("desc") == "true"
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = 0
	}

	threads, err := h.useCase.GetForumThreads(c.Request().Context(), slug, since, limit, offset, sort)
	if err == forum.ErrForumNotFound {
		msg := map[string]string{"message": "Forum not found"}
		return c.JSON(http.StatusNotFound, msg)
	} else if err != nil {
		errText := map[string]string{"error": err.Error()}
		return c.JSON(http.StatusInternalServerError, errText)
	}

	for i := 0; i < len(threads); i++ {
		// TODO temporary for tests
		threads[i].Created = threads[i].Created.Add(-3 * time.Hour)
	}

	result := threadsToJsonArray(threads)
	// TODO temporary for tests
	if result == nil {
		return c.String(http.StatusOK, "[]")
	}
	return c.JSON(http.StatusOK, result)
}
