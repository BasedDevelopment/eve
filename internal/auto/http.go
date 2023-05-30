package auto

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

func (a *Auto) httpReq(method string, urlStr string, data any) (respBodyBytes []byte, status int, err error) {
	c := a.getHttpsClient()

	switch method {
	case "GET", "DELETE":
		resp := &http.Response{}
		var err error
		if method == "GET" {
			resp, err = c.Get(urlStr)
		}
		if method == "DELETE" {
			// Parse the urlStr
			u, err := url.Parse(urlStr)
			resp, err = c.Do(&http.Request{
				Method: "DELETE",
				URL:    u,
			})
		}
		if err != nil {
			return nil, -1, err
		}
		if resp.Body == nil {
			return nil, -1, errors.New("eve/auto: response body is nil")
		}
		defer resp.Body.Close()
		respBodyBytes, err = ioutil.ReadAll(resp.Body)
		status = resp.StatusCode
	case "POST", "PUT", "PATCH":
		if data == nil {
			return nil, -1, errors.New("eve/auto: data is nil")
		}

		url, err := url.Parse(urlStr)
		if err != nil {
			return nil, -1, err
		}

		reqBodyBytes, err := json.Marshal(data)
		if err != nil {
			return nil, -1, err
		}

		req := http.Request{
			Method: method,
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
			URL:  url,
			Body: ioutil.NopCloser(bytes.NewBuffer(reqBodyBytes)),
		}

		resp, err := c.Do(&req)
		if err != nil {
			return nil, -1, err
		}

		if resp.Body == nil {
			return nil, -1, errors.New("eve/auto: response body is nil")
		}
		defer resp.Body.Close()
		respBodyBytes, err = ioutil.ReadAll(resp.Body)
		status = resp.StatusCode
	default:
		return nil, -1, errors.New("eve/auto: invalid method")
	}

	return
}
