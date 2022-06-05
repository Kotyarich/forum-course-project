package http

import (
	"forum-service/forum"
	"forum-service/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type ForumHandler struct {
	useCase forum.UseCaseForum
}


func NewForumHandler(useCase forum.UseCaseForum) *ForumHandler {
	return &ForumHandler{
		useCase: useCase,
	}
}


// Информация о форуме
// swagger:model Forum
type ForumOutput struct {
	// Общее количество сообщений в форуме
	//
	// read only: true
	// example: 200000
	Posts   int    `json:"posts"`

	// Человекопонятный URL, уникальное поле
	//
	// required: true
	// unique: true
	// example: pirate-stories
	Slug    string `json:"slug"`

	// Общее количество ветвей обсуждения в форуме
	//
	// read only: true
	// example: 200
	Threads int    `json:"threads"`

	// Название форума
	//
	// required: true
	// example: Pirate stories
	Title   string `json:"title"`

	// Nickname пользователя, создавшего форум
	// example: j.sparrow
	User    string `json:"user"`
}

// Информация о новом форуме
// swagger:model
type ForumInput struct {
	// Человекопонятный URL, уникальное поле
	//
	// required: true
	// unique: true
	// example: pirate-stories
	Slug  string `json:"slug"`

	// Название форума
	//
	// required: true
	// example: Pirate stories
	Title   string `json:"title"`

	// Nickname пользователя, создавшего форум
	// example: j.sparrow
	User    string `json:"user"`
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

	return c.JSON(http.StatusOK, usersToJsonArray(users))
}


type UserOutput struct {
	About    string `json:"about"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Nickname string `json:"nickname"`
}

func userToUserOutput(user *models.User) *UserOutput {
	return &UserOutput{
		Nickname: user.Nickname,
		Fullname: user.Fullname,
		Email:    user.Email,
		About:    user.About,
	}
}

func usersToJsonArray(users []*models.User) []*UserOutput {
	var result []*UserOutput
	for i := 0; i < len(users); i++ {
		result = append(result, userToUserOutput(users[i]))
	}

	return result
}

func forumsToJsonArray(forums []*models.Forum) []*ForumOutput {
	var result []*ForumOutput
	for i := 0; i < len(forums); i++ {
		result = append(result, forumToOutputFormat(forums[i]))
	}

	return result
}
