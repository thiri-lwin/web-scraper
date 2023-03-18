package api

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thiri-lwin/web_scraper/db"
	"github.com/thiri-lwin/web_scraper/util"
)

var testDB *db.Store

func TestMain(m *testing.M) {
	// gin.SetMode(gin.TestMode)

	var err error
	testDB, err = db.NewStore(dbDriver, dbstring)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	os.Exit(m.Run())
}

func newTestServer(store *db.Store) *Server {
	templatePath := "../templates/*.html"
	assetsPath := "../assets"
	cssPath := "../templates/css"
	keywordChan := make(chan util.UploadedFile)
	server := NewServer(store, keywordChan, templatePath, assetsPath, cssPath)
	return server
}

const (
	dbDriver = "postgres"
	dbstring = "postgresql://postgres:postgres@localhost:5432/web_scraper?sslmode=disable" // move to config later
)

func TestInitPage(t *testing.T) {
	testCases := []struct {
		name          string
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				p, err := ioutil.ReadAll(recorder.Body)
				pageOK := err == nil && strings.Index(string(p), "<title>Sign In</title>") > 0
				require.Equal(t, true, pageOK)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(testDB)
			recorder := httptest.NewRecorder()

			req, _ := http.NewRequest("GET", "/", nil)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestSignUpGetHandler(t *testing.T) {
	testCases := []struct {
		name          string
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				p, err := ioutil.ReadAll(recorder.Body)
				pageOK := err == nil && strings.Index(string(p), "<title>Sign Up</title>") > 0
				require.Equal(t, true, pageOK)
			},
		},
		{
			name: "WrongTitle",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				p, err := ioutil.ReadAll(recorder.Body)
				require.NoError(t, err)
				pageOK := strings.Index(string(p), "<title>Wrong Title</title>") > 0
				require.Equal(t, false, pageOK)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(testDB)
			recorder := httptest.NewRecorder()

			req, _ := http.NewRequest("GET", "/signup", nil)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}
