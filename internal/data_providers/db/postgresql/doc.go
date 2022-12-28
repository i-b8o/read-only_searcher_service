package postgressql

import (
	"context"
	"fmt"

	client "read-only_search/pkg/client/postgresql"

	"read-only_search/internal/domen/entity"
)

type docStorage struct {
	client client.PostgreSQLClient
}

func NewDocStorage(client client.PostgreSQLClient) *docStorage {
	return &docStorage{client: client}
}

func (ss *docStorage) DocSearch(ctx context.Context, searchQuery string, params ...string) ([]entity.Search, error) {
	sql := `SELECT id, name, title, updated_at, count(*) OVER() AS full_count from doc WHERE ts @@ phraseto_tsquery('russian',$1)`
	// Pagination
	if len(params) == 2 {
		// sql += fmt.Sprintf(` AND (updated_at, id) > ('%s' :: TIMESTAMPTZ, '%s') ORDER BY updated_at, id LIMIT %s`, params[0], params[1], params[2])
		sql += fmt.Sprintf(` OFFSET %s LIMIT %s`, params[0], params[1])
	}
	//  else if len(params) == 1 { // First page
	// 	sql += fmt.Sprintf(` LIMIT %s`, params[0])
	// }
	var searchResults []entity.Search
	rows, err := ss.client.Query(ctx, sql, searchQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		search := entity.Search{}
		if err = rows.Scan(
			&search.DocID, &search.DocName, &search.Text, &search.UpdatedAt, &search.Count,
		); err != nil {
			return nil, err
		}

		searchResults = append(searchResults, search)
	}

	return searchResults, nil
}
