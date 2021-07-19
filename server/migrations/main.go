package main

import (
	"fmt"
	"jacobsmi/server/src/dbutils"
)

func main() {
	defer dbutils.DB.Close()

	sqlStatement := ` CREATE TABLE known_qa(
		id SERIAL,
		question VARCHAR UNIQUE,
		answer VARCHAR,
		comments VARCHAR
		);`
	_, err := dbutils.DB.Exec(sqlStatement)
	if err != nil {
		fmt.Println("Error creating database")
		fmt.Println(err)
	}
}
