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

func TestGetAllTasks_Success_Integration(t *testing.T) {
	_, email, password := registrationUser(t)
	sessionId := loginUser(t, email, password)
	projId := createProject(t, sessionId, "NewProj")
	taskIdOne := createTask(t, sessionId, projId, "Desc", "")
	newTime := time.Now().Add(5 * time.Hour).UTC().Truncate(time.Second)
	taskIdSecond := createTask(t, sessionId, projId, "dsa", newTime.Format(time.RFC3339))

	resp := getAllTasksGetResponse(t, projId, sessionId)

	var respBody struct {
		Tasks []*taskdomain.TaskDomain `json:"tasks"`
	}

	expTasks := []*taskdomain.TaskDomain{
		{Id: taskIdOne, ProjectId: projId, Description: "Desc"},
		{Id: taskIdSecond, ProjectId: projId, Description: "dsa", Deadline: newTime},
	}
	expStatusCode := http.StatusOK

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	for i := range expTasks {
		require.Equal(t, expTasks[i], respBody.Tasks[i])
	}
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestGetAllTasks_InvalidProjectId_Integration(t *testing.T) {
	_, email, password := registrationUser(t)
	sessionId := loginUser(t, email, password)
	projId := createProject(t, sessionId, "NewProj")
	_ = createTask(t, sessionId, projId, "Desc", "")

	resp := getAllTasksGetResponse(t, 0, sessionId)

	var respBody struct {
		Tasks []*taskdomain.TaskDomain `json:"tasks"`
	}

	var expTasks []*taskdomain.TaskDomain = nil
	expStatusCode := http.StatusBadRequest

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expTasks, respBody.Tasks)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestGetAllTasks_TasksNotFound_Integration(t *testing.T) {
	_, email, password := registrationUser(t)
	sessionId := loginUser(t, email, password)
	projId := createProject(t, sessionId, "NewProj")

	resp := getAllTasksGetResponse(t, projId, sessionId)

	var respBody struct {
		Tasks []*taskdomain.TaskDomain `json:"tasks"`
	}

	var expTasks []*taskdomain.TaskDomain = nil
	expStatusCode := http.StatusNotFound

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expTasks, respBody.Tasks)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func getAllTasksGetResponse(t *testing.T, projId uint32, sessionId string) *http.Response {
	cookies := []*http.Cookie{}
	cookie := http.Cookie{
		Name:  "sessionId",
		Value: sessionId,
	}
	cookies = append(cookies, &cookie)

	u, err := url.Parse(urlGetAllTasks)

	cookieJar, err := cookiejar.New(nil)
	require.NoError(t, err)

	cookieJar.SetCookies(u, cookies)

	client := http.Client{
		Jar: cookieJar,
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%d", urlGetAllTasks, projId), nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)

	return resp
}
