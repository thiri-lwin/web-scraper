package db

import (
	"context"
	"time"
)

type SearchResult struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Keyword   string    `db:"keyword"`
	Results   string    `db:"results"`
	CreatedAt time.Time `db:"created_at"`
}

func (s *Store) InsertSearchResult(ctx context.Context, arg SearchResult) (SearchResult, error) {
	data := SearchResult{}
	stmt, err := s.db.PrepareNamed(`
	INSERT INTO search_results (
		user_id,
		keyword,
		results)
	VALUES (
		:user_id, 
		:keyword, 
		:results
	) 
	RETURNING *`)
	if err != nil {
		return data, err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			return
		}
	}()
	return data, stmt.Get(&data, arg)
}

func (s *Store) GetSearchResultsByUserID(ctx context.Context, userID int) ([]SearchResult, error) {
	data := make([]SearchResult, 0)
	if err := s.db.Select(&data, `
	SELECT id, keyword, created_at from search_results
	WHERE user_id = $1`, userID); err != nil {
		return []SearchResult{}, err
	}
	return data, nil
}

func (s *Store) GetSearchResultByIDAndUserID(ctx context.Context, id, userID int) (SearchResult, error) {
	data := make([]SearchResult, 0)
	if err := s.db.Select(&data, `
	SELECT * from search_results
	WHERE id = $1 and user_id = $2 LIMIT 1`, id, userID); err != nil {
		return SearchResult{}, err
	}
	if len(data) > 0 {
		return data[0], nil
	}
	return SearchResult{}, nil
}
