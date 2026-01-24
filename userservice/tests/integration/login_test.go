package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogin_Success_Integration(t *testing.T) {
	email, pass := registrateUser(t)

	body := map[string]string{
		"email":    email,
		"password": pass,
	}

	b, err := json.Marshal(body)
	require.NoError(t, err)

	resp, err := http.Post(urlLog, contentType, bytes.NewReader(b))
	require.NoError(t, err)
	defer resp.Body.Close()

	expBody := map[string]string{
		"first_name":  "Ivan",
		"middle_name": "Ivanovich",
		"last_name":   "Ivanov",
	}
	expStatusCode := http.StatusOK

	var respBody struct {
		Userdata map[string]string `json:"user"`
	}

	cookies := resp.Cookies()
	var sessionCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "sessionId" {
			sessionCookie = c
		}
	}

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expBody, respBody.Userdata)
	require.Equal(t, expStatusCode, resp.StatusCode)
	require.NotNil(t, sessionCookie)
	require.NotEmpty(t, sessionCookie.Value)
}

func TestLogin_BadRequest_Integration(t *testing.T) {
	email, pass := registrateUser(t)

	body := map[string]string{
		"EEEEMMMMAAAIIILLLL": email,
		"password":           pass,
	}

	b, err := json.Marshal(body)
	require.NoError(t, err)

	resp, err := http.Post(urlLog, contentType, bytes.NewReader(b))
	require.NoError(t, err)
	defer resp.Body.Close()

	expBody := map[string]string{
		"Email": "field is required",
	}
	expStatusCode := http.StatusBadRequest

	var resBody struct {
		Errors map[string]string `json:"errors"`
	}

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&resBody))
	require.Equal(t, expBody, resBody.Errors)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestLogin_UserNotFound_Integration(t *testing.T) {
	body := map[string]string{
		"email":    uniqueEmail(),
		"password": "somePass",
	}

	b, err := json.Marshal(body)
	require.NoError(t, err)

	resp, err := http.Post(urlLog, contentType, bytes.NewReader(b))
	require.NoError(t, err)
	defer resp.Body.Close()

	expBody := "user not found"
	expStatusCode := http.StatusNotFound

	var resBody struct {
		Error string `json:"error"`
	}

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&resBody))
	require.Equal(t, expBody, resBody.Error)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestLogin_WrongPassword_Integration(t *testing.T) {
	email, _ := registrateUser(t)
	body := map[string]string{
		"email":    email,
		"password": "wrong password",
	}

	b, err := json.Marshal(body)
	require.NoError(t, err)

	resp, err := http.Post(urlLog, contentType, bytes.NewReader(b))
	require.NoError(t, err)
	defer resp.Body.Close()

	expBody := "wrong password"
	expStatusCode := http.StatusUnauthorized

	var resBody struct {
		Error string `json:"error"`
	}

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&resBody))
	require.Equal(t, expBody, resBody.Error)
	require.Equal(t, expStatusCode, resp.StatusCode)
}
