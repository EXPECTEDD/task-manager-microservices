package integration

import (
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type bodyFields struct {
	Id        uint32    `json:"Id"`
	OwnerId   uint32    `json:"OwnerId"`
	Name      string    `json:"Name"`
	CreatedAt time.Time `json:"CreatedAt"`
}

func TestGetAll_Success_Integration(t *testing.T) {
	_, email, pass := registrationUser(t)
	sessionId := loginUser(t, email, pass)

	projName := uniqueProjectName()

	createProject(t, sessionId, projName)

	resp := getAllGetResponse(t, sessionId)
	defer resp.Body.Close()

	var respBody struct {
		Projects []bodyFields `json:"projects"`
	}

	expStatusCode := http.StatusOK

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Greater(t, int(respBody.Projects[0].Id), 0)
	require.Greater(t, int(respBody.Projects[0].OwnerId), 0)
	require.Equal(t, projName, respBody.Projects[0].Name)
	require.NotNil(t, respBody.Projects[0].CreatedAt)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestGetAll_NotFound_Integration(t *testing.T) {
	_, email, pass := registrationUser(t)
	sessionId := loginUser(t, email, pass)

	resp := getAllGetResponse(t, sessionId)
	defer resp.Body.Close()

	expStatusCode := http.StatusNotFound

	require.Equal(t, expStatusCode, resp.StatusCode)
}

func getAllGetResponse(t *testing.T, sessionId string) *http.Response {
	jar, err := cookiejar.New(nil)
	require.NoError(t, err)

	cookies := []*http.Cookie{}
	cookie := &http.Cookie{
		Name:  "sessionId",
		Value: sessionId,
	}
	cookies = append(cookies, cookie)

	u, err := url.Parse(urlGetAll)
	require.NoError(t, err)

	jar.SetCookies(u, cookies)

	client := http.Client{
		Jar: jar,
	}

	req, err := http.NewRequest(http.MethodGet, urlGetAll, nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	return resp
}
