package http

import (
	"dbProject/forum"
	userHttp "dbProject/user/delivery/http"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
)

type PostHandler struct {
	useCase forum.UseCasePost
}


func NewPostHandler(useCase forum.UseCasePost) *PostHandler {
	return &PostHandler{
		useCase: useCase,
	}
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

// Полная информация о сообщени, включая связанные объекты
// swagger:model PostFull
type DetailedInfo struct {
	// "$ref": "#/definition/Post"
	PostInfo   Post                 `json:"post"`

	// "$ref": "#/definition/User"
	AuthorInfo *userHttp.UserOutput `json:"author,omitempty"`

	// "$ref": "#/definition/Thread"
	ThreadInfo *Thread              `json:"thread,omitempty"`

	// "$ref": "#/definition/Forum"
	ForumInfo  *ForumOutput         `json:"forum,omitempty"`
}

func (h *PostHandler) GetPostHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "wrong ID")
	}

	related := strings.Split(c.QueryParam("related"), ",")
	u := false
	f := false
	t := false
	for i := range related {
		switch related[i] {
		case "user":
			u = true
		case "forum":
			f = true
		case "thread":
			t = true
		}
	}

	details, err := h.useCase.GetPostInfo(c.Request().Context(), id, u, t, f)
	if err == forum.ErrPostNotFound {
		msg := map[string]string{"message": "Post not found"}
		return c.JSON(http.StatusNotFound, msg)
	} else if err == forum.ErrThreadNotFound {
		msg := map[string]string{"message": "Thread not found"}
		return c.JSON(http.StatusNotFound, msg)
	} else if err == forum.ErrForumNotFound {
		msg := map[string]string{"message": "Forum not found"}
		return c.JSON(http.StatusNotFound, msg)
	} else if err == forum.ErrUserNotFound {
		msg := map[string]string{"message": "User not found"}
		return c.JSON(http.StatusNotFound, msg)
	} else if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	info := new(DetailedInfo)
	info.PostInfo = *ModelToPost(&details.PostInfo)
	if u {
		info.AuthorInfo = userHttp.UserToUserOutput(details.AuthorInfo)
	}
	if t {
		info.ThreadInfo = modelToThread(details.ThreadInfo)
	}
	if f {
		info.ForumInfo = forumToOutputFormat(details.ForumInfo)
	}

	return c.JSON(http.StatusOK, info)
}
