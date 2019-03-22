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
)

func SetForumRouter(router *mux.Router) {
	router.HandleFunc("/api/forum/create", createHandler)
	router.HandleFunc("/api/forum/{slug}/create", slugCreateHandler)
	router.HandleFunc("/api/forum/{slug}/details", getForum)
	router.HandleFunc("/api/forum/{slug}/threads", getThreads)
	router.HandleFunc("/api/forum/{slug}/{type}", slugHandler)
}

func createHandler(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		http.Error(writer, err.Error(), 500)
	}
	// parse input
	var input models.ForumInput
	err = json.Unmarshal(body, &input)
	if err != nil {
		http.Error(writer, err.Error(), 500)
	}

	db := db2.GetDB()
	err = db.QueryRow("SELECT nickname " +
		"FROM users WHERE nickname = $1", input.User).Scan(&input.User)
	if err != nil {
		msg, _ := json.Marshal(map[string]string{"message": "User not found"})
		utils.WriteData(writer, 404, msg)
		return
	}

	_, err = db.Exec("INSERT INTO forums (slug, title, author) " +
		"VALUES ($1, $2, $3)", input.Slug, input.Title, input.User)
	if err == nil {
		data, err := json.Marshal(input)
		if err != nil {
			http.Error(writer, err.Error(), 500)
		}

		utils.WriteData(writer, 201, data)
		return
	} else {
		row := db.QueryRow("SELECT * FROM forums WHERE slug = $1", input.Slug)

		var f models.Forum
		var id int
		err = row.Scan(&id, &f.Posts, &f.Slug, &f.Threads, &f.Title, &f.User)
		if err != nil {
			http.Error(writer, err.Error(), 500)
		}

		data, _ := json.Marshal(f)
		utils.WriteData(writer, 409, data)
	}
}

func slugHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("hi /forum")
	if request.Method == "POST" {
		//
	}
}

func slugCreateHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		slug := mux.Vars(request)["slug"]

		body, err := ioutil.ReadAll(request.Body)
		defer request.Body.Close()
		if err != nil {
			http.Error(writer, err.Error(), 500)
		}
		// parse input
		var thread models.Thread
		err = json.Unmarshal(body, &thread)
		if err != nil {
			http.Error(writer, err.Error(), 500)
		}

		db := db2.GetDB()
		err = db.QueryRow("SELECT nickname " +
			"FROM users WHERE nickname = $1", thread.Author).Scan(&thread.Author)
		if err != nil {
			msg, _ := json.Marshal(map[string]string{"message": "User not found"})
			utils.WriteData(writer, 404, msg)
			return
		}
		// Check forum
		err = db.QueryRow("SELECT slug "+
			"FROM forums WHERE slug = $1", slug).Scan(&thread.ForumName)
		if err != nil {
			msg, _ := json.Marshal(map[string]string{"message": "Forum not found"})
			utils.WriteData(writer, 404, msg)
			return
		}
		// AAAAAAAAAAAAAAAAA
		if len(thread.Created) > 0 {
			if len(thread.Slug) > 0 {
				err = db.QueryRow("INSERT INTO threads (author, created, forum, message, title, slug) "+
					"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", thread.Author, thread.Created,
					thread.ForumName, thread.Message, thread.Title, thread.Slug).Scan(&thread.Id)
			} else {
				err = db.QueryRow("INSERT INTO threads (author, created, forum, message, title) "+
					"VALUES ($1, $2, $3, $4, $5) RETURNING id", thread.Author, thread.Created,
					thread.ForumName, thread.Message, thread.Title).Scan(&thread.Id)
			}
		} else {
			if len(thread.Slug) > 0 {
				err = db.QueryRow("INSERT INTO threads (author, forum, message, title, slug) "+
					"VALUES ($1, $2, $3, $4, $5) RETURNING id", thread.Author,
					thread.ForumName, thread.Message, thread.Title, thread.Slug).Scan(&thread.Id)
			} else {
				err = db.QueryRow("INSERT INTO threads (author, forum, message, title) "+
					"VALUES ($1, $2, $3, $4) RETURNING id", thread.Author,
					thread.ForumName, thread.Message, thread.Title).Scan(&thread.Id)
			}
		}
		if err == nil {
			data, err := json.Marshal(thread)
			if err != nil {
				http.Error(writer, err.Error(), 500)
				return
			}

			utils.WriteData(writer, 201, data)
			return
		} else {
			row := db.QueryRow("SELECT * FROM threads WHERE slug = $1", thread.Slug)

			var thr models.Thread
			err = row.Scan(&thr.Author, &thr.Created, &thr.ForumName, &thr.Id,
				&thr.Message, &thr.Slug, &thr.Title, &thr.Votes)
			if err != nil {
				http.Error(writer, err.Error(), 500)
				return
			}

			data, _ := json.Marshal(thr)
			utils.WriteData(writer, 409, data)
		}
	}
}

func getForum(writer http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		vars := mux.Vars(r)
		db := db2.GetDB()
		row := db.QueryRow("SELECT posts, slug, threads, title, author " +
			"FROM forums WHERE slug = $1", vars["slug"])

		var forum models.Forum
		err := row.Scan(&forum.Posts, &forum.Slug, &forum.Threads, &forum.Title, &forum.User)
		if err != nil {
			msg, _ := json.Marshal(map[string]string{"message": "404"})
			utils.WriteData(writer, 404, msg)
			return
		}

		data, err := json.Marshal(forum)
		if err != nil {
			http.Error(writer, err.Error(), 500)
		}
		utils.WriteData(writer, 200, data)
	}
}

func getThreads(writer http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		slug := mux.Vars(r)["slug"]
		// check if forum exist
		db := db2.GetDB()
		var forum models.Forum
		err := db.QueryRow("SELECT slug FROM forums WHERE slug = $1", slug).Scan(&forum.Slug)
		if err != nil {
			msg, _ := json.Marshal(map[string]string{"message": "Forum not found"})
			utils.WriteData(writer, 404, msg)
			return
		}

		// form query
		query := "SELECT * FROM threads WHERE forum = $1 "
		since := r.FormValue("since")
		sort := r.FormValue("desc")
		if since != ""{
			if sort == "true" {
				query += "AND created <= '" + since + "' "
			} else {
				query += "AND created >= '" + since + "' "
			}
		}
		query += "ORDER BY created "
		if sort != "" && sort != "false" {
			query += "DESC "
		}
		if limit := r.FormValue("limit"); limit != "" {
			query += "LIMIT " + limit + " "
		}

		rows, err := db.Query(query, slug)
		if err != nil {
			http.Error(writer, err.Error(), 500)
			return
		}

		result := []byte("[ ")
		for rows.Next() {
			if len(result) > 2 {
				result = append(result, ',')
			}

			thr := models.Thread{}
			err = rows.Scan(&thr.Author, &thr.Created, &thr.ForumName, &thr.Id,
				&thr.Message, &thr.Slug, &thr.Title, &thr.Votes)
			if err != nil {
				http.Error(writer, err.Error(), 500)
				return
			}

			data, _ := json.Marshal(thr)
			result = append(result, data...)
		}
		result = append(result, ']')

		utils.WriteData(writer, 200, result)
	}
}