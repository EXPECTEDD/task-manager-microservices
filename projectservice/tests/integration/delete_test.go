package integration

import (
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDelete_Success_Integration(t *testing.T) {
	email, pass := registrationUser(t)
	sessionId := loginUser(t, email, pass)

	projName := uniqueProjectName()

	projId := createProject(t, sessionId, projName)

	resp := deleteGetResponse(t, sessionId, projId)
	defer resp.Body.Close()

	var respBody struct {
		IsDeleted bool `json:"is_deleted"`
	}

	expBody := true
	expStatusCode := http.StatusOK

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expBody, respBody.IsDeleted)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestDelete_NotFound_Integration(t *testing.T) {
	email, pass := registrationUser(t)
	sessionId := loginUser(t, email, pass)

	projectId := 1

	resp := deleteGetResponse(t, sessionId, uint32(projectId))
	defer resp.Body.Close()

	var respBody struct {
		IsDeleted bool `json:"is_deleted"`
	}

	expBody := false
	expStatusCode := http.StatusNotFound

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expBody, respBody.IsDeleted)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestDelete_InvalidProjectId_Integration(t *testing.T) {
	email, pass := registrationUser(t)
	sessionId := loginUser(t, email, pass)

	projectId := 0

	resp := deleteGetResponse(t, sessionId, uint32(projectId))
	defer resp.Body.Close()

	var respBody struct {
		IsDeleted bool `json:"is_deleted"`
	}

	expBody := false
	expStatusCode := http.StatusBadRequest

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expBody, respBody.IsDeleted)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func deleteGetResponse(t *testing.T, sessionId string, projectId uint32) *http.Response {
	jar, err := cookiejar.New(nil)
	require.NoError(t, err)

	cookies := []*http.Cookie{}
	cookie := &http.Cookie{
		Name:  "sessionId",
		Value: sessionId,
	}
	cookies = append(cookies, cookie)

	u, err := url.Parse(urlDelete)
	require.NoError(t, err)

	jar.SetCookies(u, cookies)

	client := http.Client{
		Jar: jar,
	}

	req, err := http.NewRequest(http.MethodDelete, urlDelete+"/"+strconv.Itoa(int(projectId)), nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)

	return resp
}
