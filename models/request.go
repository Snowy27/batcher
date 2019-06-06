package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

//Request that needs to be batched
type Request struct {
	Method       string                 `json:"method" binding:"required,eq=PUT|eq=POST|eq=DELETE|eq=GET"`
	Name         string                 `json:"name" binding:"required,gt=0"`
	URL          string                 `json:"url" binding:"required,url|uri"`
	Body         map[string]interface{} `json:"body" binding:"requiredwhenputorpost"`
	Dependencies []string               `json:"dependencies"`
	Concurrency  uint8                  `json:"concurrency"`
	Retries      uint8                  `json:"retries"`
	Timeout      uint                   `json:"timeout"`
}

//Execute request
func (request Request) Execute(dependencies []<-chan *Result, dependents []chan<- *Result) <-chan *Result {
	resultChannel := make(chan *Result)
	go func() {
		var result *Result

		_, dependencyError := getDependenciesResults(dependencies)

		if dependencyError != nil {
			result = &Result{Name: request.Name, Error: dependencyError}
		} else {
			// TODO: perform a call
			result = request.performAPICall()
		}

		for _, dependent := range dependents {
			dependent <- result
			close(dependent)
		}

		resultChannel <- result
		close(resultChannel)

	}()

	return resultChannel
}

func (request Request) performAPICall() *Result {
	fmt.Printf("Executing request %s\n", request.Name)

	var resp *http.Response
	var err error
	//TODO specify real timeout
	client := http.Client{Timeout: 5 * time.Second}
	switch request.Method {
	case "GET":
		resp, err = request.performGet(client)
	case "POST":
		resp, err = request.performPost(client)
	case "PUT":
		resp, err = request.performPut(client)
	case "DELETE":
		resp, err = request.performDelete(client)
	}

	defer resp.Body.Close()

	if err != nil {
		return &Result{Name: request.Name, Error: err, StatusCode: 500}
	}

	return request.handleResponse(resp)

}

func (request Request) performGet(client http.Client) (*http.Response, error) {
	req, err := http.NewRequest("GET", request.URL, nil)
	if err != nil {
		return nil, err
	}

	return client.Do(req)
}

func (request Request) performPost(client http.Client) (*http.Response, error) {
	bytesRepresentation, err := json.Marshal(request.Body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", request.URL, bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}

func (request Request) performPut(client http.Client) (*http.Response, error) {
	bytesRepresentation, err := json.Marshal(request.Body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", request.URL, bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}

func (request Request) performDelete(client http.Client) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", request.URL, nil)
	if err != nil {
		return nil, err
	}

	return client.Do(req)
}

func (request Request) handleResponse(resp *http.Response) *Result {

	if resp.StatusCode > 299 {
		errorMessage, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return &Result{Name: request.Name, Error: NewRequestError("Unable to get response"), StatusCode: resp.StatusCode}
		}
		return &Result{Name: request.Name, Error: NewRequestError(string(errorMessage)), StatusCode: resp.StatusCode}
	}

	var result map[string]interface{}

	err := json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		return &Result{Name: request.Name, Error: err, StatusCode: 500}
	}

	fmt.Printf("Finishing request %s\n", request.Name)
	return &Result{Name: request.Name, Body: result, StatusCode: resp.StatusCode}
}

func getDependenciesResults(dependencies []<-chan *Result) (results []*Result, dependencyError *DependencyError) {
	for _, dependency := range dependencies {
		result := <-dependency
		if result.Error != nil {
			dependencyError = NewDependencyError(result.Name)
		}
		results = append(results, result)
	}
	return
}
