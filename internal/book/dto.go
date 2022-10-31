package book

type CreateBookDto struct {
	Name     string `json:"name"`
	AuthorID int    `json:"author_id"`
}
