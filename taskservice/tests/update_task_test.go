package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUpdateTask_Success_Integration(t *testing.T) {
	_, email, password := registrationUser(t)
	sessionId := loginUser(t, email, password)
	projId := createProject(t, sessionId, "NewProj")
	taskId := createTask(t, sessionId, projId, "NewTask", time.Now().Format(time.RFC3339))

	body := map[string]string{
		"new_description": "new description",
		"new_deadline":    time.Now().Add(1 * time.Hour).Format(time.RFC3339),
	}

	resp := updateGetResponse(t, body, sessionId, projId, taskId)

	var respBody struct {
		Updated bool `json:"updated"`
	}

	expResult := true
	expStatusCode := http.StatusOK

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expResult, respBody.Updated)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestUpdateTask_SuccessWithoutDescription_Integration(t *testing.T) {
	_, email, password := registrationUser(t)
	sessionId := loginUser(t, email, password)
	projId := createProject(t, sessionId, "NewProj")
	taskId := createTask(t, sessionId, projId, "NewTask", time.Now().Format(time.RFC3339))

	body := map[string]string{
		"new_deadline": time.Now().Add(1 * time.Hour).Format(time.RFC3339),
	}

	resp := updateGetResponse(t, body, sessionId, projId, taskId)

	var respBody struct {
		Updated bool `json:"updated"`
	}

	expResult := true
	expStatusCode := http.StatusOK

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expResult, respBody.Updated)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestUpdateTask_SuccessWithoutDeadline_Integration(t *testing.T) {
	_, email, password := registrationUser(t)
	sessionId := loginUser(t, email, password)
	projId := createProject(t, sessionId, "NewProj")
	taskId := createTask(t, sessionId, projId, "NewTask", time.Now().Format(time.RFC3339))

	body := map[string]string{
		"new_description": "new description",
		"new_deadline":    time.Now().Add(1 * time.Hour).Format(time.RFC3339),
	}

	resp := updateGetResponse(t, body, sessionId, projId, taskId)

	var respBody struct {
		Updated bool `json:"updated"`
	}

	expResult := true
	expStatusCode := http.StatusOK

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expResult, respBody.Updated)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestUpdateTask_TaskNotFound_DifProjectId_Integration(t *testing.T) {
	_, email, password := registrationUser(t)
	sessionId := loginUser(t, email, password)
	projId := createProject(t, sessionId, "NewProj")
	taskId := createTask(t, sessionId, projId, "NewTask", time.Now().Format(time.RFC3339))

	body := map[string]string{
		"new_description": "new description",
		"new_deadline":    time.Now().Add(1 * time.Hour).Format(time.RFC3339),
	}

	resp := updateGetResponse(t, body, sessionId, projId+1, taskId)

	var respBody struct {
		Updated bool `json:"updated"`
	}

	expResult := false
	expStatusCode := http.StatusNotFound

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expResult, respBody.Updated)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestUpdateTask_TaskNotFound_DifTaskId_Integration(t *testing.T) {
	_, email, password := registrationUser(t)
	sessionId := loginUser(t, email, password)
	projId := createProject(t, sessionId, "NewProj")
	taskId := createTask(t, sessionId, projId, "NewTask", time.Now().Format(time.RFC3339))

	body := map[string]string{
		"new_description": "new description",
		"new_deadline":    time.Now().Add(1 * time.Hour).Format(time.RFC3339),
	}

	resp := updateGetResponse(t, body, sessionId, projId, taskId+1)

	var respBody struct {
		Updated bool `json:"updated"`
	}

	expResult := false
	expStatusCode := http.StatusNotFound

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expResult, respBody.Updated)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestUpdateTask_InvalidTaskId_Integration(t *testing.T) {
	_, email, password := registrationUser(t)
	sessionId := loginUser(t, email, password)
	projId := createProject(t, sessionId, "NewProj")
	_ = createTask(t, sessionId, projId, "NewTask", time.Now().Format(time.RFC3339))

	body := map[string]string{
		"new_description": "new description",
		"new_deadline":    time.Now().Add(1 * time.Hour).Format(time.RFC3339),
	}

	resp := updateGetResponse(t, body, sessionId, projId, 0)

	var respBody struct {
		Updated bool `json:"updated"`
	}

	expResult := false
	expStatusCode := http.StatusBadRequest

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expResult, respBody.Updated)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestUpdateTask_InvalidProjectId_Integration(t *testing.T) {
	_, email, password := registrationUser(t)
	sessionId := loginUser(t, email, password)
	projId := createProject(t, sessionId, "NewProj")
	taskId := createTask(t, sessionId, projId, "NewTask", time.Now().Format(time.RFC3339))

	body := map[string]string{
		"new_description": "new description",
		"new_deadline":    time.Now().Add(1 * time.Hour).Format(time.RFC3339),
	}

	resp := updateGetResponse(t, body, sessionId, 0, taskId)

	var respBody struct {
		Updated bool `json:"updated"`
	}

	expResult := false
	expStatusCode := http.StatusBadRequest

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expResult, respBody.Updated)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestUpdateTask_WithoutAll_Integration(t *testing.T) {
	_, email, password := registrationUser(t)
	sessionId := loginUser(t, email, password)
	projId := createProject(t, sessionId, "NewProj")
	taskId := createTask(t, sessionId, projId, "NewTask", time.Now().Format(time.RFC3339))

	body := map[string]string{}

	resp := updateGetResponse(t, body, sessionId, projId, taskId)

	var respBody struct {
		Updated bool `json:"updated"`
	}

	expResult := false
	expStatusCode := http.StatusBadRequest

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expResult, respBody.Updated)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func updateGetResponse(t *testing.T, body map[string]string, sessionId string, projId uint32, taskId uint32) *http.Response {
	b, err := json.Marshal(body)
	require.NoError(t, err)

	cookies := []*http.Cookie{}
	cookie := http.Cookie{
		Name:  "sessionId",
		Value: sessionId,
	}
	cookies = append(cookies, &cookie)

	u, err := url.Parse(urlUpdateTask)

	cookieJar, err := cookiejar.New(nil)
	require.NoError(t, err)

	cookieJar.SetCookies(u, cookies)

	client := http.Client{
		Jar: cookieJar,
	}

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/%d/%d", urlUpdateTask, taskId, projId), bytes.NewReader(b))
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)

	return resp
}
