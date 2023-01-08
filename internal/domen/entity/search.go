package entity

import (
	pb "github.com/i-b8o/read-only_contracts/pb/searcher/v1"
)

type Search struct {
	DocID     *uint64 `json:"doc_id"`
	DocName   *string `json:"doc_name"`
	CID       *uint64 `json:"c_id"`
	CName     *string `json:"c_name"`
	UpdatedAt *string `json:"updated_at"`
	PID       *uint64 `json:"p_id"`
	Text      *string `json:"text"`
	Count     *uint32 `json:"count"`
}

func (s Search) ToResponse() (resp pb.SearchResponse) {
	if s.DocID != nil && *s.DocID > 0 {
		resp.DocID = *s.DocID
	}
	if s.DocName != nil && len(*s.DocName) > 0 {
		resp.DocName = *s.DocName
	}
	if s.CID != nil && *s.CID > 0 {
		resp.ChapterID = *s.CID
	}
	if s.CName != nil && len(*s.CName) > 0 {
		resp.ChapterName = *s.CName
	}
	if s.UpdatedAt != nil && len(*s.UpdatedAt) > 0 {
		resp.UpdatedAt = *s.UpdatedAt
	}
	if s.PID != nil && *s.PID > 0 {
		resp.ParagraphID = *s.PID
	}
	if s.Text != nil && len(*s.Text) > 0 {
		resp.Text = *s.Text
	}
	if s.Count != nil && *s.Count > 0 {
		resp.Count = *s.Count
	}
	return resp
}
