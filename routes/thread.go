package routes

import (
	db2 "dbProject/db"
	"dbProject/models"
	"dbProject/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

func SetThreadRouter(router *mux.Router) {
	router.HandleFunc("/api/thread/{slug}/create", threadCreateHandler)
	router.HandleFunc("/api/thread/{slug}/details", detailsHandler)
	router.HandleFunc("/api/thread/{slug}/posts", getPosts)
	router.HandleFunc("/api/thread/{slug}/vote", threadVoteHandler)
}

func threadHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("hi /thread")
}

func detailsHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		slug := mux.Vars(request)["slug"]
		id, _ := strconv.Atoi(slug)

		db := db2.GetDB()
		row := db.QueryRow("SELECT * " +
			"FROM threads WHERE slug = $1 OR id = $2;", slug, id)

		var thread models.Thread
		err := row.Scan(&thread.Author, &thread.Created, &thread.ForumName, &thread.Id,
			&thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
		if err != nil {
			msg, _ := json.Marshal(map[string]string{"message": "Thread not found"})
			utils.WriteData(writer, 404, msg)
			return
		}

		data, err := json.Marshal(thread)
		if err != nil {
			http.Error(writer, err.Error(), 500)
		}
		utils.WriteData(writer, 200, data)
	}
}

func threadCreateHandler(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}
	// parse input
	var posts []models.Post
	err = json.Unmarshal(body, &posts)
	if err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}

	for i := 0; i < len(posts); i++ {
		err = createPost(&posts[i], writer, request)
		if err != nil {
			break
		}
	}

	if err == nil {
		data, err := json.Marshal(posts)
		if err != nil {
			http.Error(writer, err.Error(), 500)
			return
		}

		utils.WriteData(writer, 201, data)
	}
}

func createPost(post *models.Post, writer http.ResponseWriter, request *http.Request) error {
	slug := mux.Vars(request)["slug"]
	tid, _ := strconv.Atoi(slug)
	db := db2.GetDB()
	// check user
	err := db.QueryRow("SELECT nickname " +
		"FROM users WHERE nickname = $1", post.Author).Scan(&post.Author)
	if err != nil {
		msg, _ := json.Marshal(map[string]string{"message": "User not found"})
		utils.WriteData(writer, 404, msg)
		return err
	}
	// check thread and get forum
	err = db.QueryRow("SELECT id, forum " +
		"FROM threads WHERE slug = $1 OR id = $2", slug, tid).Scan(&post.Tid, &post.ForumName)
	if err != nil {
		msg, _ := json.Marshal(map[string]string{"message": "Thread not found"})
		utils.WriteData(writer, 404, msg)
		return err
	}
	// check parent
	if post.Parent != 0 {
		err = db.QueryRow("SELECT id "+
			"FROM posts WHERE tid = $1 AND id = $2", post.Tid, post.Parent).Scan(&post.Parent)
		if err != nil {
			fmt.Println(err)
			msg, _ := json.Marshal(map[string]string{"message": "Parent not found"})
			utils.WriteData(writer, 409, msg)
			return err
		}
	}
	if post.Created == "" {
		post.Created = "1970-01-01T00:00:00.000Z"
	}
	err = db.QueryRow("INSERT INTO posts (author, forum, message, parent, tid, created, slug) " +
		"VALUES ($1, $2, $3, $4, $5, $6, " +
		"(SELECT slug FROM posts WHERE id = $4) || (SELECT currval('posts_id_seq')::integer)) " +
		"RETURNING id", post.Author, post.ForumName, post.Message,
		post.Parent, post.Tid, post.Created).Scan(&post.Id)
	if err != nil {
		http.Error(writer, err.Error(), 500)
		return err
	}

	var postSlug string
	if post.Parent != 0 {
		idStr := strconv.Itoa(post.Id)
		parentStr := strconv.Itoa(post.Parent)
		postSlug += parentStr
		for i := 0; i < 32 - len(parentStr) - len(idStr); i++ {
			postSlug += "0"
		}
		postSlug += idStr
	} else {
		postSlug = strconv.Itoa(post.Id)
	}
	_, err = db.Exec("UPDATE posts SET slug = $1 WHERE id = $2", postSlug, post.Id)

	return nil
}

func threadVoteHandler(writer http.ResponseWriter, request *http.Request) {
	slug := mux.Vars(request)["slug"]
	id, _ := strconv.Atoi(slug)

	body, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}
	// parse input
	var vote models.Vote
	err = json.Unmarshal(body, &vote)
	if err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}

	var thread models.Thread
	db := db2.GetDB()
	transaction, _ := db.Begin()
	// check thread and get forum
	err = transaction.QueryRow("SELECT author, created, id, forum, message, slug, title, votes " +
		"FROM threads WHERE slug = $1 OR id = $2", slug, id).Scan(
			&thread.Author, &thread.Created, &thread.Id, &thread.ForumName,
			&thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	if err != nil {
		_ = transaction.Rollback()
		msg, _ := json.Marshal(map[string]string{"message": "Thread not found"})
		utils.WriteData(writer, 404, msg)
		return
	}
	// create/update vote
	rows, _ := transaction.Query("SELECT * FROM votes")
	for rows.Next() {
		var id int
		var v models.Vote
		rows.Scan(&v.Nickname, &id, &v.Voice)
	}
	r, err := transaction.Exec("UPDATE votes SET voice=$1 " +
		"WHERE tid=$2 AND nickname=$3;", vote.Voice, thread.Id, vote.Nickname)
	if count, _ := r.RowsAffected(); count == 0 {
		_, err := transaction.Exec("INSERT INTO votes (nickname, tid, voice)" +
			"VALUES ($1, $2, $3);", vote.Nickname, thread.Id, vote.Voice)
		if err != nil {
			_ = transaction.Rollback()
			http.Error(writer, err.Error(), 500)
			return
		}
	}
	// get new votes
	err = transaction.QueryRow("SELECT votes FROM threads " +
		"WHERE id = $1", thread.Id).Scan(&thread.Votes)
	if err != nil {
		_ = transaction.Rollback()
		http.Error(writer, err.Error(), 500)
		return
	}

	data, err := json.Marshal(thread)
	if err != nil {
		_ = transaction.Rollback()
		http.Error(writer, err.Error(), 500)
		return
	}
	_ = transaction.Commit()
	utils.WriteData(writer, 200, data)
}

type postsInput struct {
	Slug string
	Id int
	ParentId int
	Limit string
	Since string
	Sort string
	Desc bool
}

func getFlatPosts(input postsInput, writer http.ResponseWriter, r *http.Request) error {
	db := db2.GetDB()

	query := "SELECT id, author, created, forum, isEdited, message, parent, tid " +
		"FROM posts WHERE tid = $1 "
	if input.Since != ""{
		if input.Desc {
			query += "AND id < '" + input.Since + "' "
		} else {
			query += "AND id > '" + input.Since + "' "
		}
	}
	query += "ORDER BY id "
	if input.Desc {
		query += "DESC "
	}
	if limit := r.FormValue("limit"); limit != "" {
		query += "LIMIT " + limit + " "
	}

	rows, err := db.Query(query, input.Id)
	if err != nil {
		http.Error(writer, err.Error(), 500)
		return err
	}

	result := []byte("[ ")
	for rows.Next() {
		if len(result) > 2 {
			result = append(result, ',')
		}

		post := models.Post{}
		err = rows.Scan(&post.Id, &post.Author, &post.Created, &post.ForumName,
			&post.IsEdited, &post.Message, &post.Parent, &post.Tid)
		if err != nil {
			http.Error(writer, err.Error(), 500)
			return err
		}

		data, _ := json.Marshal(post)
		result = append(result, data...)
	}
	result = append(result, ']')

	utils.WriteData(writer, 200, result)
	return nil
}

func getTreePosts(input postsInput, writer http.ResponseWriter, r *http.Request) error {
	db := db2.GetDB()

	query := "SELECT id, author, created, forum, isEdited, message, parent, tid, slug" +
		" FROM posts WHERE tid = $1 "
	if input.Since != ""{
		if input.Desc {
			query += "AND slug < (SELECT slug FROM posts WHERE id = " + input.Since + ") "
		} else {
			query += "AND slug > (SELECT slug FROM posts WHERE id = " + input.Since + ") "
		}
	}
	query += "ORDER BY slug "
	if input.Desc {
		query += "DESC "
	}
	if limit := r.FormValue("limit"); limit != "" {
		query += "LIMIT " + limit + " "
	}

	rows, err := db.Query(query, input.Id)
	if err != nil {
		http.Error(writer, err.Error(), 500)
		return err
	}

	result := []byte("[ ")
	for rows.Next() {
		if len(result) > 2 {
			result = append(result, ',')
		}

		post := models.Post{}
		slug := ""
		err = rows.Scan(&post.Id, &post.Author, &post.Created, &post.ForumName,
			&post.IsEdited, &post.Message, &post.Parent, &post.Tid, &slug)
		if err != nil {
			http.Error(writer, err.Error(), 500)
			return err
		}

		data, _ := json.Marshal(post)
		result = append(result, data...)
	}
	result = append(result, ']')

	utils.WriteData(writer, 200, result)
	return nil
}

func getParentTreePosts(input postsInput, writer http.ResponseWriter, r *http.Request) error {
	db := db2.GetDB()

	query := "WITH roots AS ( " +
		"SELECT id FROM posts WHERE tid = $1 AND parent = 0 "
		//"SELECT id, author, created, forum, isEdited, message, parent, tid " +
		//"FROM posts WHERE tid = $1 AND parent = 0 "
	if limit := r.FormValue("limit"); limit != "" {
		query += "LIMIT " + limit + " "
	}
	query += ") SELECT posts.id, author, created, forum, isEdited, message, parent, tid " +
		"FROM posts JOIN roots ON roots.id = rootId "
	if input.Since != ""{
		if input.Desc {
			query += "AND posts.id < " + input.Since + " "
		} else {
			query += "AND posts.id > " + input.Since + " "
		}
	}
	query += "ORDER BY rootId"
	if input.Desc {
		query += " DESC"
	}
	query += ", slug"

	rows, err := db.Query(query, input.Id)
	if err != nil {
		fmt.Println(err)
		http.Error(writer, err.Error(), 500)
		return err
	}

	result := []byte("[ ")
	for rows.Next() {
		if len(result) > 2 {
			result = append(result, ',')
		}

		post := models.Post{}
		//slug := ""
		err = rows.Scan(&post.Id, &post.Author, &post.Created, &post.ForumName,
			&post.IsEdited, &post.Message, &post.Parent, &post.Tid)
		if err != nil {
			http.Error(writer, err.Error(), 500)
			return err
		}

		data, _ := json.Marshal(post)
		result = append(result, data...)
	}
	result = append(result, ']')

	utils.WriteData(writer, 200, result)
	return nil
}

func getChilPosts(input postsInput, writer http.ResponseWriter, request *http.Request) ([]byte, error) {
	db := db2.GetDB()


	query := "SELECT id, author, created, forum, isEdited, message, parent, tid " +
		"FROM posts WHERE tid = $1 AND parent = $2 "
	if input.Since != ""{
		if input.Desc {
			query += "AND id < '" + input.Since + "' "
		} else {
			query += "AND id > '" + input.Since + "' "
		}
	}
	query += "ORDER BY id "
	if input.Desc {
		query += "DESC "
	}
	if input.Limit != "" {
		query += "LIMIT " + input.Limit + " "
	}

	rows, err := db.Query(query, input.Id, input.ParentId)
	if err != nil {
		http.Error(writer, err.Error(), 500)
		return []byte{}, err
	}

	var result []byte
	for rows.Next() {
		result = append(result, ',')

		post := models.Post{}
		err = rows.Scan(&post.Id, &post.Author, &post.Created, &post.ForumName,
			&post.IsEdited, &post.Message, &post.Parent, &post.Tid)
		if err != nil {
			http.Error(writer, err.Error(), 500)
			return []byte{}, err
		}

		data, _ := json.Marshal(post)
		result = append(result, data...)
	}

	return result, nil
}


func getPosts(writer http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var input postsInput
		input.Slug = mux.Vars(r)["slug"]
		input.Id, _ = strconv.Atoi(input.Slug)
		// check if forum exist
		db := db2.GetDB()
		err := db.QueryRow("SELECT id FROM threads " +
			"WHERE slug = $1 OR id = $2", input.Slug, input.Id).Scan(&input.Id)
		if err != nil {
			msg, _ := json.Marshal(map[string]string{"message": "Thread not found"})
			utils.WriteData(writer, 404, msg)
			return
		}

		// get params
		input.Since = r.FormValue("since")
		input.Desc = r.FormValue("desc") == "true"
		input.Sort = r.FormValue("sort")
		input.Limit = r.FormValue("limit")
		switch input.Sort {
		case "flat", "":
			err = getFlatPosts(input, writer, r)
		case "tree":
			err = getTreePosts(input, writer, r)
		case "parent_tree":
			err = getParentTreePosts(input, writer, r)
		}
	}
}
