package models

import "fmt"

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

//Result of request
type Result struct {
	Name                string
	Body                string
	DependenciesResults []Result
}

//Execute request
func (request Request) Execute(dependencies []<-chan Result, dependents []chan<- Result) <-chan Result {
	resultChannel := make(chan Result)
	go func() {
		dependenciesResults := make([]Result, 0, len(dependencies))
		for _, dependency := range dependencies {
			dependenciesResults = append(dependenciesResults, <-dependency)
		}

		result := Result{Name: request.Name, Body: fmt.Sprintf("I am a result of %s", request.Name), DependenciesResults: dependenciesResults}

		for _, dependent := range dependents {
			dependent <- result
			close(dependent)
		}

		resultChannel <- result
		close(resultChannel)

	}()

	return resultChannel
}
