//go:build integration
// +build integration

package test

import (
	"io/ioutil"
	"net/http"

	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestHealthCheck() {
	resp, err := http.Get(host + "/")
	assert.Nil(ts.T(), err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(ts.T(), err)

	assert.Equal(ts.T(), http.StatusOK, resp.StatusCode, string(body))
}
