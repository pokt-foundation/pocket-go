package mock

import (
	"io/ioutil"

	"github.com/jarcoal/httpmock"
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
