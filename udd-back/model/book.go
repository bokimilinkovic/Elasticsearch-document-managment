package model

import (
	"time"
)

type Book struct {
	ISBN        string    `json:"isbn"`
	PublishYear time.Time `json:"publish_year"`
	PageNumber  int       `json:"page_number"`
	///Keywords    []string  `json:"keywords"`
	Genre   string `json:"genre"`
	Author  string `json:"author"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
