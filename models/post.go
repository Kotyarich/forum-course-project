package models

type Post struct {
	Author string `json:"author"`
	Created string `json:"created"`
	ForumName string `json:"forum"`
	Id int `json:"id"`
	IsEdited bool `json:"isEdited"`
	Message string `json:"message"`
	Parent int `json:"parent"`
	Tid int `json:"thread"`
}
