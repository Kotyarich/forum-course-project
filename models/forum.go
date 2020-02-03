package models

type Forum struct {
	Title   string
	Slug    string
	User    string
	Threads int
	Posts   int
}

type Forum1 struct {
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
