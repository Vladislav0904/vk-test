package items

import "time"

type Item struct {
	ID        uint64    `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`
	Text      string    `db:"text" json:"text"`
	ImageURL  string    `db:"image_url" json:"image_url"`
	Price     float32   `db:"price" json:"price"`
	AuthorID  uint64    `db:"author_id" json:"author_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
