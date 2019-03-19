package routes

import (
	db2 "dbProject/db"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"dbProject/models"
	"io/ioutil"
	"net/http"
)

func SetUserRouter(router *mux.Router) {
	router.HandleFunc("/api/user/{nickname}/create", userCreateHandler)
	router.HandleFunc("/api/user/{nickname}/profile", userProfileHandler)
}

func userCreateHandler(writer http.ResponseWriter, request *http.Request) {
	request.Header.Set("Content-Type", "application/json")
	vars := mux.Vars(request)
	fmt.Println("hi /", vars["nickname"])

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
	fmt.Println(err)
}

func userProfileHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("hi /")
}