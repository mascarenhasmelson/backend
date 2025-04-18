package mongodb
import (
	"context"
	"log"
	"fmt"
	"time"
	"gitlab.com/melsonmascarenhas/backend/models"
	"gitlab.com/melsonmascarenhas/backend/db"
	"github.com/google/uuid" 
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "modernc.org/sqlite"
	
)
var (
	mongoDB *mongo.Database
	mongoClient *mongo.Client
	account models.Account
	ledgerAccount []models.Account_ledger
	ledgers []models.Account_ledger
	getdata []models.Account_ledger
)

func InsertAccount_mongo(account models.Account_ledger) error {
	initMongoDB() 
	collection := GetCollection("ledger")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	prin, err := collection.InsertOne(ctx, account)
	if err != nil {
		return fmt.Errorf("failed to insert account: %v", err)
	}

	log.Println("Account inserted successfully", prin)
	return nil
}

//MongoDB conection
func initMongoDB() {

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB")
	mongoClient = client
}
func GetCollection(collectionName string) *mongo.Collection {
	if mongoClient == nil {
		initMongoDB()
	}
	return mongoClient.Database("test").Collection(collectionName)
}

//init create account
func CreateAccount(name string,customer_ID string, initialBalance float64) (string, error) {
	account := models.Account{
		ID:         uuid.NewString(),
		CustomerID: customer_ID,
		Name:       name,
		Balance:    initialBalance,
		CreatedAt:  time.Now(),
	}
	accountMongo := models.Account_ledger{
		ID:         uuid.NewString(),
		Balance:    initialBalance,
		Type:       "deposit", 
		CustomerID: customer_ID,
		Name:       name,
		Withdraw_balance: 0,
		Deposit_balance:initialBalance,
		CreatedAt:  time.Now(),
	}
fmt.Println(account)

if initialBalance <=0 && initialBalance <500{
fmt.Println("Deposit amount more than 500")
}
db.Insert(account)
InsertAccount_mongo(accountMongo)
return account.ID, nil
}
func WithDraw(CustomerID string,amount float64,isDeposit bool)error{
	// Connect to the SQLite database
	dbb, _ := db.ConnectDB()
	// Compare the balance and amount
	query := `SELECT Balance FROM account WHERE CustomerID = ?`
 dbb.Get(&account.Balance, query, CustomerID)

log.Printf("CustomerID: %s,withdraw Balance: %.2f,amount %2.f\n", CustomerID, amount,account.Balance)
	//  Update the balance
	var transactionType string
	var Deposit_balance float64
	var Withdraw_balance float64
	if amount < 0 {
        return fmt.Errorf("amount cannot be negative")
    } else if isDeposit {
        account.Balance += amount
		Deposit_balance=amount
		transactionType = "deposit"
        fmt.Println("I'm in positive") //Tracing 
    } else {
        if account.Balance < amount {
            return fmt.Errorf("insufficient funds")
        }
        account.Balance -= amount
		Withdraw_balance=amount
		transactionType = "withdrawal"
        fmt.Println("I'm in negative") //Tracing 
    }
	updateStatement := `UPDATE account SET Balance = ? WHERE CustomerID = ?`
	_, err := dbb.Exec(updateStatement, account.Balance, CustomerID)
	if err != nil {
		 fmt.Errorf("failed to update balance in SQLite: %v", err)
		
	}
//Log printing data for tracing
	afterupdate := `SELECT Balance FROM account WHERE CustomerID = ?`
	dbb.Get(&amount, afterupdate, CustomerID)
   log.Printf("CustomerID: %s, Balance: %.2f\n", CustomerID, account.Balance)


	//For ledger purpose
	mongoDB := GetCollection("ledger")
	ledgerEntry := bson.M{
		"CustomerID": CustomerID,
		"Balance":     account.Balance,
		"Type":        transactionType, // E.g., "deposit" or "withdraw"
		"Amount":      amount,
		"Deposit_balance":Deposit_balance,
		"Withdraw_balance":Withdraw_balance,
		"CreatedAt":   time.Now(),
	}
	
	// entering the ledger entry into the MongoDB collection
	_, err = mongoDB.InsertOne(context.TODO(), ledgerEntry)
	if err != nil {
		log.Printf("Failed to insert ledger entry in MongoDB: %v", err)
	} else {
		log.Printf("Ledger entry added for CustomerID: %s, Type: %s, Amount: %.2f, Balance: %.2f\n",
			CustomerID, transactionType, amount, account.Balance)
	}
	return nil
}


func LedgerByCustomerID(customerID string) ([]models.Account_ledger, error) {
	
		// time for excution for each api hit 
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
	
		// Connect to MongoDB
		client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
		if err != nil {
			log.Printf("Failed to connect to MongoDB: %v", err)
			return nil, err
		}
		defer client.Disconnect(ctx)
	
		
		collection := client.Database("test").Collection("ledger")
	
		
		filter := bson.M{"customer_id": customerID}
		log.Printf("Filter being used: %v", filter)  // Log
	
		
		cursor, err := collection.Find(ctx, filter)
		if err != nil {
			log.Printf("Failed to fetch ledger for CustomerID %s: %v", customerID, err)
			return nil, err
		}
		defer cursor.Close(ctx)
	
	
		var ledgers []models.Account_ledger
		for cursor.Next(ctx) {
			var ledger models.Account_ledger
			if err := cursor.Decode(&ledger)
			 err != nil {
				log.Printf("Failed to decode ledger document: %v", err)
				continue
			}
			ledgers = append(ledgers, ledger)
		}
	
		// Check for cursor errors
		if err := cursor.Err(); err != nil {
			log.Printf("Cursor error while fetching ledgers for CustomerID %s: %v", customerID, err)
			return nil, err
		}
	
		return ledgers, nil
	}
