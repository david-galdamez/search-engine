package models

type SearchResponse struct {
	Results      []SearchResults `json:"results"`
	Query        string          `json:"query"`
	TotalResults int             `json:"total_results"`
}

type SearchResults struct {
	DocId string  `json:"doc_id"`
	Title string  `json:"title"`
	Score float64 `json:"score"`
}

type DocsScore map[string]float64
