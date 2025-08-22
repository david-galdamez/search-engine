package models

type IndexRequest struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text,omitempty"`
	Url   string `json:"url,omitempty"`
}
