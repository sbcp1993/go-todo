FROM ubuntu:18.04

WORKDIR /todo

COPY go-todo /todo/

EXPOSE 8080

CMD ["./go-todo"]