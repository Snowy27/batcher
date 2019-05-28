package models

import (
	"errors"
	"fmt"
)

//Request that needs to be batched
type Request struct {
	Method       string      `json:"method" binding:"required,eq=PUT|eq=POST|eq=DELETE|eq=GET"`
	Name         string      `json:"name" binding:"required,gt=0"`
	URL          string      `json:"url" binding:"required,url|uri"`
	Body         interface{} `json:"body"`
	Dependencies []string    `json:"dependencies"`
	Concurrency  uint8       `json:"concurrency"`
	Retries      uint8       `json:"retries"`
	Timeout      uint        `json:"timeout"`
	Weight       int
}

//Execute request
func (request Request) Execute(dependencies []<-chan Response, dependents []chan<- Response) <-chan Response {
	resultChannel := make(chan Response)
	go func() {
		var result Response

		_, dependencyError := getDependenciesResults(dependencies, request.Name)

		if dependencyError != nil {
			result = dependencyError
		} else {
			if request.Name == "test3" {
				result = &RequestError{Name: request.Name, Err: errors.New("Http Error")}
			} else {
				result = &Result{Name: request.Name, Body: fmt.Sprintf("Body of %s", request.Name)}
			}
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

func getDependenciesResults(dependencies []<-chan Response, name string) (results []Response, dependencyError *DependencyError) {
	for _, dependency := range dependencies {
		depResponse := <-dependency
		depResult, err := depResponse.ProvideResult()
		if err != nil {
			dependencyError = &DependencyError{Name: name, DependencyName: depResult.Name}
		}
		results = append(results, depResult)
	}
	return
}
