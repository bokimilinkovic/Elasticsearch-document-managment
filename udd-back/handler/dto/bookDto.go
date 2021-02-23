package dto

type BookDto struct {
	ISBN        string `json:"isbn"`
	PublishYear string `json:"publish_year"`
	PageNumber  int    `json:"page_number"`
	///Keywords    []string  `json:"keywords"`
	Genre   string `json:"genre"`
	Author  string `json:"author"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
