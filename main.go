package main

import (
	// "innoscripta/db"
	// "innoscripta/mongodb"
	"encoding/json"
	// "go.mongodb.org/mongo-driver/bson"
    "log"
	"fmt"
	"time"
	"context"
	"gitlab.com/melsonmascarenhas/backend/db"
	"net/http"
	"gitlab.com/melsonmascarenhas/backend/models"
	"gitlab.com/melsonmascarenhas/backend/consumer"
	"gitlab.com/melsonmascarenhas/backend/mongodb"
	 amqp "github.com/rabbitmq/amqp091-go"
	//  "github.com/gorilla/mux"
	//  	_ "modernc.org/sqlite"
	// "github.com/jmoiron/sqlx"
)

var(
	 account models.PubAccount
     withdraw models.Transaction
	  ledgers []models.Account_ledger
)

type RequestBody struct {
	CustomerID string `json:"customerID"`
}
type CustomerData struct {
	Account models.Account         `json:"account"`
	Ledgers []models.Account_ledger `json:"ledgers"`
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// https://www.rabbitmq.com/tutorials/tutorial-one-go


func CreatePublisherHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure only POST requests are allowed
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&account); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// // Connect to RabbitMQ
	conn, err:= amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err:= conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare a queue
	q, err := ch.QueueDeclare(
		"createAccount", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Serialize the account data to JSON
	body, err := json.Marshal(account)
	if err != nil {
		http.Error(w, "Failed to serialize data", http.StatusInternalServerError)
		return
	}

	// Create a context for publishing with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Publish the message to the RabbitMQ queue
	err = ch.PublishWithContext(ctx,
		"",       // exchange
		q.Name,   // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "application/json", // Set the content type
			Body:        body,
		},
	)
	if err != nil {
		http.Error(w, "Failed to publish message", http.StatusInternalServerError)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message published successfully"))
	fmt.Printf("Published message: %s\n", string(body))
}
//Withdraw publisher
func WithdrawPublisherHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure only POST requests are allowed
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&withdraw); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	type RequestBody struct {
		CustomerID string `json:"customerID"`
	}
	// Connect to RabbitMQ
	conn, err:= amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err:= conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare a queue
	q, err := ch.QueueDeclare(
		"withdraw", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Serialize the withdraw data to JSON
	body, err := json.Marshal(withdraw)
	if err != nil {
		http.Error(w, "Failed to serialize data", http.StatusInternalServerError)
		return
	}

	// Create a context for publishing with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Publish the message to the RabbitMQ queue
	err = ch.PublishWithContext(ctx,
		"",       // exchange
		q.Name,   // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "application/json", // Set the content type
			Body:        body,
		},
	)
	if err != nil {
		http.Error(w, "Failed to publish message", http.StatusInternalServerError)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message published successfully"))
	fmt.Printf("Published message: %s\n", string(body))
}

func getLedgerHandler(w http.ResponseWriter, r *http.Request) {
    // Ensure the request method is POST
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    // Parse the incoming JSON body
    var requestData struct {
        CustomerID string `json:"customerID"`
    }

    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&requestData); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    // Call the LedgerByCustomerID function with the customerID from the request
    ledgers, err := mongodb.LedgerByCustomerID(requestData.CustomerID)
    if err != nil {
        http.Error(w, "Failed to fetch ledger data", http.StatusInternalServerError)
        return
    }

    // Set response header to JSON
    w.Header().Set("Content-Type", "application/json")

    // Send the ledger data as a JSON response
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(ledgers); err != nil {
        log.Printf("Failed to encode response: %v", err)
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
    }
}

func main(){

log.Println("Initializing database...")
	db.Create()

	//From consumer package
	log.Println("Starting consumers...")
	go consumer.StartCreateAccountConsumer()
	go consumer.StartWithdrawConsumer()

	// Routes
	http.HandleFunc("/createAccount", CreatePublisherHandler)
	http.HandleFunc("/Transaction", WithdrawPublisherHandler)
	http.HandleFunc("/getCustomerData", getLedgerHandler)

	// Server 
	go func() {
		log.Println("HTTP server is running on http://localhost:8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Prevent the main function from exiting
	select {}
}


