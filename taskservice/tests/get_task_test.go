package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	taskdomain "taskservice/internal/domain/task"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetTask_Success_Integration(t *testing.T) {
	_, email, password := registrationUser(t)
	sessionId := loginUser(t, email, password)
	projId := createProject(t, sessionId, "NewProj")
	newTime := time.Now().Add(5 * time.Hour).UTC().Truncate(time.Second)
	taskId := createTask(t, sessionId, projId, "dsa", newTime.Format(time.RFC3339))

	resp := getTaskGetResponse(t, taskId, projId, sessionId)

	var respBody struct {
		Task *taskdomain.TaskDomain `json:"task"`
	}

	expTask := &taskdomain.TaskDomain{Id: taskId, ProjectId: projId, Description: "dsa", Deadline: newTime}
	expStatusCode := http.StatusOK

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expTask, respBody.Task)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestGetTask_InvalidProjectId_Integration(t *testing.T) {
	_, email, password := registrationUser(t)
	sessionId := loginUser(t, email, password)
	projId := createProject(t, sessionId, "NewProj")
	taskId := createTask(t, sessionId, projId, "Desc", "")

	resp := getTaskGetResponse(t, taskId, 0, sessionId)

	var respBody struct {
		Task *taskdomain.TaskDomain `json:"task"`
	}

	var expTasks *taskdomain.TaskDomain = nil
	expStatusCode := http.StatusBadRequest

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expTasks, respBody.Task)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestGetTask_TasksNotFound_Integration(t *testing.T) {
	_, email, password := registrationUser(t)
	sessionId := loginUser(t, email, password)
	projId := createProject(t, sessionId, "NewProj")

	resp := getTaskGetResponse(t, 1, projId, sessionId)

	var respBody struct {
		Task *taskdomain.TaskDomain `json:"task"`
	}

	var expTasks *taskdomain.TaskDomain = nil
	expStatusCode := http.StatusNotFound

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expTasks, respBody.Task)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestGetTask_InvalidTaskId_Integration(t *testing.T) {
	_, email, password := registrationUser(t)
	sessionId := loginUser(t, email, password)
	projId := createProject(t, sessionId, "NewProj")
	_ = createTask(t, sessionId, projId, "Desc", "")

	resp := getTaskGetResponse(t, 0, projId, sessionId)

	var respBody struct {
		Task *taskdomain.TaskDomain `json:"task"`
	}

	var expTasks *taskdomain.TaskDomain = nil
	expStatusCode := http.StatusBadRequest

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expTasks, respBody.Task)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func getTaskGetResponse(t *testing.T, taskId uint32, projId uint32, sessionId string) *http.Response {
	cookies := []*http.Cookie{}
	cookie := http.Cookie{
		Name:  "sessionId",
		Value: sessionId,
	}
	cookies = append(cookies, &cookie)

	u, err := url.Parse(urlGetTask)

	cookieJar, err := cookiejar.New(nil)
	require.NoError(t, err)

	cookieJar.SetCookies(u, cookies)

	client := http.Client{
		Jar: cookieJar,
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%d/%d", urlGetTask, taskId, projId), nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)

	return resp
}
