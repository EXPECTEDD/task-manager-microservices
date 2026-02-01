package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreate_Success_Integration(t *testing.T) {
	email, pass := registrationUser(t)
	sessionId := loginUser(t, email, pass)

	body := map[string]string{
		"name": uniqueProjectName(),
	}

	resp := createGetResponse(t, sessionId, body)
	defer resp.Body.Close()

	var respBody struct {
		IsCreated bool `json:"is_created"`
	}

	expBody := true
	expStatusCode := http.StatusOK

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expBody, respBody.IsCreated)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestCreate_MissingFieldName_Integration(t *testing.T) {
	email, pass := registrationUser(t)
	sessionId := loginUser(t, email, pass)

	body := map[string]string{
		"nam": uniqueProjectName(),
	}

	resp := createGetResponse(t, sessionId, body)
	defer resp.Body.Close()

	var respBody struct {
		IsCreated bool `json:"is_created"`
	}

	expBody := false
	expStatusCode := http.StatusBadRequest

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expBody, respBody.IsCreated)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestCreate_IvalidName_Integration(t *testing.T) {
	email, pass := registrationUser(t)
	sessionId := loginUser(t, email, pass)

	body := map[string]string{
		"name": strings.Repeat(uniqueProjectName(), 300),
	}

	resp := createGetResponse(t, sessionId, body)
	defer resp.Body.Close()

	var respBody struct {
		IsCreated bool `json:"is_created"`
	}

	expBody := false
	expStatusCode := http.StatusBadRequest

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expBody, respBody.IsCreated)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestCreate_AlreadyExists_Integration(t *testing.T) {
	email, pass := registrationUser(t)
	sessionId := loginUser(t, email, pass)

	body := map[string]string{
		"name": uniqueProjectName(),
	}

	resp := createGetResponse(t, sessionId, body)
	resp.Body.Close()
	resp = createGetResponse(t, sessionId, body)
	defer resp.Body.Close()

	var respBody struct {
		IsCreated bool `json:"is_created"`
	}

	expBody := false
	expStatusCode := http.StatusConflict

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expBody, respBody.IsCreated)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func createGetResponse(t *testing.T, sessionId string, body map[string]string) *http.Response {
	b, err := json.Marshal(body)
	require.NoError(t, err)

	jar, err := cookiejar.New(nil)
	require.NoError(t, err)

	var cookies []*http.Cookie
	cookie := &http.Cookie{
		Name:  "sessionId",
		Value: sessionId,
	}
	cookies = append(cookies, cookie)

	u, err := url.Parse(urlCreate)
	require.NoError(t, err)
	jar.SetCookies(u, cookies)

	client := http.Client{
		Jar: jar,
	}

	req, err := http.NewRequest(http.MethodPost, urlCreate, bytes.NewReader(b))
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	return resp
}
