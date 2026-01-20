package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	contentType = "application/json"
	url         = "http://localhost:44044/registration"
)

func TestRegistration_Success_Integration(t *testing.T) {
	body := map[string]string{
		"first_name":  "Ivan",
		"middle_name": "Ivanovich",
		"last_name":   "Ivanov",
		"password":    "somePass",
		"email":       uniqueEmail(),
	}

	b, err := json.Marshal(body)
	require.NoError(t, err)

	resp, err := http.Post(url, contentType, bytes.NewReader(b))
	require.NoError(t, err)
	defer resp.Body.Close()

	var reqBody struct {
		IsRegistered bool `json:"is_registered"`
	}
	expStatusCode := http.StatusOK

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&reqBody))
	require.True(t, reqBody.IsRegistered)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestRegistration_BadRequest_Integration(t *testing.T) {
	body := map[string]string{
		"middle_name": "Ivanovich",
		"last_name":   "Ivanov",
		"password":    "somePass",
		"email":       uniqueEmail(),
	}

	b, err := json.Marshal(body)
	require.NoError(t, err)

	resp, err := http.Post(url, contentType, bytes.NewReader(b))
	require.NoError(t, err)
	defer resp.Body.Close()

	expBody := map[string]string{"FirstName": "field is required"}
	expStatusCode := http.StatusBadRequest

	var reqBody struct {
		Errors map[string]string `json:"errors"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&reqBody))
	require.Equal(t, expBody, reqBody.Errors)
	require.Equal(t, expStatusCode, resp.StatusCode)
}

func TestRegistration_AlreadyExists_Integration(t *testing.T) {
	body := map[string]string{
		"first_name":  "Ivan",
		"middle_name": "Ivanovich",
		"last_name":   "Ivanov",
		"password":    "somePass",
		"email":       uniqueEmail(),
	}

	b, err := json.Marshal(body)
	require.NoError(t, err)

	resp1, err := http.Post(url, contentType, bytes.NewReader(b))
	require.NoError(t, err)
	defer resp1.Body.Close()

	var reqBody1 struct {
		IsRegistered bool `json:"is_registered"`
	}
	expStatusCode := http.StatusOK

	require.NoError(t, json.NewDecoder(resp1.Body).Decode(&reqBody1))
	require.True(t, reqBody1.IsRegistered)
	require.Equal(t, expStatusCode, resp1.StatusCode)

	resp2, err := http.Post(url, contentType, bytes.NewReader(b))
	require.NoError(t, err)
	defer resp2.Body.Close()

	var reqBody2 struct {
		ExpErr string `json:"error"`
	}
	expBody2 := "user already exists"
	expStatusCode = http.StatusConflict

	require.NoError(t, json.NewDecoder(resp2.Body).Decode(&reqBody2))
	require.Equal(t, expBody2, reqBody2.ExpErr)
	require.Equal(t, expStatusCode, resp2.StatusCode)
}

func uniqueEmail() string {
	return "testgmail" + strconv.Itoa(int(time.Now().Unix())) + "@gmail.com"
}
