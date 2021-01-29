package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const welcomemsg = `
Organize your tasks with todo:
Options:
  register 		user registration
  login			log in to todo
  add     		create a new todo
  update  		Update existing todo
  delete  		Delete a todo
  mark    		Complete a todo
  high			Get high priority tasks
  low			Get low priority tasks
  complete 		Get completed tasks
  uncomplete 	Get uncompleted tasks
`

func main() {
	if len(os.Args) < 2 {
		fmt.Println(welcomemsg)
		return
	}

	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "register":
			register()
		case "login":
			login()
		case "add":
			add()
		case "update":
			alter()
		case "delete":
			del()
		case "mark":
			mark()
		case "high":
			getbypriority("high")
		case "low":
			getbypriority("low")
		case "complete":
			getbystatus(true)
		case "uncomplete":
			getbystatus(false)
		case "logout":
			logout()
		default:
			fmt.Printf("'%s' is not a todo command.", command)
		}
	}
}

func logout() {
	var b []byte
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/logout", bytes.NewReader(b))
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	fmt.Println(res.StatusCode)
}

func register() {
	if len(os.Args) < 4 {
		fmt.Println("cmd syntax: register \"username\" \"password\"")
		return
	}
	var name string
	var passwd string

	if len(os.Args) == 4 {
		name = os.Args[2]
		passwd = os.Args[3]
	}

	user := struct {
		Name         string `json:"username"`
		PasswordHash string `json:"passwordhash"`
	}{
		name, passwd,
	}

	b, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	res, err := http.Post("http://localhost:8080/register", "application/json", bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	fmt.Println(res.StatusCode)
	defer res.Body.Close()
}

func login() {
	if len(os.Args) < 4 {
		fmt.Println("cmd syntax: login \"username\" \"password\"")
		return
	}
	var name string
	var passwd string

	if len(os.Args) == 4 {
		name = os.Args[2]
		passwd = os.Args[3]
	}

	user := struct {
		Name         string `json:"username"`
		PasswordHash string `json:"passwordhash"`
	}{
		name, passwd,
	}

	b, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	res, err := http.Post("http://localhost:8080/login", "application/json", bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.Header)
	fmt.Println("token: ", string(b))
}

func getbypriority(p string) {
	var cmdsyntax = fmt.Sprintf("cmd syntax: %s \"token\"", p)
	if len(os.Args) < 3 {
		fmt.Println(cmdsyntax)
		return
	}
	var token string

	if len(os.Args) == 3 {
		token = os.Args[2]
	}
	req := struct {
		Priority string `json:"priority"`
	}{
		p,
	}
	b, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	request, err := http.NewRequest(http.MethodGet, "http://localhost:8080/getbypriority", bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	request.Header.Add("Authorization", token)
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	if len(b) == 0 {
		fmt.Println("No Data")
		return
	}
	fmt.Println("task list: ", string(b))
}

func getbystatus(s bool) {
	var cmdsyntax string
	if s {
		cmdsyntax = fmt.Sprintf("cmd syntax: complete \"token\"")
	} else {
		cmdsyntax = fmt.Sprintf("cmd syntax: uncomplete \"token\"")
	}

	if len(os.Args) < 3 {
		fmt.Println(cmdsyntax)
		return
	}
	var token string

	if len(os.Args) == 3 {
		token = os.Args[2]
	}
	req := struct {
		Status bool `json:"status"`
	}{
		s,
	}
	b, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	request, err := http.NewRequest(http.MethodGet, "http://localhost:8080/getbycompletestatus", bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	request.Header.Add("Authorization", token)
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	if len(b) == 0 {
		fmt.Println("No Data")
		return
	}
	fmt.Println("task list: ", string(b))
}

func add() {
	if len(os.Args) < 7 {
		fmt.Println("cmd syntax: add \"title\" \"description\" \"priority\" \"date (mm/dd/yy)\" \"token\"")
		return
	}
	var title string
	var desc string
	var priority string
	var due time.Time
	var err error
	var token string

	if len(os.Args) == 7 {
		title = os.Args[2]
		desc = os.Args[3]
		priority = os.Args[4]

		if due, err = time.Parse("01/02/06", os.Args[5]); err != nil {
			panic(err)
		}

		token = os.Args[6]
	}

	todo := struct {
		Title       string    `json:"title"`
		Description string    `json:"description,omitempty"`
		DueDate     time.Time `json:"due_date,omitempty"`
		Priority    string    `json:"priority"`
	}{
		title, desc, due, priority,
	}

	b, err := json.Marshal(todo)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/create", bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	request.Header.Add("Authorization", token)
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("status code : %d \n uuid: %s", res.StatusCode, string(b))
}

func del() {
	if len(os.Args) < 4 {
		fmt.Println("cmd syntax: delete \"uuid\" \"token\"")
		return
	}
	var id string
	var token string
	if len(os.Args) == 4 {
		id = os.Args[2]
		token = os.Args[3]
	}
	// fmt.Println("Task id:")
	// fmt.Scanln(&id)

	req := struct {
		UUID string `json:"id"`
	}{
		id,
	}
	b, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/delete", bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	request.Header.Add("Authorization", token)
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	if res.StatusCode == http.StatusOK {
		fmt.Println("Deleted successfully")
	} else {
		fmt.Println("Deletion Failed...!!")
	}
	res.Body.Close()
}

func alter() {
	var cmdsyntax = `
	cmd syntax: update "id" "title" "description" "due_date" "priority" "token"
	provide empty string for values that not needed to update
	`
	if len(os.Args) < 8 {
		fmt.Println(cmdsyntax)
		return
	}
	var id string
	var title string
	var desc string
	var priority string
	var due time.Time
	var err error
	var token string

	if len(os.Args) == 8 {
		id = os.Args[2]
		title = os.Args[3]
		desc = os.Args[4]

		if os.Args[5] != "" {
			if due, err = time.Parse("01/02/06", os.Args[5]); err != nil {
				panic(err)
			}
		}

		priority = os.Args[6]
		token = os.Args[7]

	}

	type todo struct {
		Title       string    `json:"title,omitempty"`
		Description string    `json:"description,omitempty"`
		DueDate     time.Time `json:"due_date,omitempty"`
		Priority    string    `json:"priority,omitempty"`
	}
	req := struct {
		UUID string `json:"id"`
		Todo todo   `json:"todo,omitempty"`
	}{
		id, todo{title, desc, due, priority},
	}
	b, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	request, err := http.NewRequest(http.MethodPatch, "http://localhost:8080/update", bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	request.Header.Add("Authorization", token)
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	if res.Status == http.StatusText(http.StatusOK) {
		fmt.Println("Updated successfully")
	}
	res.Body.Close()
}

func mark() {
	var cmdsyntax = `
	cmd syntax: mark "id" "token"
	`
	if len(os.Args) < 4 {
		fmt.Println(cmdsyntax)
		return
	}
	var id string
	var token string

	if len(os.Args) == 4 {
		id = os.Args[2]
		token = os.Args[3]

	}
	// var id string

	// fmt.Println("Enter completed task's id:")
	// fmt.Scanln(&id)

	req := struct {
		UUID   string `json:"id"`
		Status bool   `json:"status"`
	}{
		id, true,
	}
	b, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest(http.MethodPatch, "http://localhost:8080/markcomplete", bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	request.Header.Add("Authorization", token)
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	if res.Status == http.StatusText(http.StatusOK) {
		fmt.Println("Updated successfully")
	}
	res.Body.Close()
}
