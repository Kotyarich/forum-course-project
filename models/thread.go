package models

type Thread struct {
	Author string `json:"author"`
	Slug string `json:"slug"`
	Votes int `json:"votes"`
	Title string `json:"title"`
	Created string `json:"created"`
	ForumName string `json:"forum"`
	Id int `json:"id"`
	Message string `json:"message"`
}

type ThreadResult struct {
	Author string `json:"author"`
	Title string `json:"title"`
	Created string `json:"created"`
	ForumName string `json:"forum"`
	Id int `json:"id"`
	Message string `json:"message"`
}

type ThreadUpdate struct {
	Message string `json:"message"`
	Title string `json:"title"`
}