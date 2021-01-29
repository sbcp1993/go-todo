package handler

import (
	"bytes"
	"fmt"
	"go-todo/db"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_todoHandler_createTodo(t *testing.T) {
	dbc, err := db.NewDBConnection()
	fmt.Println(dbc, err)

	req, err := http.NewRequest("POST", "http://localhost:8080/create", bytes.NewBuffer([]byte(`{"title":"task1","description":"new task","due_date":"2020-01-01T00:00:00Z","priority":"high"}`)))
	fmt.Println(req, err)
	req.Header.Add("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJOYW1lIjoic3J1dGhpIiwiZXhwIjoxNjE3OTE3MzAyfQ.8yhQtZDsigo5qD2LhBqLJFK1KRfrydMqCpsRHbPgzL")
	rw := httptest.NewRecorder()
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name    string
		handler *todoHandler
		args    args
	}{
		{
			name: "test",
			handler: &todoHandler{
				postgres: dbc,
			},
			args: args{
				w: rw,
				r: req,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.handler.createTodo(tt.args.w, tt.args.r)
		})
		fmt.Println(rw)
	}
}

func Test_todoHandler_updateTodo(t *testing.T) {
	dbc, err := db.NewDBConnection()
	fmt.Println(dbc, err)

	req, err := http.NewRequest("PATCH", "http://localhost:8080/update", bytes.NewBuffer([]byte(`{"id": "2b825172-553d-4e8d-a442-14f665dc843d", "todo": {"title":"","description":"","due_date":"2001-01-01T00:00:00Z","priority":"low"}}`)))
	fmt.Println(req, err)

	rw := httptest.NewRecorder()
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name    string
		handler *todoHandler
		args    args
	}{
		{
			name: "test",
			handler: &todoHandler{
				postgres: dbc,
			},
			args: args{
				w: rw,
				r: req,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.handler.updateTodo(tt.args.w, tt.args.r)
		})
	}
}

func Test_todoHandler_completeTodo(t *testing.T) {
	dbc, err := db.NewDBConnection()
	fmt.Println(dbc, err)

	req, err := http.NewRequest("PATCH", "http://localhost:8080/markcomplete", bytes.NewBuffer([]byte(`{"id": "7c6c644f-4105-4cea-aed1-9b4e3f4d9a3a", "status": true}`)))
	fmt.Println(req, err)

	rw := httptest.NewRecorder()
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name    string
		handler *todoHandler
		args    args
	}{
		{
			name: "test",
			handler: &todoHandler{
				postgres: dbc,
			},
			args: args{
				w: rw,
				r: req,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.handler.completeTodo(tt.args.w, tt.args.r)
		})
	}
}

func Test_todoHandler_getTodoByPriority(t *testing.T) {
	dbc, err := db.NewDBConnection()
	fmt.Println(dbc, err)

	req, err := http.NewRequest("GET", "http://localhost:8080/getbypriority", bytes.NewBuffer([]byte(`{"priority": "low"}`)))
	fmt.Println(req, err)

	rw := httptest.NewRecorder()
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name    string
		handler *todoHandler
		args    args
	}{
		{
			name: "test",
			handler: &todoHandler{
				postgres: dbc,
			},
			args: args{
				w: rw,
				r: req,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.handler.getTodoByPriority(tt.args.w, tt.args.r)
		})
		fmt.Println(tt.args.w.Header())
	}
}

func Test_todoHandler_deleteTodo(t *testing.T) {
	dbc, err := db.NewDBConnection()
	fmt.Println(dbc, err)

	req, err := http.NewRequest("DELETE", "http://localhost:8080/delete", bytes.NewBuffer([]byte(`{"id": "7c6c644f-4105-4cea-aed1-9b4e3f4d9a3a"}`)))
	fmt.Println(req, err)

	rw := httptest.NewRecorder()
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name    string
		handler *todoHandler
		args    args
	}{
		{
			name: "test",
			handler: &todoHandler{
				postgres: dbc,
			},
			args: args{
				w: rw,
				r: req,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.handler.deleteTodo(tt.args.w, tt.args.r)
		})
	}
}

func Test_todoHandler_getToken(t *testing.T) {
	dbc, err := db.NewDBConnection()
	fmt.Println(dbc, err)

	req, err := http.NewRequest("DELETE", "http://localhost:8080/login", bytes.NewBuffer([]byte(`{"username":"sruthi","passwordhash":"sruthis"}`)))
	fmt.Println(req, err)

	rw := httptest.NewRecorder()
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name    string
		handler *todoHandler
		args    args
	}{
		{
			name: "test",
			handler: &todoHandler{
				postgres: dbc,
			},
			args: args{
				w: rw,
				r: req,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.handler.getToken(tt.args.w, tt.args.r)
		})
	}
}
