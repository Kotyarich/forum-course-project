package http

import (
	"dbProject/common"
	"dbProject/forum"
	"dbProject/models"
	userHttp "dbProject/user/delivery/http"
	"encoding/json"
	"io/ioutil"
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

func (h *ForumHandler) ForumCreateHandler(writer http.ResponseWriter, request *http.Request, ps map[string]string) {
	body, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	var input ForumInput
	err = json.Unmarshal(body, &input)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	f, err := h.useCase.CreateForum(request.Context(), forumInputToModel(input))
	if err == forum.ErrUserNotFound {
		msg, _ := json.Marshal(map[string]string{"message": "User not found"})
		common.WriteData(writer, http.StatusNotFound, msg)
		return
	} else if err == forum.ErrForumAlreadyExists {
		data, _ := json.Marshal(forumToOutputFormat(f))
		common.WriteData(writer, http.StatusConflict, data)
		return
	} else if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(forumToOutputFormat(f))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	common.WriteData(writer, http.StatusCreated, data)
	return
}

func (h *ForumHandler) ThreadCreateHandler(writer http.ResponseWriter, request *http.Request, ps map[string]string) {
	slug := ps["slug"]

	body, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	var threadInp Thread
	err = json.Unmarshal(body, &threadInp)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	thread, err := h.useCase.CreateForumThread(request.Context(), slug, threadToModel(&threadInp))
	if err == forum.ErrUserNotFound {
		msg, _ := json.Marshal(map[string]string{"message": "User not found"})
		common.WriteData(writer, http.StatusNotFound, msg)
	} else if err == forum.ErrForumNotFound {
		msg, _ := json.Marshal(map[string]string{"message": "Forum not found"})
		common.WriteData(writer, http.StatusNotFound, msg)
	} else if err == forum.ErrThreadAlreadyExists {
		data, err := json.Marshal(modelToThread(thread))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		common.WriteData(writer, http.StatusConflict, data)
	} else if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	} else {
		data, err := json.Marshal(modelToThread(thread))

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		common.WriteData(writer, http.StatusCreated, data)
	}
}

func (h *ForumHandler) ForumsHandler(writer http.ResponseWriter, r *http.Request, ps map[string]string) {
	forums, err := h.useCase.GetForums(r.Context())
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := forumsToJsonArray(forums)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	common.WriteData(writer, http.StatusOK, data)
}

func (h *ForumHandler) ForumDetailsHandler(writer http.ResponseWriter, r *http.Request, ps map[string]string) {
	slug := ps["slug"]

	f, err := h.useCase.GetForumDetails(r.Context(), slug)
	if err == forum.ErrForumNotFound {
		msg, _ := json.Marshal(map[string]string{"message": "404"})
		common.WriteData(writer, http.StatusNotFound, msg)
		return
	} else if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(forumToOutputFormat(f))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	common.WriteData(writer, http.StatusOK, data)
}

func (h *ForumHandler) ForumUsersHandler(writer http.ResponseWriter, r *http.Request, ps map[string]string) {
	slug := ps["slug"]
	since := r.FormValue("since")
	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		limit = -1
	}
	sort := r.FormValue("desc") == "true"

	users, err := h.useCase.GetForumUsers(r.Context(), slug, since, limit, sort)
	if err == forum.ErrForumNotFound {
		msg, _ := json.Marshal(map[string]string{"message": "Forum not found"})
		common.WriteData(writer, http.StatusNotFound, msg)
		return
	} else if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := userHttp.UsersToJsonArray(users)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	common.WriteData(writer, http.StatusOK, result)
}

func forumsToJsonArray(forums []*models.Forum) ([]byte, error) {
	result := []byte{'['}
	for i := 0; i < len(forums); i++ {
		if len(result) > 1 {
			result = append(result, ',')
		}

		forumOutput := forumToOutputFormat(forums[i])
		data, err := json.Marshal(forumOutput)
		if err != nil {
			return nil, err
		}

		result = append(result, data...)
	}
	result = append(result, ']')

	return result, nil
}

func threadsToJsonArray(threads []*models.Thread) ([]byte, error) {
	result := []byte{'['}
	for i := 0; i < len(threads); i++ {
		if len(result) > 1 {
			result = append(result, ',')
		}

		threadOutput := modelToThread(threads[i])
		data, err := json.Marshal(threadOutput)
		if err != nil {
			return nil, err
		}

		result = append(result, data...)
	}
	result = append(result, ']')

	return result, nil
}

func (h *ForumHandler) ForumThreadsHandler(writer http.ResponseWriter, r *http.Request, ps map[string]string) {
	slug := ps["slug"]
	since := r.FormValue("since")
	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		limit = -1
	}
	sort := r.FormValue("desc") == "true"

	threads, err := h.useCase.GetForumThreads(r.Context(), slug, since, limit, sort)
	if err == forum.ErrForumNotFound {
		msg, _ := json.Marshal(map[string]string{"message": "Forum not found"})
		common.WriteData(writer, http.StatusNotFound, msg)
		return
	} else if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := threadsToJsonArray(threads)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	common.WriteData(writer, http.StatusOK, result)
}
