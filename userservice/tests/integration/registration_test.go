package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
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

	resp, err := http.Post(urlReg, contentType, bytes.NewReader(b))
	require.NoError(t, err)
	defer resp.Body.Close()

	var resBody struct {
		IsRegistered bool `json:"is_registered"`
	}
	expStatusCode := http.StatusOK

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&resBody))
	require.True(t, resBody.IsRegistered)
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

	resp, err := http.Post(urlReg, contentType, bytes.NewReader(b))
	require.NoError(t, err)
	defer resp.Body.Close()

	expBody := map[string]string{"FirstName": "field is required"}
	expStatusCode := http.StatusBadRequest

	var resBody struct {
		Errors map[string]string `json:"errors"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&resBody))
	require.Equal(t, expBody, resBody.Errors)
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

	resp1, err := http.Post(urlReg, contentType, bytes.NewReader(b))
	require.NoError(t, err)
	defer resp1.Body.Close()

	var resBody1 struct {
		IsRegistered bool `json:"is_registered"`
	}
	expStatusCode := http.StatusOK

	require.NoError(t, json.NewDecoder(resp1.Body).Decode(&resBody1))
	require.True(t, resBody1.IsRegistered)
	require.Equal(t, expStatusCode, resp1.StatusCode)

	resp2, err := http.Post(urlReg, contentType, bytes.NewReader(b))
	require.NoError(t, err)
	defer resp2.Body.Close()

	var resBody2 struct {
		ExpErr string `json:"error"`
	}
	expBody2 := "user already exists"
	expStatusCode = http.StatusConflict

	require.NoError(t, json.NewDecoder(resp2.Body).Decode(&resBody2))
	require.Equal(t, expBody2, resBody2.ExpErr)
	require.Equal(t, expStatusCode, resp2.StatusCode)
}
