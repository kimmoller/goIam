package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func connect() (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
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
		syncCreateAccount(channel)
		syncEnableAccount(channel)
		syncDisableAccount(channel)
		syncDeleteAccount(channel)

		time.Sleep(5 * time.Second)
	}
}

func sendAccountsToQueueu(channel *amqp.Channel, queueName string, accounts []AccountProvision) {
	if len(accounts) > 0 {
		log.Printf("Sending %s accounts to queueu %s", fmt.Sprint(len(accounts)), queueName)
	}
	for i := range accounts {
		account := accounts[i]
		body, err := json.Marshal(account)
		if err != nil {
			log.Printf("Error in json marshal: %s", err)
		} else {
			queue, err := channel.QueueDeclare(
				queueName+"_"+account.SystemId.String,
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
				if queueName == "accountCreate" {
					pgInstance.markAccountAsProvisioned(accounts[i].AccountID)
				}
				if queueName == "accountEnable" {
					pgInstance.markAccountEnableAsProvisioned(accounts[i].AccountID)
				}
				if queueName == "accountDisable" {
					pgInstance.markAccountDisableAsProvisioned(accounts[i].AccountID)
				}
				if queueName == "accountDelete" {
					pgInstance.markAccountDeleteAsProvisioned(accounts[i].AccountID)
				}
			}
		}
	}
}

func syncCreateAccount(channel *amqp.Channel) {
	accounts, err := pgInstance.getAccountsForProvisioning(context.Background())
	retryAccounts, retryErr := pgInstance.getAccountsForRetryProvisioning(context.Background())

	if err != nil || retryErr != nil {
		log.Printf("Error while syncing accounts: %s", err)
	} else {
		accounts = append(accounts, retryAccounts...)
		sendAccountsToQueueu(channel, "accountCreate", accounts)
	}
}

func syncEnableAccount(channel *amqp.Channel) {
	accounts, err := pgInstance.getEnableAccountsForProvisioning(context.Background())
	retryAccounts, retryErr := pgInstance.getEnableAccountsForRetryProvisioning(context.Background())

	if err != nil || retryErr != nil {
		log.Printf("Error while syncing accounts: %s", err)
	} else {
		accounts = append(accounts, retryAccounts...)
		sendAccountsToQueueu(channel, "accountEnable", accounts)
	}
}

func syncDisableAccount(channel *amqp.Channel) {
	accounts, err := pgInstance.getDisableAccountsForProvisioning(context.Background())
	retryAccounts, retryErr := pgInstance.getDisableAccountsForRetryProvisioning(context.Background())

	if err != nil || retryErr != nil {
		log.Printf("Error while syncing accounts: %s", err)
	} else {
		accounts = append(accounts, retryAccounts...)
		sendAccountsToQueueu(channel, "accountDisable", accounts)
	}
}

func syncDeleteAccount(channel *amqp.Channel) {
	accounts, err := pgInstance.getDeleteAccountsForProvisioning(context.Background())
	retryAccounts, retryErr := pgInstance.getDeleteAccountsForRetryProvisioning(context.Background())

	if err != nil || retryErr != nil {
		log.Printf("Error while syncing accounts: %s", err)
	} else {
		accounts = append(accounts, retryAccounts...)
		sendAccountsToQueueu(channel, "accountDelete", accounts)
	}
}

func declareQueue(channel *amqp.Channel, queueuName string) (*amqp.Queue, error) {
	queue, err := channel.QueueDeclare(
		queueuName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queueu %s, %w", queueuName, err)
	}
	return &queue, nil
}

func consumeMessages(channel *amqp.Channel, queue *amqp.Queue) (<-chan amqp.Delivery, error) {
	messages, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queueu %s, %w", queue.Name, err)
	}
	return messages, nil
}

func commitAccounts() {
	conn, channel := connect()
	defer conn.Close()
	defer channel.Close()

	createQueue, err := declareQueue(channel, "accountCreateCommit")
	if err != nil {
		log.Printf("Failed to declare create queue %s", err)
	}
	enableQueue, err := declareQueue(channel, "accountEnableCommit")
	if err != nil {
		log.Printf("Failed to declare enable queue %s", err)
	}
	disableQueue, err := declareQueue(channel, "accountDisableCommit")
	if err != nil {
		log.Printf("Failed to declare disable queue %s", err)
	}
	deleteQueue, err := declareQueue(channel, "accountDeleteCommit")
	if err != nil {
		log.Printf("Failed to declare delete queue %s", err)
	}
	log.Println("Ready to receive commit messages...")

	var forever chan struct{}

	go func() {
		for {
			commitAccountCreation(channel, createQueue)
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			commitAccountEnable(channel, enableQueue)
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			commitAccountDisable(channel, disableQueue)
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			commitAccountDelete(channel, deleteQueue)
			time.Sleep(1 * time.Second)
		}
	}()

	<-forever
}

func commitAccountCreation(channel *amqp.Channel, queue *amqp.Queue) {
	messages, err := consumeMessages(channel, queue)

	if err != nil {
		log.Printf("Error while consuming queue %s, %s", queue.Name, err)
	} else {
		for message := range messages {
			log.Printf("Received a message: %s", message.Body)
			var account AccountProvision
			json.Unmarshal(message.Body, &account)
			pgInstance.markAccountAsCommitted(account.AccountID)
		}
	}
}

func commitAccountEnable(channel *amqp.Channel, queue *amqp.Queue) {
	messages, err := consumeMessages(channel, queue)

	if err != nil {
		log.Printf("Error while consuming queue %s, %s", queue.Name, err)
	} else {
		for message := range messages {
			log.Printf("Received a message: %s", message.Body)
			var account AccountProvision
			json.Unmarshal(message.Body, &account)
			pgInstance.markAccountEnableAsCommitted(account.AccountID)
		}
	}
}

func commitAccountDisable(channel *amqp.Channel, queue *amqp.Queue) {
	messages, err := consumeMessages(channel, queue)

	if err != nil {
		log.Printf("Error while consuming queue %s, %s", queue.Name, err)
	} else {
		for message := range messages {
			log.Printf("Received a message: %s", message.Body)
			var account AccountProvision
			json.Unmarshal(message.Body, &account)
			pgInstance.markAccountDisableAsCommitted(account.AccountID)
		}
	}
}

func commitAccountDelete(channel *amqp.Channel, queue *amqp.Queue) {
	messages, err := consumeMessages(channel, queue)

	if err != nil {
		log.Printf("Error while consuming queue %s", queue.Name)
	} else {
		for message := range messages {
			log.Printf("Received a message: %s", message.Body)
			var account AccountProvision
			json.Unmarshal(message.Body, &account)
			pgInstance.markAccountDeleteAsCommitted(account.AccountID)
		}
	}
}
