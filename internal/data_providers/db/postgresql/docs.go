package postgressql

import (
	"context"
	"fmt"

	client "read-only_search/pkg/client/postgresql"

	pb "github.com/i-b8o/read-only_contracts/pb/searcher/v1"
)

type docStorage struct {
	client client.PostgreSQLClient
}

func NewDocStorage(client client.PostgreSQLClient) *docStorage {
	return &docStorage{client: client}
}

func (s *docStorage) Docs(ctx context.Context, searchQuery string) ([]*pb.SearchResponse, error) {
	sql := `SELECT id, name, title, updated_at, count(*) OVER() AS full_count from doc WHERE ts @@ phraseto_tsquery('russian',$1)`
	var searchResults []*pb.SearchResponse
	rows, err := s.client.Query(ctx, sql, searchQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		search := &pb.SearchResponse{}
		if err = rows.Scan(
			&search.DocID, &search.DocName, &search.Text, &search.UpdatedAt, &search.Count,
		); err != nil {
			return nil, err
		}

		searchResults = append(searchResults, search)
	}

	return searchResults, nil
}

func (s *docStorage) DocsWithOffset(ctx context.Context, searchQuery string, offset, limit uint32) ([]*pb.SearchResponse, error) {
	sql := `SELECT id, name, title, updated_at, count(*) OVER() AS full_count from doc WHERE ts @@ phraseto_tsquery('russian',$1)`
	// Pagination
	// sql += fmt.Sprintf(` AND (updated_at, id) > ('%s' :: TIMESTAMPTZ, '%s') ORDER BY updated_at, id LIMIT %s`, params[0], params[1], params[2])
	sql += fmt.Sprintf(` OFFSET %d LIMIT %d`, offset, limit)
	//  else if len(params) == 1 { // First page
	// 	sql += fmt.Sprintf(` LIMIT %s`, params[0])
	// }
	var searchResults []*pb.SearchResponse
	rows, err := s.client.Query(ctx, sql, searchQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		search := &pb.SearchResponse{}
		if err = rows.Scan(
			&search.DocID, &search.DocName, &search.Text, &search.UpdatedAt, &search.Count,
		); err != nil {
			return nil, err
		}

		searchResults = append(searchResults, search)
	}

	return searchResults, nil
}
