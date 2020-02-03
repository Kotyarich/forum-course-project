package http

import (
	"dbProject/forum"
	"dbProject/models"
	userHttp "dbProject/user/delivery/http"
	"dbProject/utils"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	useCase forum.UseCase
}

func NewHandler(useCase forum.UseCase) *Handler {
	return &Handler{
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

func (h *Handler) ForumCreateHandler(writer http.ResponseWriter, request *http.Request, ps map[string]string) {
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
		utils.WriteData(writer, http.StatusNotFound, msg)
		return
	} else if err == forum.ErrForumAlreadyExists {
		data, _ := json.Marshal(forumToOutputFormat(f))
		utils.WriteData(writer, http.StatusConflict, data)
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

	utils.WriteData(writer, http.StatusCreated, data)
	return
}

type Thread struct {
	Author    string    `json:"author"`
	Slug      *string   `json:"slug"`
	Votes     int       `json:"votes"`
	Title     string    `json:"title"`
	Created   time.Time `json:"created"`
	ForumName string    `json:"forum"`
	Id        int       `json:"id"`
	Message   string    `json:"message"`
}

func threadToModel(t *Thread) *models.Thread {
	return &models.Thread{
		Author:    t.Author,
		Slug:      t.Slug,
		Title:     t.Title,
		Message:   t.Message,
		ForumName: t.ForumName,
		Id:        t.Id,
		Created:   t.Created,
		Votes:     t.Votes,
	}
}

func modelToThread(t *models.Thread) *Thread {
	return &Thread{
		Author:    t.Author,
		Slug:      t.Slug,
		Title:     t.Title,
		Message:   t.Message,
		ForumName: t.ForumName,
		Id:        t.Id,
		Created:   t.Created,
		Votes:     t.Votes,
	}
}

func (h *Handler) ThreadCreateHandler(writer http.ResponseWriter, request *http.Request, ps map[string]string) {
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
		utils.WriteData(writer, http.StatusNotFound, msg)
	} else if err == forum.ErrForumNotFound {
		msg, _ := json.Marshal(map[string]string{"message": "Forum not found"})
		utils.WriteData(writer, http.StatusNotFound, msg)
	} else if err == forum.ErrThreadAlreadyExists {
		data, err := json.Marshal(modelToThread(thread))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		utils.WriteData(writer, http.StatusConflict, data)
	} else if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	} else {
		data, err := json.Marshal(modelToThread(thread))

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		utils.WriteData(writer, http.StatusCreated, data)
	}
}

func (h *Handler) ForumDetailsHandler(writer http.ResponseWriter, r *http.Request, ps map[string]string) {
	slug := ps["slug"]

	f, err := h.useCase.GetForumDetails(r.Context(), slug)
	if err == forum.ErrForumNotFound {
		msg, _ := json.Marshal(map[string]string{"message": "404"})
		utils.WriteData(writer, http.StatusNotFound, msg)
		return
	} else if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(forumToOutputFormat(f))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	utils.WriteData(writer, http.StatusOK, data)
}

func (h *Handler) ForumUsersHandler(writer http.ResponseWriter, r *http.Request, ps map[string]string) {
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
		utils.WriteData(writer, http.StatusNotFound, msg)
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

	utils.WriteData(writer, http.StatusOK, result)
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

func (h *Handler) ForumThreadsHandler(writer http.ResponseWriter, r *http.Request, ps map[string]string) {
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
		utils.WriteData(writer, http.StatusNotFound, msg)
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

	utils.WriteData(writer, http.StatusOK, result)
}
