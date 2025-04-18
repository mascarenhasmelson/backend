package consumer
import (
	"encoding/json"
	"log"
    "gitlab.com/melsonmascarenhas/backend/models"
	"gitlab.com/melsonmascarenhas/backend/mongodb"
	"github.com/rabbitmq/amqp091-go"
)

var(
	account models.PubAccount
	withdraw models.Transaction
)

// https://www.rabbitmq.com/tutorials/tutorial-one-go
func StartWithdrawConsumer() {
	// Connect to RabbitMQ
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Open a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare the queue
	q, err := ch.QueueDeclare(
		"withdraw", // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Consume messages from the queue
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// Process messages in a goroutine
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			// Parse the JSON data
			if err := json.Unmarshal(d.Body, &withdraw); err != nil {
				log.Printf("Failed to parse message: %v", err)
				continue
			}

			// Call the WithDraw function
			err := mongodb.WithDraw(withdraw.CustomerID, withdraw.Amount, withdraw.IsDeposit)
			if err != nil {
				log.Printf("Error processing transaction: %v", err)
			}
		}
	}()

	log.Println(" [*] Waiting for messages on the 'withdraw' queue. To exit press CTRL+C")
	<-forever
}
func StartCreateAccountConsumer() {
	// Connect to RabbitMQ
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Open a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare the queue
	q, err := ch.QueueDeclare(
		"createAccount", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Consume messages from the queue
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// Process messages in a goroutine
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			// Parse the JSON data
			// var account AccountData
			if err := json.Unmarshal(d.Body, &account); err != nil {
				log.Printf("Failed to parse message: %v", err)
				continue
			}

			// Call the CreateAccount function
			response, err := mongodb.CreateAccount(account.Name, account.CustomerID, account.Amount)
			if err != nil {
				log.Printf("Error creating account: %v", err)
			} else {
				log.Printf("Success: %s", response)
			}
		}
	}()

	log.Println(" [*] Waiting for messages on the 'createAccount' queue. To exit press CTRL+C")
	<-forever
}
