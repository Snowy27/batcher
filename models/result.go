package models

import "fmt"

//Response interface for result of request
type Response interface {
	ProvideResult() (*Result, error)
}

//Result of request
type Result struct {
	Name string
	Body string
}

//ProvideResult provides result of API call
func (result *Result) ProvideResult() (*Result, error) {
	return result, nil
}

//RequestError error of particular request
type RequestError struct {
	Name string
	Err  error
}

func (err RequestError) Error() string {
	return err.Err.Error()
}

//ProvideResult returns result with name of request and request error
func (err RequestError) ProvideResult() (*Result, error) {
	return &Result{Name: err.Name}, err.Err
}

//DependencyError error of dependency request
type DependencyError struct {
	Name           string
	DependencyName string
}

func (err DependencyError) Error() string {
	return fmt.Sprintf("Dependency %s has failed", err.DependencyName)
}

//ProvideResult returns result with name of request and dependency error
func (err DependencyError) ProvideResult() (*Result, error) {
	return &Result{Name: err.Name}, err
}
