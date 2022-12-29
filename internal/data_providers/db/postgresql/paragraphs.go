package postgressql

import (
	"context"
	"fmt"
	client "read-only_search/pkg/client/postgresql"

	pb "github.com/i-b8o/read-only_contracts/pb/searcher/v1"
)

type paragraphStorage struct {
	client client.PostgreSQLClient
}

func NewParagraphStorage(client client.PostgreSQLClient) *paragraphStorage {
	return &paragraphStorage{client: client}
}

func (ss *paragraphStorage) Paragraphs(ctx context.Context, searchQuery string) ([]*pb.SearchResponse, error) {
	sql := `SELECT p.id, p.content, c.name, d.name, c.updated_at, count(*) OVER() AS full_count from paragraph AS p INNER JOIN chapter as c ON c.id = p.c_id INNER JOIN doc AS d ON c.d_id = d.id  WHERE p.ts @@ phraseto_tsquery('russian',$1)`
	var searchResults []*pb.SearchResponse
	rows, err := ss.client.Query(ctx, sql, searchQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		search := &pb.SearchResponse{}
		if err = rows.Scan(
			&search.ParagraphID, &search.Text, &search.ChapterName, &search.DocID, &search.UpdatedAt, &search.Count,
		); err != nil {
			return nil, err
		}

		searchResults = append(searchResults, search)
	}

	return searchResults, nil
}

func (ss *paragraphStorage) ParagraphsWithOffset(ctx context.Context, searchQuery string, offset, limit uint32) ([]*pb.SearchResponse, error) {
	sql := `SELECT p.id, p.content, c.name, d.name, c.updated_at, count(*) OVER() AS full_count from paragraph AS p INNER JOIN chapter as c ON c.id = p.c_id INNER JOIN doc AS d ON c.d_id = d.id  WHERE p.ts @@ phraseto_tsquery('russian',$1)`
	// Pagination
	// sql += fmt.Sprintf(` AND (c.updated_at, p.id) > ('%s' :: TIMESTAMPTZ, '%s') ORDER BY c.updated_at, p.id LIMIT %s`, params[0], params[1], params[2])
	sql += fmt.Sprintf(` OFFSET %d LIMIT %d`, offset, limit)

	// else if len(params) == 1 { // First page
	// 	sql += fmt.Sprintf(` LIMIT %s`, params[0])
	// }
	var searchResults []*pb.SearchResponse
	rows, err := ss.client.Query(ctx, sql, searchQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		search := &pb.SearchResponse{}
		if err = rows.Scan(
			&search.ParagraphID, &search.Text, &search.ChapterName, &search.DocID, &search.UpdatedAt, &search.Count,
		); err != nil {
			return nil, err
		}

		searchResults = append(searchResults, search)
	}

	return searchResults, nil
}
