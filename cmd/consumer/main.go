package main

import (
	"go-backend/internal/kafka"
)

func main() {
	kafka.StartConsumer()

	select {} // keep running
}