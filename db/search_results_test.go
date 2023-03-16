package db

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/thiri-lwin/web_scraper/util"
)

type resultDB struct {
	HTMLCode           string
	NumAds             int32
	NumLinks           int32
	TotalSearchResults string
}

func createRandomSearchResult(t *testing.T) SearchResult {
	user := createRandomUser(t)

	dbRes := resultDB{
		HTMLCode:           util.RandomString(50),
		NumAds:             int32(util.RandomInt(0, 30)),
		NumLinks:           int32(util.RandomInt(0, 30)),
		TotalSearchResults: util.RandomString(10),
	}
	data, _ := json.Marshal(dbRes)
	arg := SearchResult{
		UserID:  user.ID,
		Keyword: util.RandomString(6),
		Results: string(data),
	}

	searchRes, err := testDB.InsertSearchResult(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, searchRes)

	require.NotZero(t, searchRes.ID)
	require.Equal(t, arg.Keyword, searchRes.Keyword)
	require.Equal(t, arg.UserID, user.ID)
	require.NotZero(t, searchRes.Results)
	require.NotZero(t, searchRes.CreatedAt)

	return searchRes
}

func TestInsertSearchResult(t *testing.T) {
	createRandomSearchResult(t)
}

func TestGetSearchResultsByUserID(t *testing.T) {
	searchResult1 := createRandomSearchResult(t)
	searchResults, err := testDB.GetSearchResultsByUserID(context.Background(), searchResult1.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, searchResults)

	require.NotZero(t, len(searchResults))
}
