package models

import (
	"sync"
)

//Payload that contains all requests
type Payload struct {
	Requests []Request `json:"requests" binding:"required,gt=0,dive"`
}

//Execute the batched requests
func (payload *Payload) Execute() map[string]interface{} {
	listeners, senders := createListenersAndSenders(payload.Requests)
	responseChannels := make([]<-chan Response, 0, len(payload.Requests))

	for _, request := range payload.Requests {
		responseChannel := request.Execute(listeners[request.Name], senders[request.Name])
		responseChannels = append(responseChannels, responseChannel)
	}

	mergedResponsesChannel := merge(responseChannels)
	responses := make(map[string]interface{})

	for response := range mergedResponsesChannel {
		result, err := response.ProvideResult()
		if err != nil {
			responses[result.Name] = map[string]string{"Error": err.Error()}
		} else {
			responses[result.Name] = map[string]string{"Result": result.Body}
		}
	}

	return responses

}

func createListenersAndSenders(requests []Request) (map[string][]<-chan Response, map[string][]chan<- Response) {
	//TODO: figure out circular dependencies
	listeners := make(map[string][]<-chan Response)
	senders := make(map[string][]chan<- Response)

	for _, req := range requests {
		if _, ok := listeners[req.Name]; !ok {
			listeners[req.Name] = make([]<-chan Response, 0)
		}

		if _, ok := senders[req.Name]; !ok {
			senders[req.Name] = make([]chan<- Response, 0)
		}

		for _, dependency := range req.Dependencies {
			//TODO: figure out case where dependency is not present
			ch := make(chan Response)
			listeners[req.Name] = append(listeners[req.Name], ch)
			senders[dependency] = append(senders[dependency], ch)
		}
	}

	return listeners, senders
}

func merge(cs []<-chan Response) <-chan Response {
	var wg sync.WaitGroup
	out := make(chan Response)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan Response) {
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
