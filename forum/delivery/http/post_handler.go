package http

import (
	"dbProject/common"
	"dbProject/forum"
	userHttp "dbProject/user/delivery/http"
	"encoding/json"
	"io/ioutil"
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


type postInput struct {
	Message string `json:"message"`
}

func (h *PostHandler) ChangePostHandler(writer http.ResponseWriter, request *http.Request, ps map[string]string) {
	id, err := strconv.Atoi(ps["id"])
	if err != nil {
		http.Error(writer, "wrong ID", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	var input postInput
	err = json.Unmarshal(body, &input)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	post, err := h.useCase.ChangePost(request.Context(), id, input.Message)
	if err == forum.ErrPostNotFound {
		msg, _ := json.Marshal(map[string]string{"message": "Post not found"})
		common.WriteData(writer, http.StatusNotFound, msg)
		return
	} else if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(ModelToPost(post))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	common.WriteData(writer, http.StatusOK, data)
}

type DetailedInfo struct {
	PostInfo   Post                 `json:"post"`
	AuthorInfo *userHttp.UserOutput `json:"author,omitempty"`
	ThreadInfo *Thread              `json:"thread,omitempty"`
	ForumInfo  *ForumOutput         `json:"forum,omitempty"`
}

func (h *PostHandler) GetPostHandler(writer http.ResponseWriter, request *http.Request, ps map[string]string) {
	id, err := strconv.Atoi(ps["id"])
	if err != nil {
		http.Error(writer, "wrong ID", http.StatusBadRequest)
		return
	}

	related := strings.Split(request.FormValue("related"), ",")
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

	details, err := h.useCase.GetPostInfo(request.Context(), id, u, t, f)
	if err == forum.ErrPostNotFound {
		msg, _ := json.Marshal(map[string]string{"message": "Post not found"})
		common.WriteData(writer, http.StatusNotFound, msg)
		return
	} else if err == forum.ErrThreadNotFound {
		msg, _ := json.Marshal(map[string]string{"message": "Thread not found"})
		common.WriteData(writer, http.StatusNotFound, msg)
		return
	} else if err == forum.ErrForumNotFound {
		msg, _ := json.Marshal(map[string]string{"message": "Forum not found"})
		common.WriteData(writer, http.StatusNotFound, msg)
		return
	} else if err == forum.ErrUserNotFound {
		msg, _ := json.Marshal(map[string]string{"message": "User not found"})
		common.WriteData(writer, http.StatusNotFound, msg)
		return
	} else if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
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

	data, err := json.Marshal(info)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	common.WriteData(writer, http.StatusOK, data)
}
