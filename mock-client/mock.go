package mock

import (
	"errors"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/jarcoal/httpmock"
)

var (
	// ErrResponseNotFound when try get a response not avaible
	ErrResponseNotFound = errors.New("response not found")
)

// AddMockedResponseFromFile adds a mocked response given a file path relative to the test file
func AddMockedResponseFromFile(method string, url string, statusCode int, filePath string) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	AddMockedResponse(method, url, statusCode, string(data))
}

// AddMockedResponse adds a mocked response given its content
func AddMockedResponse(method string, url string, statusCode int, content string) {
	responder := httpmock.NewStringResponder(statusCode, content)
	httpmock.RegisterResponder(method, url, responder)
}

// AddMultipleMockedResponses add a mocked response given one to one from each file
func AddMultipleMockedResponses(method string, url string, statusCode int, responseList []string) {
	var mutex = sync.Mutex{}

	nextResponseIndex := 0
	responseFunction := func(req *http.Request) (*http.Response, error) {
		mutex.Lock()
		defer mutex.Unlock()

		if nextResponseIndex >= len(responseList) {
			return nil, ErrResponseNotFound
		}

		data, err := ioutil.ReadFile(responseList[nextResponseIndex])
		if err != nil {
			panic(err)
		}

		req.Response = httpmock.NewStringResponse(statusCode, string(data))

		nextResponseIndex++

		return req.Response, nil
	}

	httpmock.RegisterResponder(method, url, responseFunction)
}
