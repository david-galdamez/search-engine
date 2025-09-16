package models

type Document struct {
	Id      string  `json:"id"`
	Title   string  `json:"title"`
	Content string  `json:"content,omitempty"`
	Url     *string `json:"url,omitempty"`
	Length  int     `json:"length"`
}
