package models

import (
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
func (request Request) Execute(dependencies []<-chan *Result, dependents []chan<- *Result) <-chan *Result {
	resultChannel := make(chan *Result)
	go func() {
		var result *Result

		_, dependencyError := getDependenciesResults(dependencies, request.Name)

		if dependencyError != nil {
			result = &Result{Name: request.Name, Error: dependencyError}
		} else {
			result = &Result{Name: request.Name, Body: fmt.Sprintf("Body of %s", request.Name)}
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

func getDependenciesResults(dependencies []<-chan *Result, name string) (results []*Result, dependencyError *DependencyError) {
	for _, dependency := range dependencies {
		result := <-dependency
		if result.Error != nil {
			dependencyError = NewDependencyError(result.Name)
		}
		results = append(results, result)
	}
	return
}
