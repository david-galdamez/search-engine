# 🔎 Search Engine in Go

A basic **search engine** written in **Go**, that indexes text documents or web pages and allows performing searches using **TF-IDF** to rank results.

It uses **BoltDB** as an embedded database to store documents, terms, and index metadata.

---

## ✨ Features
- Indexing plain text documents or from URLs.
- Tokenization with normalization (removes accents, stopwords).
- **TF-IDF** calculation for ranking results.
- REST API with endpoints for indexing and searching.
- Scraper with [goquery](https://github.com/PuerkitoBio/goquery) to extract the main content of web pages.
- Results sorted by relevance.

---

## 🛠️ Installation

Clone the repository:

```bash
git clone https://github.com/david-galdamez/search-engine.git
cd search-engine
```

Install dependencies:

```bash
go mod tidy
```

Build the project

```bash
go build
```

Run the project

```bash
./search-engine
```

## 🚀 Usage
The engine exposes a REST API with the following endpoints:

### Index documents

**POST**   `/index` 
Indexes a plain text document or a web page.

**Request body:**
```json
[
  {
    "id": "doc1",
    "title": "My First Document",
    "text": "This is an example of plain text indexing"
  },
  {
    "id": "doc2",
    "title": "Go Tutorial",
    "url": "https://golang.org"
  }
]
```
### Search

**GET** `/search?q="query"` 
Performs a search using TF-IDF ranking.

**Response example:**
```json
{
	"results": [
		{
			"doc_id": "doc2",
			"title": "Google",
			"score": 0.007171275876963701
		},
		{
			"doc_id": "doc3",
			"title": "C++ the programming language",
			"score": 0.002914025565134647
		},
		{
			"doc_id": "doc4",
			"title": "Rust the programming language",
			"score": 0.0009629426822470302
		},
		{
			"doc_id": "doc1",
			"title": "Go the programming language",
			"score": 0.002248993082057349
		}
	],
	"query": "google programming language",
	"total_results": 4
}
```

### Get a document by ID

**GET** `/docs/{id}`
Retrieves a single document from the index by its ID.
**Response example:**
```json
{
	"id": "doc4",
	"title": "Rust the programming language",
	"content": "Rust is a general-purpose ...",
	"url": "https://en.wikipedia.org/wiki/Rust_(programming_language)",
	"length": 35566
}
```
### Get a documents counter

**GET** `/counter`
```json
{
	"docs_counter": 4
}
```

## 📂 Project Structure
```graphql
search-engine/
├── main.go            # Entry point
├── handlers/          # HTTP handlers (index, search)
├── services/          # Indexing, TF-IDF, scraping
├── models/            # Data models
├── database/          # BoltDB setup and helpers
└── utils/             # Utility functions
```

