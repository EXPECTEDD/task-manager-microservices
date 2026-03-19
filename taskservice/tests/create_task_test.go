package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateTask_Success_Without_Deadline_Integration(t *testing.T) {
	_, email, password := registrationUser(t)
	sessionId := loginUser(t, email, password)
	projId := createProject(t, sessionId, "NewProj")

	body := map[string]string{
		"description": "some description",
	}

	resp := createGetResponse(t, body, sessionId, projId)

	var respBody struct {
		TaskId uint32 `json:"task_id"`
	}

	expStatusCode := http.StatusOK

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Greater(t, respBody.TaskId, uint32(0))
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestCreateTask_Success_With_Deadline_Integration(t *testing.T) {
	_, email, password := registrationUser(t)
	sessionId := loginUser(t, email, password)
	projId := createProject(t, sessionId, "NewProj")

	deadline := time.Now()

	body := map[string]string{
		"description": "some description",
		"deadline":    deadline.Format(time.RFC3339),
	}

	resp := createGetResponse(t, body, sessionId, projId)

	var respBody struct {
		TaskId uint32 `json:"task_id"`
	}

	expStatusCode := http.StatusOK

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Greater(t, respBody.TaskId, uint32(0))
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestCreateTask_ProjectNotFound_Integration(t *testing.T) {
	_, email, password := registrationUser(t)
	sessionId := loginUser(t, email, password)
	projId := createProject(t, sessionId, "NewProj")

	deadline := time.Now()

	body := map[string]string{
		"description": "some description",
		"deadline":    deadline.Format(time.RFC3339),
	}

	resp := createGetResponse(t, body, sessionId, projId+1)

	var respBody struct {
		TaskId uint32 `json:"task_id"`
	}

	expStatusCode := http.StatusNotFound

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, respBody.TaskId, uint32(0))
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestCreateTask_InvalidProjectId_Integration(t *testing.T) {
	_, email, password := registrationUser(t)
	sessionId := loginUser(t, email, password)
	_ = createProject(t, sessionId, "NewProj")

	deadline := time.Now()

	body := map[string]string{
		"description": "some description",
		"deadline":    deadline.Format(time.RFC3339),
	}

	resp := createGetResponse(t, body, sessionId, 0)

	var respBody struct {
		TaskId uint32 `json:"task_id"`
	}

	expStatusCode := http.StatusBadRequest

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, respBody.TaskId, uint32(0))
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestCreateTask_InvalidDescription_Integration(t *testing.T) {
	_, email, password := registrationUser(t)
	sessionId := loginUser(t, email, password)
	projId := createProject(t, sessionId, "NewProj")

	deadline := time.Now()

	body := map[string]string{
		"description": strings.Repeat("some description", 100),
		"deadline":    deadline.Format(time.RFC3339),
	}

	resp := createGetResponse(t, body, sessionId, projId)

	var respBody struct {
		TaskId uint32 `json:"task_id"`
	}

	expStatusCode := http.StatusBadRequest

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, respBody.TaskId, uint32(0))
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func createGetResponse(t *testing.T, body map[string]string, sessionId string, projId uint32) *http.Response {
	b, err := json.Marshal(body)
	require.NoError(t, err)

	cookies := []*http.Cookie{}
	cookie := http.Cookie{
		Name:  "sessionId",
		Value: sessionId,
	}
	cookies = append(cookies, &cookie)

	u, err := url.Parse(urlCreateTask)

	cookieJar, err := cookiejar.New(nil)
	require.NoError(t, err)

	cookieJar.SetCookies(u, cookies)

	client := http.Client{
		Jar: cookieJar,
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%d", urlCreateTask, projId), bytes.NewReader(b))
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)

	return resp
}
