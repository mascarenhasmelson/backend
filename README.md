# GoLang Backend Simple Banking Ledger (In development stage)

## _Introduction_ 
This document provides details about a simple banking ledger backend built using Golang and its libraries. Additionally, it includes a setup guide for running the Go backend server with Docker instructions

**Note:** APIs are case-sensitive, so be sure to use the correct endpoint format.
Additionally, before running the backend server users must set up Docker Swarm to initialize RabbitMQ and MongoDB databases. All code and Docker scripts have been tested and compiled on a Linux system.

Adjust the **docker-compose.yml**  CPU, replicas, and memory settings as needed to suit your system's requirements, these values were configured for my setup

- Postman is used to test APIs and endpoints for POST and GET operations.
- SQLite serves as the relational database for handling withdrawals and deposits.
- MongoDB is utilized for logging transactional data.
- RabbitMQ facilitates the producer-consumer queue system as a message broker.
- Docker Swarm enables horizontal scaling.

### Backend package information:

- ***Consumer Package***
    - This package contains the core logic for handling the queue system, including both consumer and producer operations for RabbitMQ.
- ***DB Package***
    - Responsible for database interactions, this package manages SQL operations and handles the connection to the SQLite database.
- ***Error Package***
    - Designed to showcase error-handling mechanisms, this package is intended to   centralize error messages and maintain clean code. While not fully implemented integrated for code, 
- ***Models Package***
    - Contains struct definitions that abstract application entities like accounts,    transactions, and ledgers.
- ***MongoDB Package***
    - Houses the logic for MongoDB interactions, including database connections and transaction-related functions.

#### **Below are the command to setup docker swarm for scaling**

#### Docker Swarm Initialization and Deployment

#### 1. Initialize Docker Swarm
To initialize a Docker Swarm, run the following command:
```
docker swarm init
```
**Note:** Run this command only if your system's network card interface has more than one IP address.

#### 2. Initialize with Specific Advertised IP
If you need to specify a particular IP address to advertise for the swarm, use the following command:
```
docker swarm init --advertise-addr <network-card-ip>
```
#### 3. Get Join Token for Manager Node
```
docker swarm join-token manager
```
#### 4. Create an Overlay Network
```
docker network create -d overlay backend
```
#### 5. Deploy the Stack Using Docker Compose
```
docker stack deploy -c docker-compose.yml backend_stack
```
> **Note:** It may take a few seconds or minutes for Docker Swarm to deploy all the 
> containers. Please be patient while the stack is being initialized and containers are brought up.

_To run the program directly without building an executable file, use the following command:_

``` go run main.go ```

_The Go backend will start running on_ **port 8080**. _Ensure that no other service is already using this port._

_To build the package, use the command:_ **go build**

This will generate an executable file named **backend**. To run the generated build file, use the following command: 

``` sudo ./backend ```

Ensure you have the necessary permissions to execute the file.
1) **Support the creation of accounts with specified initial balances.**  
Initial account can be created by using this endpoint
```   
 http://localhost:8080/createAccount
{
 "customerID": "123456",
 "name": "melson",
 "amount": 5000
}
```
2) **Facilitate deposits and withdrawals of funds.**  
The **Transaction** endpoint uses a boolean-based parameter to differentiate between deposit and withdrawal operations.
    - To **deposit**, set IsDeposit **to** true.
    - To **withdraw**, set IsDeposit **to** false.
    
    **For Deposit API**
    ```
    http://localhost:8080/Transaction
    {
    "customerID":"123456",
    "amount":500.00,
    "isDeposit":true
     }
     ```

   **For Withdraw API**
   ```
    http://localhost:8080/Transaction
    {
    "customerID":"123456",
    "amount":500.00,
    "isDeposit":flase
     }
     ```
3) **Maintain a detailed transaction log (ledger) for each account**  
For ledger purpose mongodb NoSQL database has been chosen  
API for pulling ledger, its based on customer ID
   
   ```
   http://localhost:8080/getCustomerData
    {
    "customerID": "123456"
    }
    ```
**_If the API is unable to provide accurate values, you can delve deeper into the detailed information, which can be retrieved directly from the container. This issue arises due to problems with JSON marshalling and unmarshalling._**  \
Command to pull ledger from mongodb container

``` sudo docker exec -it <mongo-container-id> mongosh ```  
``` show collection ```  
``` test> db.ledger.find().pretty() ```

> When the backend server is executed, the test.db SQLite database is automatically created. To view the data in **test.db**, users can either download a database browser tool from [SQLite Browser](https://sqlitebrowser.org/) or upload the database file to the [SQLite Viewer](https://sqliteviewer.app/) website for easy access.

4) **Ensure ACID-like consistency, Scale horizontally, Integrate an asynchronous queue or broker, Include a comprehensive testing strategy**

Partial testing has been conducted, with minimal validation implemented, the system functions correctly when valid values are provided. RabbitMQ has been integrated to handle deposit and withdrawal operations through a queue system. Horizontal scaling is achieved using Docker Swarm. 

