package main

import (
	"fmt"
	"log"
	"net/http"

	"go-todo/db"
	"go-todo/handler"
)

func main() {
	var postgres *db.DBConnection
	var err error

	postgres, err = db.NewDBConnection()

	if err != nil {
		panic(err)
	} else if postgres == nil {
		panic("postgres is nil")
	}

	mux := handler.SetUpConnection(postgres)

	fmt.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
