# RFP Site

## Prerequisites

- node.js/npm
- Golang
- PostgreSQL

## Running the frontend

- `cd` into the client directory and run `npm i` to install the node dependencies
- `npm run dev` to run the dependencies

## Running the Server
- Create a postgreSQL database 
- Create a `dbVars.go` file with the following content that corresponds to the database that you just created

```go
package dbutils

const (
	host     = "localhost"
	port     = 5432
	username = "postgres"
	password = "your-password"
	dbname   = "rfp_db"
)
```

- `cd` into the migrations directory 

- `cd` into the `server/src` directory
- `go run main.go` to run the server