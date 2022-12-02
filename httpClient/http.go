package httpClient

import (
	"net/http"
)

func newHttp() *http.Client {
	return NewClient(ClientParams{})
}

func SendRequest(url string) (*http.Response, error) {
	request, err := NewRequest(Request{
		Headers: nil,
		Url:     url,
		Method:  "GET",
		Body:    nil,
	})

	if err != nil {
		return nil, err
	}

	response, err := Make(request, newHttp())

	if err != nil {
		return nil, err
	}

	return response, nil
}
