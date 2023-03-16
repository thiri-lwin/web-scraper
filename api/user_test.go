package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thiri-lwin/web_scraper/util"
)

func TestCreateUser(t *testing.T) {
	createUserTest(t)
}

func TestLoginUser(t *testing.T) {
	user := createUserTest(t)

	email := user.Get("email")
	password := user.Get("password")

	data := url.Values{}
	data.Add("email", email)
	data.Add("password", password)

	invalidEmail := url.Values{}
	invalidEmail.Add("email", util.RandomEmail())
	invalidEmail.Add("password", password)

	incorrectPassword := url.Values{}
	incorrectPassword.Add("email", email)
	incorrectPassword.Add("password", util.RandomString(6))
	testCases := []struct {
		name          string
		body          string
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: data.Encode(),
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusMovedPermanently, recorder.Code)
			},
		},
		{
			name: "IncompleteData",
			body: "",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				p, err := ioutil.ReadAll(recorder.Body)
				pageOK := err == nil && strings.Index(string(p), "<title>Sign In</title>") > 0 && strings.Index(string(p), "Something went wrong. Please try again later.") > 0
				require.Equal(t, true, pageOK)
			},
		},
		{
			name: "InvalidEmail",
			body: invalidEmail.Encode(),
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				p, err := ioutil.ReadAll(recorder.Body)
				pageOK := err == nil && strings.Index(string(p), "<title>Sign In</title>") > 0 && strings.Index(string(p), "Incorrect email or password.") > 0
				require.Equal(t, true, pageOK)
			},
		},
		{
			name: "IncorrectPassword",
			body: incorrectPassword.Encode(),
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)

				p, err := ioutil.ReadAll(recorder.Body)
				pageOK := err == nil && strings.Index(string(p), "<title>Sign In</title>") > 0 && strings.Index(string(p), "Incorrect email or password.") > 0
				require.Equal(t, true, pageOK)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(testDB)
			recorder := httptest.NewRecorder()

			req, _ := http.NewRequest("POST", "/login", strings.NewReader(tc.body))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}

}

func createUserTest(t *testing.T) url.Values {
	data := getCreateUserPostPayload()

	testCases := []struct {
		name          string
		body          string
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: data.Encode(),
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusFound, recorder.Code)
			},
		},
		{
			name: "DuplicateEmail",
			body: data.Encode(),
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

				p, err := ioutil.ReadAll(recorder.Body)
				pageOK := err == nil && strings.Index(string(p), "<title>Sign Up</title>") > 0 && strings.Index(string(p), "Email is already registered.") > 0
				require.Equal(t, true, pageOK)
			},
		},
		{
			name: "IncompleteData",
			body: "",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				p, err := ioutil.ReadAll(recorder.Body)
				pageOK := err == nil && strings.Index(string(p), "<title>Sign Up</title>") > 0 && strings.Index(string(p), "Something went wrong.") > 0
				require.Equal(t, true, pageOK)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(testDB)
			recorder := httptest.NewRecorder()

			req, _ := http.NewRequest("POST", "/signup", strings.NewReader(tc.body))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
	return data
}

func getCreateUserPostPayload() url.Values {
	email := util.RandomEmail()
	firstName := util.RandomString(6)
	lastName := util.RandomString(6)
	password := util.RandomString(6)

	data := url.Values{}
	data.Add("email", email)
	data.Add("first_name", firstName)
	data.Add("last_name", lastName)
	data.Add("password", password)
	return data
}
