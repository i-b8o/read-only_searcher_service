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

func (ss *paragraphStorage) ParagraphSearch(ctx context.Context, searchQuery string, params ...string) ([]pb.SearchResponse, error) {
	sql := `SELECT p.id, p.content, c.name, r.name, c.updated_at, count(*) OVER() AS full_count from paragraph AS p INNER JOIN chapter as c ON c.id = p.c_id INNER JOIN doc AS r ON c.r_id = r.id  WHERE p.ts @@ phraseto_tsquery('russian',$1)`
	// Pagination
	if len(params) == 3 {
		// sql += fmt.Sprintf(` AND (c.updated_at, p.id) > ('%s' :: TIMESTAMPTZ, '%s') ORDER BY c.updated_at, p.id LIMIT %s`, params[0], params[1], params[2])
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
			&search.ParagraphID, &search.Text, &search.ChapterName, &search.DocID, &search.UpdatedAt, &search.Count,
		); err != nil {
			return nil, err
		}

		searchResults = append(searchResults, search)
	}

	return searchResults, nil
}
