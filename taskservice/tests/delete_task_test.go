package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeleteTask_Success_Integration(t *testing.T) {
	_, email, password := registrationUser(t)
	sessionId := loginUser(t, email, password)
	projId := createProject(t, sessionId, "NewProj")
	taskId := createTask(t, sessionId, projId, "Desc", "")

	resp := deleteGetResponse(t, projId, taskId, sessionId)

	var respBody struct {
		Deleted bool `json:"deleted"`
	}

	expOut := true
	expStatusCode := http.StatusOK

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expOut, respBody.Deleted)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestDeleteTask_InvalidProjectId_Integration(t *testing.T) {
	_, email, password := registrationUser(t)
	sessionId := loginUser(t, email, password)
	projId := createProject(t, sessionId, "NewProj")
	taskId := createTask(t, sessionId, projId, "Desc", "")

	resp := deleteGetResponse(t, 0, taskId, sessionId)

	var respBody struct {
		Deleted bool `json:"deleted"`
	}

	expOut := false
	expStatusCode := http.StatusBadRequest

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expOut, respBody.Deleted)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestDeleteTask_InvalidTaskId_Integration(t *testing.T) {
	_, email, password := registrationUser(t)
	sessionId := loginUser(t, email, password)
	projId := createProject(t, sessionId, "NewProj")
	_ = createTask(t, sessionId, projId, "Desc", "")

	resp := deleteGetResponse(t, projId, 0, sessionId)

	var respBody struct {
		Deleted bool `json:"deleted"`
	}

	expOut := false
	expStatusCode := http.StatusBadRequest

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expOut, respBody.Deleted)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestDeleteTask_TaskNotFound_Integration(t *testing.T) {
	_, email, password := registrationUser(t)
	sessionId := loginUser(t, email, password)
	projId := createProject(t, sessionId, "NewProj")
	taskId := createTask(t, sessionId, projId, "Desc", "")

	resp := deleteGetResponse(t, projId, taskId+1, sessionId)

	var respBody struct {
		Deleted bool `json:"deleted"`
	}

	expOut := false
	expStatusCode := http.StatusNotFound

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expOut, respBody.Deleted)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func deleteGetResponse(t *testing.T, projId, taskId uint32, sessionId string) *http.Response {
	cookies := []*http.Cookie{}
	cookie := http.Cookie{
		Name:  "sessionId",
		Value: sessionId,
	}
	cookies = append(cookies, &cookie)

	u, err := url.Parse(urlDeleteTask)

	cookieJar, err := cookiejar.New(nil)
	require.NoError(t, err)

	cookieJar.SetCookies(u, cookies)

	client := http.Client{
		Jar: cookieJar,
	}

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%d/%d", urlDeleteTask, taskId, projId), nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)

	return resp
}
