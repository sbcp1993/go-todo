package db

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"go-todo/types"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	sqlbuilder "github.com/huandu/go-sqlbuilder" // SQL builder
	_ "github.com/lib/pq"
)

const (
	errUnknownDBType = "unknown DBType %s"
	dbtodoTable      = "goapp.todolist"
	dbuserTable      = "goapp.auth"
)

// DBConnection is a holder for a database connection
type DBConnection struct {
	dbType  string
	DB      *sql.DB
	Builder sqlbuilder.Flavor
	IDMax   int
}

// NewDBConnection connects to the DB
func NewDBConnection() (conn *DBConnection, err error) {
	conn = new(DBConnection)
	conn.dbType = "postgres"

	conn.DB, err = sql.Open("postgres", "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable")
	if err != nil {
		return
	}
	err = conn.DB.Ping()
	if err != nil {
		return
	}
	conn.IDMax = 9223372036854775807

	conn.Builder = sqlbuilder.PostgreSQL
	return
}

func (conn *DBConnection) Close() {
	conn.DB.Close()
}

func (conn *DBConnection) Insert(todo *types.Todo) (string, error) {
	var cols []string
	var vals []interface{}

	if todo.Title == "" {
		return "", errors.New("Title cannot be empty")
	}
	uuid := uuid.New().String()
	cols = append(cols, "id", "title")
	vals = append(vals, uuid, todo.Title)
	if todo.Description != "" {
		cols = append(cols, "todo_description")
		vals = append(vals, todo.Description)
	}
	if !todo.DueDate.IsZero() {
		cols = append(cols, "due_date")
		vals = append(vals, todo.DueDate)
	}
	if todo.Priority != "" {
		cols = append(cols, "todo_priority")
		vals = append(vals, todo.Priority)
	}
	ib := conn.Builder.NewInsertBuilder()
	ib.InsertInto(dbtodoTable).Cols(cols...).Values(vals...)

	sqlString, args := ib.Build()

	_, err := conn.DB.Query(sqlString, args...)
	if err != nil {
		return "", err
	}

	return uuid, nil
}

func (conn *DBConnection) Delete(uuid string) error {
	delb := conn.Builder.NewDeleteBuilder()
	delb.DeleteFrom(dbtodoTable).Where(delb.Equal("id", uuid))
	sqlString, args := delb.Build()

	if _, err := conn.DB.Exec(sqlString, args...); err != nil {
		return err
	}

	return nil
}

func (conn *DBConnection) Update(uuid string, todo *types.Todo) error {
	var setAssigns []string
	ub := conn.Builder.NewUpdateBuilder()
	ub.Update(dbtodoTable)
	if todo.Title != "" {
		setAssigns = append(setAssigns, ub.Assign("title", todo.Title))
	}
	if todo.Description != "" {
		setAssigns = append(setAssigns, ub.Assign("todo_description", todo.Description))
	}
	if !todo.DueDate.IsZero() {
		setAssigns = append(setAssigns, ub.Assign("due_date", todo.DueDate))
	}
	if todo.Priority != "" {
		setAssigns = append(setAssigns, ub.Assign("todo_priority", todo.Priority))
	}
	ub.Set(setAssigns...).Where(ub.Equal("id", uuid))

	sqlString, args := ub.Build()

	if _, err := conn.DB.Exec(sqlString, args...); err != nil {
		return err
	}
	return nil
}

func (conn *DBConnection) SetComplete(uuid string, status bool) error {
	ub := conn.Builder.NewUpdateBuilder()
	ub.Update(dbtodoTable).Set(
		ub.Assign("iscomplete", status),
	).Where(ub.Equal("id", uuid))

	sqlString, args := ub.Build()

	if _, err := conn.DB.Exec(sqlString, args...); err != nil {
		return err
	}
	return nil
}

func (conn *DBConnection) GetTodoByStatus(status bool) (titles []string, err error) {
	sb := conn.Builder.NewSelectBuilder()
	sb.Select("title").From(dbtodoTable).Where(sb.Equal("iscomplete", status))

	sqlString, args := sb.Build()

	rows, err := conn.DB.Query(sqlString, args...)
	if err != nil {
		return titles, err
	}

	var t string
	for rows.Next() {
		if err = rows.Scan(&t); err != nil {
			return titles, err
		}
		titles = append(titles, t)
	}

	return titles, nil
}

func (conn *DBConnection) GetTodoByPriority(priority string) (titles []string, err error) {
	sb := conn.Builder.NewSelectBuilder()
	sb.Select("title").From(dbtodoTable).Where(sb.Equal("todo_priority", priority))

	sqlString, args := sb.Build()

	rows, err := conn.DB.Query(sqlString, args...)
	if err != nil {
		return titles, err
	}

	var t string
	for rows.Next() {
		if err = rows.Scan(&t); err != nil {
			return titles, err
		}
		titles = append(titles, t)
	}

	return titles, nil
}

func buildPasswordHash(pwd string) (hashB64 string, err error) {
	var bytePWD []byte
	var hash []byte

	if "" == pwd {
		err = errors.New("password empty")
		return
	}

	if hash, err = bcrypt.GenerateFromPassword(bytePWD, int(14)); err != nil {
		return
	}

	hashB64 = base64.StdEncoding.EncodeToString(hash)
	return
}

func (conn *DBConnection) userExists(name string) (exists bool, err error) {
	var sqlErr error
	var sb *sqlbuilder.SelectBuilder
	var sqlString string
	var args []interface{}
	var cnt int

	sb = conn.Builder.NewSelectBuilder()
	sb.Select("count(*)")
	sb.From(dbuserTable)
	sb.Where(sb.Equal("username", name))

	sqlString, args = sb.Build()

	sqlErr = conn.DB.QueryRow(sqlString, args...).Scan(&cnt)
	if sqlErr != nil {
		return
	}

	exists = 1 == cnt

	return
}

func (conn *DBConnection) AddUser(usr *types.User) error {
	var hash string
	var err error
	var exist bool

	if exist, err = conn.userExists(usr.Name); err != nil {
		return err
	}
	if !exist {
		if hash, err = buildPasswordHash(usr.PasswordHash); err != nil {
			return err
		}
		ib := conn.Builder.NewInsertBuilder()
		ib.InsertInto(dbuserTable).Cols("username", "password_hash").Values(usr.Name, hash)

		sqlString, args := ib.Build()

		_, err := conn.DB.Query(sqlString, args...)
		if err != nil {
			return err
		}

	} else {
		err = errors.New("user exists")
		return err
	}

	return nil
}

func (conn *DBConnection) Login(name, password string) (res map[string]interface{}, err error) {
	var passwordhash string
	var exist bool
	var tokenString string

	if exist, err = conn.userExists(name); err != nil {
		return
	}
	if exist {
		sb := conn.Builder.NewSelectBuilder()
		sb = conn.Builder.NewSelectBuilder()
		sb.Select("password_hash")
		sb.From(dbuserTable)
		sb.Where(sb.Equal("username", name))

		sqlString, args := sb.Build()

		err = conn.DB.QueryRow(sqlString, args...).Scan(&passwordhash)
		if err != nil {
			return
		}

		expiresAt := time.Now().Add(time.Minute * 100000).Unix()

		err = bcrypt.CompareHashAndPassword([]byte(passwordhash), []byte(password))
		if err != nil {
			res = map[string]interface{}{"status": false, "message": "Invalid login credentials. Please try again"}
			return
		}

		tk := &types.Token{
			Name: name,
			StandardClaims: &jwt.StandardClaims{
				ExpiresAt: expiresAt,
			},
		}

		os.Setenv("ACCESS_SECRET", "secret")

		token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

		tokenString, err = token.SignedString([]byte(os.Getenv("ACCESS_SECRET")))

		if err != nil {
			return
		}

		res = map[string]interface{}{"status": false, "message": "logged in"}
		res["token"] = tokenString
		res["user"] = name
		return

	} else {
		err = errors.New("user does not exists")
		return
	}
}
