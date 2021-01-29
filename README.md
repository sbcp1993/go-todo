# go-todo
Implementing a to-do list api in golang

clone this repository to your GOPATH

cd to go-todo and run 'docker-compose up -d' to up the service

a sample cli for testing this api is developed in go-todo/cmd/main.go
To check this, cd to go-todo/cmd and run 'go run main.go' this will print the following

```
Organize your tasks with todo:
Options:
  add                   create a new todo
  update                Update existing todo
  delete                Delete a todo
  mark                  Complete a todo
  high                  Get high priority tasks
  low                   Get low priority tasks
  complete              Get completed tasks
  uncomplete            Get uncompleted tasks

  
Register a user with 'go run main.go register <username> <password>' 
Login to the api with registered user credentials 'go run main.go login <username> <password>'
```

Use command options with 'go run main.go' command to create,update,delete,select ...



