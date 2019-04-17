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
	Weight int
	Name   string
}

//Execute the batched requests
func (payload *Payload) Execute() {
	listeners := make(map[string][]chan int)
	senders := make(map[string][]chan int)
	requestsByName := make(map[string]Request)
	for _, req := range payload.Requests {
		requestsByName[req.Name] = req
		if _, ok := listeners[req.Name]; !ok {
			listeners[req.Name] = make([]chan int, 0)
		}
		if _, ok := senders[req.Name]; !ok {
			senders[req.Name] = make([]chan int, 0)
		}
		for _, dependency := range req.Dependencies {
			ch := make(chan int)
			listeners[req.Name] = append(listeners[req.Name], ch)
			senders[dependency] = append(senders[dependency], ch)
		}
	}
	result := calculateWeights(requestsByName, senders, listeners)
	fmt.Println(result)
}

func calculateWeights(requestsByName map[string]Request, senders, listeners map[string][]chan int) map[string]int {
	weights := make(map[string]int)
	weightChannels := make([]<-chan WeightedRequest, 0)
	for name := range requestsByName {
		weightChannel := make(chan WeightedRequest)
		weightChannels = append(weightChannels, weightChannel)
		go func(name string) {
			fmt.Println("Starting routine for ", name)
			weight := 1
			for _, listenerChannel := range listeners[name] {
				fmt.Println("Waiting on ", name)
				weight += <-listenerChannel
			}
			for _, senderChannel := range senders[name] {
				fmt.Println("Sending the ", name)
				senderChannel <- weight
				close(senderChannel)
			}
			weightChannel <- WeightedRequest{Name: name, Weight: weight}
			close(weightChannel)
		}(name)
	}

	mergedChannel := merge(weightChannels)

	for weight := range mergedChannel {
		weights[weight.Name] = weight.Weight
	}
	return weights
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
