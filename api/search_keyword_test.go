package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thiri-lwin/web_scraper/util"
)

func TestGetKeywords(t *testing.T) {
	user := createUserTest(t)
	testCases := []struct {
		name          string
		data          string
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			data: user.Get("email"),
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				p, err := ioutil.ReadAll(recorder.Body)
				pageOK := err == nil && strings.Index(string(p), "<title>Keyword List</title>") > 0
				require.Equal(t, true, pageOK)
			},
		},
		{
			name: "UnauthorizedError",
			data: "",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusMovedPermanently, recorder.Code)
			},
		},
		{
			name: "InvalidUser",
			data: util.RandomEmail(),
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)

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

			if tc.data != "" {
				http.SetCookie(recorder, &http.Cookie{Name: util.Userkey, Value: tc.data})
			} else {
				http.SetCookie(recorder, &http.Cookie{Name: util.Userkey, Value: tc.data, MaxAge: -1})
			}

			req, _ := http.NewRequest("GET", "/keywords", nil)
			req.Header = http.Header{"Cookie": recorder.HeaderMap["Set-Cookie"]}

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}
