package models

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
}
