package api

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thiri-lwin/web_scraper/util"
)

func TestUploadKeywords(t *testing.T) {

	filename := "test/sample_success.csv"
	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("Failed to open file %s: %s", filename, err)
	}
	defer file.Close()

	// Create multipart form data with CSV file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("csvfile", filepath.Base(filename))
	if err != nil {
		t.Fatalf("Failed to create form file: %s", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatalf("Failed to copy file to form: %s", err)
	}
	writer.Close()
	user := createUserTest(t)
	testCases := []struct {
		name          string
		data          string
		body          *bytes.Buffer
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			data: user.Get("email"),
			body: body,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "UnauthorizedError",
			data: "",
			body: body,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusMovedPermanently, recorder.Code)
			},
		},
		{
			name: "InvalidUser",
			data: util.RandomEmail(),
			body: body,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)

				p, err := ioutil.ReadAll(recorder.Body)
				pageOK := err == nil && strings.Index(string(p), "<title>Sign In</title>") > 0
				require.Equal(t, true, pageOK)
			},
		},
		{
			name: "InvalidFile",
			data: user.Get("email"),
			body: &bytes.Buffer{},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

				p, err := ioutil.ReadAll(recorder.Body)
				pageOK := err == nil && strings.Index(string(p), "<title>Upload</title>") > 0
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

			req, _ := http.NewRequest("POST", "/upload", tc.body)
			req.Header = http.Header{"Cookie": recorder.HeaderMap["Set-Cookie"]}
			req.Header.Add("Content-Type", writer.FormDataContentType())

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestUploadGetHandler(t *testing.T) {
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
				pageOK := err == nil && strings.Index(string(p), "<title>Upload</title>") > 0
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

			req, _ := http.NewRequest("GET", "/upload", nil)
			req.Header = http.Header{"Cookie": recorder.HeaderMap["Set-Cookie"]}

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

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
