package integrationtests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

var apiHost = "http://0.0.0.0:8888"

// var apiHost = "http://localhost:50052"

var (
	EventsURL = fmt.Sprintf("%s/events", apiHost)
	UsersURL  = fmt.Sprintf("%s/users", apiHost)
)

func makeRequest(url string, method string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(context.Background(), method, url, body)
	if err != nil {
		return
	}
	req.Header.Add("Content-type", "application/json")
	return http.DefaultClient.Do(req)
}

func PrepareItem(url, method string, payload interface{}) (statusCode int, content []byte, err error) {
	bts, err := json.Marshal(payload)
	if err != nil {
		return
	}
	// fmt.Printf("Request to %s %s\n", url, string(bts))
	resp, err := makeRequest(url, method, bytes.NewReader(bts))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	statusCode = resp.StatusCode
	content, err = ioutil.ReadAll(resp.Body)
	// fmt.Printf("Response %s\n", string(content))
	return
}

func CreateUser(url string, payload User) (statusCode int, result User, err error) {
	statusCode, content, err := PrepareItem(url, http.MethodPost, payload)
	if err != nil {
		return
	}
	err = json.Unmarshal(content, &result)
	return
}

func ProcessEvent(url, method string, payload interface{}) (statusCode int, result Event, err error) {
	statusCode, content, err := PrepareItem(url, method, payload)
	if err != nil {
		return
	}
	err = json.Unmarshal(content, &result)
	return
}

func ProcessEvents(url, method string, payload interface{}) (statusCode int, result QueryEventsResponse, err error) {
	statusCode, content, err := PrepareItem(url, method, payload)
	if err != nil {
		return
	}
	err = json.Unmarshal(content, &result)
	return
}

func ProcessError(url, method string, payload interface{}) (statusCode int, result ErrorResponse, err error) {
	statusCode, content, err := PrepareItem(url, method, payload)
	if err != nil {
		return
	}
	err = json.Unmarshal(content, &result)
	return
}
