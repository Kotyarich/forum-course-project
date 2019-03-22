package routes

import (
	"database/sql"
	db2 "dbProject/db"
	"dbProject/models"
	"dbProject/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func SetUserRouter(router *mux.Router) {
	router.HandleFunc("/api/user/{nickname}/create", userCreateHandler)
	router.HandleFunc("/api/user/{nickname}/profile", userProfileHandler)
}

func printAll(rows *sql.Rows) {
	for rows.Next() {
		var user models.User
		i := 0
		rows.Scan(&i, &user.About, &user.Email, &user.Fullname, &user.Nickname)
		fmt.Println(user.Email)
	}
	fmt.Println()
}

func userCreateHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	body, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		http.Error(writer, err.Error(), 500)
	}

	var user models.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(writer, err.Error(), 500)
	}
	user.Nickname = vars["nickname"]
	db := db2.GetDB()
	_, err = db.Exec("INSERT INTO users (about, email, fullname, nickname) " +
		"VALUES ($1, $2, $3, $4)",
		user.About, user.Email, user.Fullname, user.Nickname)
	if err == nil {
		data, err := json.Marshal(user)
		if err != nil {
			http.Error(writer, err.Error(), 500)
		}
		writer.Header().Set("content-type", "application/json")
		writer.WriteHeader(201)
		writer.Write(data)
	} else {
		rows, err := db.Query("SELECT about, email, fullname, nickname " +
			"FROM users WHERE nickname = $1 OR email = $2",
			user.Nickname, user.Email)
		if err != nil {
			http.Error(writer, err.Error(), 500)
		}
		defer rows.Close()

		conflicts := []byte{'['}
		for rows.Next() {
			if len(conflicts) > 1 {
				conflicts = append(conflicts, ',')
			}
			var u models.User
			_ = rows.Scan(&u.About, &u.Email, &u.Fullname, &u.Nickname)
			data, err := json.Marshal(u)
			if err != nil {
				http.Error(writer, err.Error(), 500)
			}
			conflicts = append(conflicts, data...)
		}
		conflicts = append(conflicts, ']')

		writer.Header().Set("content-type", "application/json")
		writer.WriteHeader(409)
		writer.Write(conflicts)
	}
}

func userProfileHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		vars := mux.Vars(request)
		db := db2.GetDB()
		row := db.QueryRow("SELECT about, email, fullname, nickname " +
			"FROM users WHERE nickname = $1", vars["nickname"])

		var user models.User
		err := row.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
		if err != nil {
			msg, _ := json.Marshal(map[string]string{"message": "404"})
			utils.WriteData(writer, 404, msg)
		} else {
			data, err := json.Marshal(user)
			if err != nil {
				http.Error(writer, err.Error(), 500)
			}
			utils.WriteData(writer, 200, data)
		}
	}

	if request.Method == "POST" {
		postProfile(writer, request)
	}
}

func postProfile(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	db := db2.GetDB()
	// read body
	body, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		http.Error(writer, err.Error(), 500)
	}
	// parse body
	var user models.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(writer, err.Error(), 500)
	}
	user.Nickname = vars["nickname"]
	// get current data
	var oldUser models.User
	err = db.QueryRow("SELECT about, email, fullname, nickname FROM users " +
		"WHERE nickname = $1", user.Nickname).Scan(&oldUser.About, &oldUser.Email, &oldUser.Fullname, &oldUser.Nickname)
	if err != nil {
		msg, _ := json.Marshal(map[string]string{"message": "User not found"})
		utils.WriteData(writer, 404, msg)
		return
	}
	// check empty request
	if user.Email == "" && user.Fullname == "" && user.About == "" {
		data, err := json.Marshal(oldUser)
		if err != nil {
			http.Error(writer, err.Error(), 500)
		}
		utils.WriteData(writer, 200, data)
		return
	}
	// check empty fields
	if user.Fullname == "" {
		user.Fullname = oldUser.Fullname
	}
	if user.Email == "" {
		user.Email = oldUser.Email
	}
	if user.About == "" {
		user.About = oldUser.About
	}

	result, err := db.Exec("UPDATE users " +
		"SET about = $1, email = $2, fullname = $3 " +
		"WHERE  nickname = $4", user.About, user.Email, user.Fullname, user.Nickname)
	// user with new email already exist
	if err != nil {
		msg, _ := json.Marshal(map[string]string{"message": "conflict"})
		utils.WriteData(writer, 409, msg)
		return
	}

	number, _ := result.RowsAffected()
	if number == 0 {
		msg, _ := json.Marshal(map[string]string{"message": "User not found"})
		utils.WriteData(writer, 404, msg)
		return
	} else {
		data, err := json.Marshal(user)
		if err != nil {
			http.Error(writer, err.Error(), 500)
		}
		utils.WriteData(writer, 200, data)
		return
	}
}
