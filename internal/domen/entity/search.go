package entity

import "time"

type Search struct {
	DocID       *uint64    `json:"doc_id"`
	DocName     *string    `json:"doc_name"`
	ChapterID   *uint64    `json:"chapter_id"`
	ChapterName *string    `json:"chapter_name"`
	UpdatedAt   *time.Time `json:"updated_at"`
	ParagraphID *uint64    `json:"paragraph_id"`
	Text        *string    `json:"text"`
	Count       *uint64    `json:"count"`
}
