package models

import (
	"fmt"
	"sync"
)

//Payload that contains all requests
type Payload struct {
	Requests []Request `json:"requests" binding:"required,gt=0,dive"`
}

type WeightedRequest struct {
	Weight  int
	Request Request
}

//Execute the batched requests
func (payload *Payload) Execute() {
	requestsByWeight, maxWeight := calculateRequestsByWeightAndMaxWeight(payload.Requests)
	fmt.Println(requestsByWeight)
	fmt.Println(maxWeight)
}

func calculateRequestsByWeightAndMaxWeight(requests []Request) (map[int][]Request, int) {
	weights := createWeightsChannel(requests)

	requestsByWeight := make(map[int][]Request)
	maxWeight := 1
	for weight := range weights {
		if _, ok := requestsByWeight[weight.Weight]; !ok {
			requestsByWeight[weight.Weight] = make([]Request, 0)
		}
		requestsByWeight[weight.Weight] = append(requestsByWeight[weight.Weight], weight.Request)

		if maxWeight < weight.Weight {
			maxWeight = weight.Weight
		}
	}
	return requestsByWeight, maxWeight
}

func createWeightsChannel(requests []Request) <-chan WeightedRequest {
	listeners, senders := createListenersAndSenders(requests)
	weights := make([]<-chan WeightedRequest, 0)

	for _, request := range requests {
		weightChannel := make(chan WeightedRequest)
		weights = append(weights, weightChannel)
		go func(request Request) {

			weight := 1
			for _, listenerChannel := range listeners[request.Name] {
				weight += <-listenerChannel
			}

			for _, senderChannel := range senders[request.Name] {
				senderChannel <- weight
				close(senderChannel)
			}

			weightChannel <- WeightedRequest{Request: request, Weight: weight}
			close(weightChannel)

		}(request)
	}

	return merge(weights)
}

func createListenersAndSenders(requests []Request) (map[string][]<-chan int, map[string][]chan<- int) {
	//TODO: figure out circular dependencies
	listeners := make(map[string][]<-chan int)
	senders := make(map[string][]chan<- int)

	for _, req := range requests {
		if _, ok := listeners[req.Name]; !ok {
			listeners[req.Name] = make([]<-chan int, 0)
		}

		if _, ok := senders[req.Name]; !ok {
			senders[req.Name] = make([]chan<- int, 0)
		}

		for _, dependency := range req.Dependencies {
			//TODO: figure out case where dependency is not present
			ch := make(chan int)
			listeners[req.Name] = append(listeners[req.Name], ch)
			senders[dependency] = append(senders[dependency], ch)
		}
	}

	return listeners, senders
}

func merge(cs []<-chan WeightedRequest) <-chan WeightedRequest {
	var wg sync.WaitGroup
	out := make(chan WeightedRequest)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan WeightedRequest) {
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
