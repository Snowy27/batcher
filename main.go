package main

import (
	"github.com/Snowy27/batcher/router"
)

func main() {
	server := router.InitRouter()
	server.Run(":8000")
}
