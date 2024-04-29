package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func connect() (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial("amqp://guest:guest@0.0.0.0:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	return conn, channel
}

func syncAccounts() {
	conn, channel := connect()
	defer conn.Close()
	defer channel.Close()

	for {
		accounts, err := pgInstance.getAccountsForProvisioning(context.Background())
		retryAccounts, retryErr := pgInstance.getAccountsForRetryProvisioning(context.Background())

		if err != nil || retryErr != nil {
			log.Printf("Error while syncing accounts: %s", err)
		} else {
			accounts = append(accounts, retryAccounts...)
			if len(accounts) > 0 {
				log.Printf("Syncing %s accounts", fmt.Sprint(len(accounts)))
			}
			for i := range accounts {
				account := accounts[i]
				body, err := json.Marshal(account)
				if err != nil {
					log.Printf("Error in json marshal: %s", err)
				} else {
					queue, err := channel.QueueDeclare(
						"accountProvision_"+account.SystemId.String,
						false,
						false,
						false,
						false,
						nil,
					)
					failOnError(err, "Failed to declare a queue")
					err = channel.PublishWithContext(context.Background(),
						"",         // exchange
						queue.Name, // routing key
						false,      // mandatory
						false,      // immediate
						amqp.Publishing{
							ContentType: "text/plain",
							Body:        []byte(body),
						})
					if err != nil {
						log.Printf("Error while sending message: %s", err)
					} else {
						pgInstance.markAccountAsProvisioned(accounts[i].AccountID)
					}
				}
			}
		}

		time.Sleep(5 * time.Second)
	}
}

func commitAccounts() {
	conn, channel := connect()
	defer conn.Close()
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"accountCommit",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	for {
		messages, err := channel.Consume(
			queue.Name,
			"",
			true,
			false,
			false,
			false,
			nil,
		)
		failOnError(err, "Failed to register a consumer")

		for message := range messages {
			log.Printf("Received a message: %s", message.Body)
			var account AccountProvision
			json.Unmarshal(message.Body, &account)
			pgInstance.markAccountAsCommitted(account.AccountID)
		}

		time.Sleep(1 * time.Second)
	}
}
