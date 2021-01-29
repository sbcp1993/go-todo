package handler

import (
	"net/http"

	"go-todo/db"
)

func SetUpConnection(conn *db.DBConnection) *http.ServeMux {
	todoHandler := &todoHandler{
		postgres: conn,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/register", todoHandler.addUser)
	mux.HandleFunc("/login", todoHandler.getToken)
	mux.HandleFunc("/markcomplete", todoHandler.completeTodo)
	mux.HandleFunc("/getbycompletestatus", todoHandler.getTodoByStatus)
	mux.HandleFunc("/getbypriority", todoHandler.getTodoByPriority)
	mux.HandleFunc("/create", todoHandler.createTodo)
	mux.HandleFunc("/update", todoHandler.updateTodo)
	mux.HandleFunc("/delete", todoHandler.deleteTodo)
	mux.HandleFunc("/logout", todoHandler.logout)

	return mux
}
