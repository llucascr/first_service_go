package model

type Notebook struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type NotebookRequestDTO struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}