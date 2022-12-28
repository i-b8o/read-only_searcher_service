package postgressql

import (
	"context"
	"fmt"
	client "read-only_search/pkg/client/postgresql"

	pb "github.com/i-b8o/read-only_contracts/pb/searcher/v1"
)

type chapterStorage struct {
	client client.PostgreSQLClient
}

func NewChapterStorage(client client.PostgreSQLClient) *chapterStorage {
	return &chapterStorage{client: client}
}

func (ss *chapterStorage) Chapters(ctx context.Context, searchQuery string) ([]pb.SearchResponse, error) {
	sql := `SELECT c.id, c.name, r.name, c.updated_at, count(*) OVER() AS full_count from chapter AS c INNER JOIN doc as r ON c.r_id = r.id WHERE c.ts @@ phraseto_tsquery('russian',$1)`
	var searchResults []pb.SearchResponse
	rows, err := ss.client.Query(ctx, sql, searchQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		search := pb.SearchResponse{}
		if err = rows.Scan(
			&search.ChapterID, &search.Text, &search.DocName, &search.UpdatedAt, &search.Count,
		); err != nil {
			return nil, err
		}

		searchResults = append(searchResults, search)
	}

	return searchResults, nil
}

func (ss *chapterStorage) ChaptersWithOffset(ctx context.Context, searchQuery string, params ...string) ([]pb.SearchResponse, error) {
	sql := `SELECT c.id, c.name, r.name, c.updated_at, count(*) OVER() AS full_count from chapter AS c INNER JOIN doc as r ON c.r_id = r.id WHERE c.ts @@ phraseto_tsquery('russian',$1)`
	// Pagination
	if len(params) == 2 {
		// sql += fmt.Sprintf(` AND (c.updated_at, c.id) > ('%s' :: TIMESTAMPTZ, '%s') ORDER BY c.updated_at, c.id LIMIT %s`, params[0], params[1], params[2])
		sql += fmt.Sprintf(` OFFSET %s LIMIT %s`, params[0], params[1])
	}
	// else if len(params) == 1 { // First page
	// 	sql += fmt.Sprintf(` LIMIT %s`, params[0])
	// }
	var searchResults []pb.SearchResponse
	rows, err := ss.client.Query(ctx, sql, searchQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		search := pb.SearchResponse{}
		if err = rows.Scan(
			&search.ChapterID, &search.Text, &search.DocName, &search.UpdatedAt, &search.Count,
		); err != nil {
			return nil, err
		}

		searchResults = append(searchResults, search)
	}

	return searchResults, nil
}
