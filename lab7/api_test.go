package lab7

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUsers(t *testing.T) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://reqres.in/api/users?page=1", nil)
	assert.NoError(t, err)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Contains(t, result, "data")
}

func TestRegisterUser(t *testing.T) {
	client := &http.Client{}
	payload := map[string]string{
		"email":    "eve.holt@reqres.in",
		"password": "pistol",
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "https://reqres.in/api/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Contains(t, result, "token")
}

func TestRegisterUserNegative(t *testing.T) {
	client := &http.Client{}
	payload := map[string]string{
		"email": "eve.holt@reqres.in",
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "https://reqres.in/api/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Contains(t, result, "error")
	assert.Equal(t, "Missing password", result["error"])
}

func TestDeleteUser(t *testing.T) {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", "https://reqres.in/api/users/2", nil)
	assert.NoError(t, err)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestGetUserById(t *testing.T) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://reqres.in/api/users/2", nil)
	assert.NoError(t, err)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	data := result["data"].(map[string]interface{})
	assert.Equal(t, 2, int(data["id"].(float64)))
	assert.NotEmpty(t, data["email"])
	assert.NotEmpty(t, data["first_name"])
}

func TestGetUserByIdNotFound(t *testing.T) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://reqres.in/api/users/23", nil)
	assert.NoError(t, err)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestLoginUser(t *testing.T) {
	client := &http.Client{}
	payload := map[string]string{
		"email":    "eve.holt@reqres.in",
		"password": "cityslicka",
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "https://reqres.in/api/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Contains(t, result, "token")
}

func TestLoginUserMissingPassword(t *testing.T) {
	client := &http.Client{}
	payload := map[string]string{
		"email": "eve.holt@reqres.in",
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "https://reqres.in/api/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Contains(t, result, "error")
	assert.Equal(t, "Missing password", result["error"])
}

func TestUpdateUser(t *testing.T) {
	client := &http.Client{}
	payload := map[string]string{
		"name": "morpheus",
		"job":  "zion resident",
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("PUT", "https://reqres.in/api/users/2", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, "morpheus", result["name"])
	assert.Equal(t, "zion resident", result["job"])
	assert.NotEmpty(t, result["updatedAt"])
}

func TestRegisterUserMissingEmail(t *testing.T) {
	client := &http.Client{}
	payload := map[string]string{
		"password": "pistol",
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "https://reqres.in/api/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Contains(t, result, "error")
	assert.Equal(t, "Missing email or username", result["error"])
}

func TestGetUsersWithPageParam(t *testing.T) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://reqres.in/api/users?page=2", nil)
	assert.NoError(t, err)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, 2, int(result["page"].(float64)))
	data := result["data"].([]interface{})
	assert.Greater(t, len(data), 0)
}
func TestDeleteNonexistentUser(t *testing.T) {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", "https://reqres.in/api/users/9999", nil)
	assert.NoError(t, err)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}
