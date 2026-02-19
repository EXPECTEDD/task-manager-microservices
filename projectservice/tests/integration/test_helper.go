package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var (
	urlRegistration = "http://localhost:44044/user/registration"
	urlLogin        = "http://localhost:44044/user/login"
	urlCreate       = "http://localhost:44046/project/create"
	urlDelete       = "http://localhost:44046/project/delete"
	urlGetAll       = "http://localhost:44046/project/getall"
	urlUpdate       = "http://localhost:44046/project/update"
)

const (
	contentType = "application/json"
)

func registrationUser(t *testing.T) (string, string) {
	email := uniqueEmail()
	password := "somePass"

	body := map[string]string{
		"first_name":  "Ivan",
		"middle_name": "Ivanovich",
		"last_name":   "Ivanov",
		"password":    password,
		"email":       email,
	}

	b, err := json.Marshal(body)
	require.NoError(t, err)

	resp, err := http.Post(urlRegistration, contentType, bytes.NewReader(b))
	require.NoError(t, err)
	defer resp.Body.Close()

	var respBody struct {
		IsRegistered bool `json:"is_registered"`
	}

	expBody := true
	expStatusCode := http.StatusOK

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expBody, respBody.IsRegistered)
	require.Equal(t, expStatusCode, resp.StatusCode)

	return email, password
}

func loginUser(t *testing.T, email string, pass string) string {
	body := map[string]string{
		"email":    email,
		"password": pass,
	}

	b, err := json.Marshal(body)
	require.NoError(t, err)

	resp, err := http.Post(urlLogin, contentType, bytes.NewReader(b))
	require.NoError(t, err)
	defer resp.Body.Close()

	var respBody struct {
		User map[string]string `json:"user"`
	}

	expBody := map[string]string{
		"first_name":  "Ivan",
		"middle_name": "Ivanovich",
		"last_name":   "Ivanov",
	}
	expStatusCode := http.StatusOK

	var sessionId string
	cookies := resp.Cookies()
	for _, c := range cookies {
		if c.Name == "sessionId" {
			sessionId = c.Value
		}
	}

	require.NotEmpty(t, sessionId)
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Equal(t, expBody, respBody.User)
	require.Equal(t, expStatusCode, resp.StatusCode)

	return sessionId
}

func uniqueEmail() string {
	return fmt.Sprintf("%stest@gmail.com", uuid.NewString())
}

func uniqueProjectName() string {
	return fmt.Sprintf("Project-%s", uuid.NewString())
}

func createProject(t *testing.T, sessionId string, projName string) uint32 {
	body := map[string]string{
		"name": projName,
	}

	b, err := json.Marshal(body)
	require.NoError(t, err)

	jar, err := cookiejar.New(nil)
	require.NoError(t, err)

	cookies := []*http.Cookie{}
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
	defer resp.Body.Close()

	var respBody struct {
		ProjectId uint32 `json:"project_id"`
	}

	expStatusCode := http.StatusOK

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
	require.Greater(t, int(respBody.ProjectId), 0)
	require.Equal(t, expStatusCode, resp.StatusCode)

	return respBody.ProjectId
}
