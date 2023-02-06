//go:build integration
// +build integration

package test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var adminToken string

func adminLogin(ts *TestSuite) {
	if adminToken != "" {
		return
	}
	request := map[string]string{
		"email":    adminEmail,
		"password": adminPassword,
	}
	reqBody, err := json.Marshal(request)
	assert.Nil(ts.T(), err)

	resp, err := http.Post(host+"/login", "application/json", bytes.NewBuffer(reqBody))
	assert.Nil(ts.T(), err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(ts.T(), err)

	assert.Equal(ts.T(), http.StatusOK, resp.StatusCode, string(body))

	var response map[string]string
	err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&response)
	assert.Nil(ts.T(), err)

	assert.NotEmpty(ts.T(), response["token"])
	adminToken = response["token"]
}

func (ts *TestSuite) TestAdminGetProfile() {
	adminLogin(ts)

	req, err := http.NewRequest("GET", host+"/me", nil)
	assert.Nil(ts.T(), err)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	assert.Nil(ts.T(), err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(ts.T(), err)

	assert.Equal(ts.T(), http.StatusOK, resp.StatusCode, string(body))

	var response map[string]interface{}
	err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&response)
	assert.Nil(ts.T(), err)

	assert.Equal(ts.T(), adminEmail, response["email"])
	assert.Equal(ts.T(), adminName, response["name"])
	assert.Equal(ts.T(), adminId.String(), response["id"])
}

func (ts *TestSuite) TestAdminMakeUser() {
	adminLogin(ts)
	request := map[string]string{
		"name":     userName,
		"email":    userEmail,
		"password": userPassword,
	}
	reqBody, err := json.Marshal(request)
	assert.Nil(ts.T(), err)

	req, err := http.NewRequest("POST", host+"/admin/users", bytes.NewBuffer(reqBody))
	assert.Nil(ts.T(), err)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	assert.Nil(ts.T(), err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(ts.T(), err)

	assert.Equal(ts.T(), http.StatusCreated, resp.StatusCode, string(body))

	var response map[string]string
	err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&response)
	assert.Nil(ts.T(), err)

	assert.NotEmpty(ts.T(), response["uuid"])
	userId = uuid.MustParse(response["uuid"])
}

func (ts *TestSuite) TestAdminLogout() {
	adminLogin(ts)
	req, err := http.NewRequest("POST", host+"/logout", nil)
	assert.Nil(ts.T(), err)

	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	assert.Nil(ts.T(), err)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(ts.T(), err)

	assert.Equal(ts.T(), http.StatusOK, resp.StatusCode, string(body))
	var response map[string]interface{}
	err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&response)
	assert.Nil(ts.T(), err)
	assert.Equal(ts.T(), "logout success", response["message"])
	adminToken = ""
}
