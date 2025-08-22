package models

type Document struct {
	Id      string  `json:"id"`
	Title   string  `json:"title"`
	Length  int     `json:"length"`
	Content string  `json:"content"`
	Url     *string `json:"url,omitempty"`
}
