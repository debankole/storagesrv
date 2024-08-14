package sqs

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"storage/common"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var closing bool

type MessageReader struct {
	sqsClient *sqs.Client
	queueURL  string
	commandCh chan common.Command
	wg        sync.WaitGroup
}

func NewMessageReader(queueURL string) (*MessageReader, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	sqsClient := sqs.NewFromConfig(cfg)
	return &MessageReader{
		sqsClient: sqsClient,
		queueURL:  queueURL,
		commandCh: make(chan common.Command),
	}, nil
}

func (s *MessageReader) Start(workers int) {
	log.Println("Starting SQS reader...")
	for i := 0; i < workers; i++ {
		s.wg.Add(1)
		go s.worker()
	}
}

func (s *MessageReader) worker() {
	defer s.wg.Done()

	s.receiveMessages()
}

func (s *MessageReader) receiveMessages() {
	for !closing {
		log.Printf("Receiving messages from queue: %s", s.queueURL)
		output, err := s.sqsClient.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl:            &s.queueURL,
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     10,
		})

		if err != nil {
			log.Fatalf("Unable to receive message from queue: %v", err)
		}

		for _, message := range output.Messages {
			var cmd common.Command
			err := json.Unmarshal([]byte(*message.Body), &cmd)
			log.Printf("Received message: %s", cmd.Type)
			if err != nil {
				log.Printf("Failed to unmarshal command: %v", err)
				continue
			}

			s.commandCh <- cmd

			s.sqsClient.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
				QueueUrl:      &s.queueURL,
				ReceiptHandle: message.ReceiptHandle,
			})
		}
	}
}

func (s *MessageReader) ReadMessages() <-chan common.Command {
	return s.commandCh
}

func (s *MessageReader) Stop() {
	log.Println("Stopping SQS reader...")

	// stop the receiveMessages loop
	closing = true
	// wait for all workers to finish
	s.wg.Wait()
	// now it's safe to close the command channel
	close(s.commandCh)
}
