package models

import "fmt"

//Response interface for result of request
type Response interface {
	ProvideResult() (Result, error)
}

//Result of request
type Result struct {
	Name string
	Body string
}

func (result Result) ProvideResult() (Result, error) {
	return result, nil
}

type RequestError struct {
	Name string
	Err  error
}

func (err RequestError) Error() string {
	return err.Err.Error()
}

func (err RequestError) ProvideResult() (Result, error) {
	return Result{Name: err.Name}, err.Err
}

type DependencyError struct {
	Name           string
	DependencyName string
}

func (err DependencyError) Error() string {
	return fmt.Sprintf("Dependency %s has failed", err.DependencyName)
}

func (err DependencyError) ProvideResult() (Result, error) {
	return Result{Name: err.Name}, err
}
