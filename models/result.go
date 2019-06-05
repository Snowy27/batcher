package models

import "fmt"

//Result of request
type Result struct {
	Name  string
	Body  string
	Error error
}

//RequestError error of particular request
type RequestError struct {
	Message string
}

func NewRequestError(message string) *RequestError {
	return &RequestError{Message: message}
}

func (err RequestError) Error() string {
	return err.Message
}

//DependencyError error of dependency request
type DependencyError struct {
	DependencyName string
}

func NewDependencyError(dependencyName string) *DependencyError {
	return &DependencyError{DependencyName: dependencyName}
}

func (err DependencyError) Error() string {
	return fmt.Sprintf("Dependency %s has failed", err.DependencyName)
}
