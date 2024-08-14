package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"storage/server/output"
	"storage/server/proc"
	"storage/server/sqs"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	queueURL := os.Getenv("SQS_QUEUE_URL")

	sqsReader, err := sqs.NewMessageReader(queueURL)
	if err != nil {
		log.Fatalf("Failed to initialize sqs reader: %v", err)
	}

	sqsReader.Start(2)

	out, err := output.NewFileWriter()
	if err != nil {
		log.Fatalf("Failed to initialize file writer: %v", err)
	}

	processor := proc.NewCommandProcessor(sqsReader, out)
	processor.Start(runtime.NumCPU())

	log.Printf("Server started with %d workers", runtime.NumCPU())

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	log.Println("Server is shutting down...")
	sqsReader.Stop()
	processor.Stop()
	out.Close()

	log.Println("Server stopped gracefully")
}
