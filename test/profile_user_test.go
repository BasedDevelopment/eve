//go:build integration
// +build integration

package test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/stretchr/testify/assert"
)

var userToken string

func userLogin(ts *TestSuite) {
	if userToken != "" {
		return
	}
	request := map[string]string{
		"email":    userEmail,
		"password": userPassword,
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
	userToken = response["token"]
}

func (ts *TestSuite) TestUserGetProfile() {
	userLogin(ts)

	req, err := http.NewRequest("GET", host+"/me", nil)
	assert.Nil(ts.T(), err)
	req.Header.Set("Authorization", "Bearer "+userToken)
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

	assert.Equal(ts.T(), userEmail, response["email"])
	assert.Equal(ts.T(), userName, response["name"])
	assert.Equal(ts.T(), userId.String(), response["id"])
}

func (ts *TestSuite) TestUserLogout() {
	userLogin(ts)
	req, err := http.NewRequest("POST", host+"/logout", nil)
	assert.Nil(ts.T(), err)

	req.Header.Set("Authorization", "Bearer "+userToken)
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
	userToken = ""
}
