package http

import (
	"dbProject/forum"
	"dbProject/models"
	"dbProject/utils"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

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
	Id        int       `json:"id"`
	IsEdited  bool      `json:"isEdited"`
	Message   string    `json:"message"`
	Parent    int       `json:"parent"`
	Tid       int       `json:"thread"`
}

func toModelPost(p *Post) *models.Post {
	return &models.Post{
		Author:    p.Author,
		Created:   p.Created,
		ForumName: p.ForumName,
		Id:        p.Id,
		IsEdited:  p.IsEdited,
		Message:   p.Message,
		Parent:    p.Parent,
		Tid:       p.Tid,
	}
}

func modelToPost(p *models.Post) *Post {
	return &Post{
		Author:    p.Author,
		Created:   p.Created,
		ForumName: p.ForumName,
		Id:        p.Id,
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
		posts = append(posts, *modelToPost(p[i]))
	}

	return posts
}

func (h *Handler) ThreadPostCreateHandler(writer http.ResponseWriter, request *http.Request, ps map[string]string) {
	body, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}

	var postsInput []Post
	err = json.Unmarshal(body, &postsInput)
	if err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}

	slug := ps["slug"]
	var posts []*models.Post
	for i := 0; i < len(postsInput); i++ {
		posts = append(posts, toModelPost(&postsInput[i]))
	}
	posts, err = h.useCase.CreateThreadPost(request.Context(), slug, posts)

	if err == forum.ErrThreadNotFound {
		msg, _ := json.Marshal(map[string]string{"message": "Thread not found"})
		utils.WriteData(writer, http.StatusNotFound, msg)
		return
	} else if err == forum.ErrUserNotFound {
		msg, _ := json.Marshal(map[string]string{"message": "User not found"})
		utils.WriteData(writer, http.StatusNotFound, msg)
		return
	} else if err == forum.ErrWrongParentsThread {
		msg, _ := json.Marshal(map[string]string{"message": "Parent in another thread"})
		utils.WriteData(writer, http.StatusConflict, msg)
		return
	} else if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(modelsToPostsArray(posts))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteData(writer, http.StatusCreated, data)
}

func (h *Handler) GetThreadHandler(writer http.ResponseWriter, request *http.Request, ps map[string]string) {
	slug := ps["slug"]

	thread, err := h.useCase.GetThread(request.Context(), slug)
	if err == forum.ErrThreadNotFound {
		msg, _ := json.Marshal(map[string]string{"message": "Thread not found"})
		utils.WriteData(writer, http.StatusNotFound, msg)
		return
	} else if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(modelToThread(thread))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	utils.WriteData(writer, http.StatusOK, data)
}

func (h *Handler) PostThreadHandler(writer http.ResponseWriter, request *http.Request, ps map[string]string) {
	slug := ps["slug"]

	body, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	// parse input
	var input ThreadUpdate
	err = json.Unmarshal(body, &input)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	thread, err := h.useCase.ChangeThread(request.Context(), slug, input.Title, input.Message)
	if err == forum.ErrThreadNotFound {
		msg, _ := json.Marshal(map[string]string{"message": "Thread not found"})
		utils.WriteData(writer, http.StatusNotFound, msg)
		return
	} else if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(modelToThread(thread))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteData(writer, http.StatusOK, data)
}

func (h *Handler) GetThreadPosts(writer http.ResponseWriter, r *http.Request, ps map[string]string) {
	slug := ps["slug"]

	since, err := strconv.Atoi(r.FormValue("since"))
	if err != nil {
		since = -1
	}
	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		limit = -1
	}
	desc := r.FormValue("desc") == "true"
	var sort models.PostSortType
	switch r.FormValue("sort") {
	case "flat", "":
		sort = models.Flat
	case "tree":
		sort = models.Tree
	case "parent_tree":
		sort = models.ParentTree
	}

	posts, err := h.useCase.GetThreadPosts(r.Context(), slug, limit, since, desc, sort)
	if err == forum.ErrThreadNotFound {
		msg, _ := json.Marshal(map[string]string{"message": "Thread not found"})
		utils.WriteData(writer, http.StatusNotFound, msg)
		return
	} else if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	result := []byte("[ ")
	for i := 0; i < len(posts); i++ {
		if len(result) > 2 {
			result = append(result, ',')
		}

		data, err := json.Marshal(modelToPost(posts[i]))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		result = append(result, data...)
	}
	result = append(result, ']')
	utils.WriteData(writer, http.StatusOK, result)
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

func (h *Handler) ThreadVoteHandler(writer http.ResponseWriter, request *http.Request, ps map[string]string) {
	slug := ps["slug"]

	body, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	var vote Vote
	err = json.Unmarshal(body, &vote)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	thread, err := h.useCase.VoteForThread(request.Context(), slug, toModelVote(&vote))
	if err == forum.ErrThreadNotFound {
		msg, _ := json.Marshal(map[string]string{"message": "Thread not found"})
		utils.WriteData(writer, http.StatusNotFound, msg)
		return
	} else if err == forum.ErrUserNotFound {
		msg, _ := json.Marshal(map[string]string{"message": "User not found"})
		utils.WriteData(writer, http.StatusNotFound, msg)
		return
	} else if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(modelToThread(thread))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteData(writer, http.StatusOK, data)
}
