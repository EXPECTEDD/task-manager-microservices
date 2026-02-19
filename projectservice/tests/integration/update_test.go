package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdate_SuccessUpdateName_Integration(t *testing.T) {
	email, pass := registrationUser(t)
	sessionId := loginUser(t, email, pass)

	projName := uniqueProjectName()

	projId := createProject(t, sessionId, projName)

	newProjName := uniqueProjectName()
	body := map[string]string{
		"new_name": newProjName,
	}

	resp := updateGetResponse(t, body, sessionId, projId)

	var respBody struct {
		IsUpdated bool `json:"is_updated"`
	}

	expBody := true
	expStatusCode := http.StatusOK

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expBody, respBody.IsUpdated)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestUpdate_NameAlreadyExists_Integration(t *testing.T) {
	email, pass := registrationUser(t)
	sessionId := loginUser(t, email, pass)

	projName := uniqueProjectName()

	_ = createProject(t, sessionId, projName)

	secondProjName := uniqueProjectName()

	secondProjId := createProject(t, sessionId, secondProjName)

	body := map[string]string{
		"new_name": projName,
	}

	resp := updateGetResponse(t, body, sessionId, secondProjId)

	var respBody struct {
		IsUpdated bool `json:"is_updated"`
	}

	expBody := false
	expStatusCode := http.StatusConflict

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expBody, respBody.IsUpdated)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestUpdate_ProjectNotFound_Integration(t *testing.T) {
	email, pass := registrationUser(t)
	sessionId := loginUser(t, email, pass)

	newProjName := uniqueProjectName()

	body := map[string]string{
		"new_name": newProjName,
	}

	resp := updateGetResponse(t, body, sessionId, 1)

	var respBody struct {
		IsUpdated bool `json:"is_updated"`
	}

	expBody := false
	expStatusCode := http.StatusNotFound

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expBody, respBody.IsUpdated)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestUpdate_MissingFieldName_Integration(t *testing.T) {
	email, pass := registrationUser(t)
	sessionId := loginUser(t, email, pass)

	projName := uniqueProjectName()

	projId := createProject(t, sessionId, projName)

	newProjName := uniqueProjectName()
	body := map[string]string{
		"new": newProjName,
	}

	resp := updateGetResponse(t, body, sessionId, projId)

	var respBody struct {
		IsUpdated bool `json:"is_updated"`
	}

	expBody := true
	expStatusCode := http.StatusOK

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expBody, respBody.IsUpdated)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func updateGetResponse(t *testing.T, body map[string]string, sessionId string, projId uint32) *http.Response {
	b, err := json.Marshal(body)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPatch, urlUpdate+"/"+strconv.Itoa(int(projId)), bytes.NewReader(b))
	require.NoError(t, err)

	jar, err := cookiejar.New(nil)
	require.NoError(t, err)

	cookies := []*http.Cookie{}
	cookie := &http.Cookie{
		Name:  "sessionId",
		Value: sessionId,
	}
	cookies = append(cookies, cookie)

	u, err := url.Parse(urlUpdate)
	require.NoError(t, err)

	jar.SetCookies(u, cookies)

	client := http.Client{
		Jar: jar,
	}

	resp, err := client.Do(req)
	require.NoError(t, err)

	return resp
}
