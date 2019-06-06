package models

import (
	"errors"
	"sync"
)

//Payload that contains all requests
type Payload struct {
	Requests []Request `json:"requests" binding:"required,gt=0,dive"`
}

//CheckForCircularDependencies verifies that there is no circular dependencies in the payload
func (payload *Payload) CheckForCircularDependencies() error {
	graph := NewDependenciesGraph(payload.Requests)

	if graph.CheckForCircularDependencies() {
		return errors.New("The payload has circular dependecies")
	}

	return nil
}

//Execute the batched requests
func (payload *Payload) Execute() map[string]interface{} {
	listeners, senders := createListenersAndSenders(payload.Requests)
	resultChannels := make([]<-chan *Result, 0, len(payload.Requests))

	for _, request := range payload.Requests {
		resultChannel := request.Execute(listeners[request.Name], senders[request.Name])
		resultChannels = append(resultChannels, resultChannel)
	}

	mergedResultsChannel := merge(resultChannels)
	responses := make(map[string]interface{})

	for result := range mergedResultsChannel {
		if result.Error != nil {
			responses[result.Name] = map[string]interface{}{"Error": result.Error.Error(), "StatusCode": result.StatusCode}
		} else {
			responses[result.Name] = map[string]interface{}{"StatusCode": result.StatusCode, "Response": result.Body}
		}
	}

	return responses

}

func createListenersAndSenders(requests []Request) (map[string][]<-chan *Result, map[string][]chan<- *Result) {
	//TODO: figure out circular dependencies
	listeners := make(map[string][]<-chan *Result)
	senders := make(map[string][]chan<- *Result)

	for _, req := range requests {
		if _, ok := listeners[req.Name]; !ok {
			listeners[req.Name] = make([]<-chan *Result, 0)
		}

		if _, ok := senders[req.Name]; !ok {
			senders[req.Name] = make([]chan<- *Result, 0)
		}

		for _, dependency := range req.Dependencies {
			//TODO: figure out case where dependency is not present
			ch := make(chan *Result)
			listeners[req.Name] = append(listeners[req.Name], ch)
			senders[dependency] = append(senders[dependency], ch)
		}
	}

	return listeners, senders
}

func merge(cs []<-chan *Result) <-chan *Result {
	var wg sync.WaitGroup
	out := make(chan *Result)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan *Result) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
