package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"storage/common"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/uuid"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	queueURL := os.Getenv("SQS_QUEUE_URL")

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Printf("Unable to load SDK config, %v", err)
		return
	}

	sqsClient := sqs.NewFromConfig(cfg)

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Usage: add-item <key> <value>, get-item <key>, delete-item <key>, get-all-items")

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			break
		}

		parts := strings.Split(input, " ")

		cmdType := common.CommandType(parts[0])
		switch cmdType {
		case common.AddItem:
			if len(parts) < 3 {
				fmt.Println("Invalid command. Usage: add-item <key> <value>")
				continue
			}
		case common.GetItem:
			if len(parts) < 2 {
				fmt.Println("Invalid command. Usage: get-item <key>")
				continue
			}
		case common.DeleteItem:
			if len(parts) < 2 {
				fmt.Println("Invalid command. Usage: delete-item <key>")
				continue
			}
		case common.GetAllItems:
			if len(parts) > 1 {
				fmt.Println("Invalid command. Usage: get-all-items")
				continue
			}
		default:
			fmt.Println("Invalid command. Usage: add-item <key> <value>, get-item <key>, delete-item <key>, get-all-items")
			continue
		}

		var key string
		if len(parts) > 1 {
			key = parts[1]
		}

		var value string
		if len(parts) > 2 {
			value = parts[2]
		}

		command := common.Command{Type: cmdType, Key: key, Value: value, Id: uuid.New().String()}

		cmdJSON, err := json.Marshal(command)
		if err != nil {
			fmt.Printf("Failed to marshal command: %v", err)
			continue
		}

		mg := "commands"
		_, err = sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
			QueueUrl:       &queueURL,
			MessageBody:    aws.String(string(cmdJSON)),
			MessageGroupId: &mg,
		})

		if err != nil {
			fmt.Printf("Failed to send message to queue: %v", err)
			continue
		}

		fmt.Println("Command sent!")
	}
}
