package main

import (
	"neatly/internal/router"
	"neatly/pkg/logging"
)

func main() {
	logging.Init()

	logger := logging.GetLogger()
	logger.Println("Logger initialized.")

	logger.Println("Application context initialized.")

	defer router.Init()

	logger.Println("Application started.")
}
