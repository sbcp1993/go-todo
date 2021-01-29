package db

import (
	"context"

	"go-todo/types"
)

const backendkey = "backendDB"

type BackendDB interface {
	Close()
	AddUser(usr *types.User) error
	Login(name string, pwd string) (res map[string]interface{}, err error)
	Insert(todo *types.Todo) (string, error)
	Delete(uuid string) error
	Update(uuid string, todo *types.Todo) error
	SetComplete(uuid string, status bool) error
	GetTodoByStatus(status bool) (titles []string, err error)
	GetTodoByPriority(priority string) (titles []string, err error)
}

func SetBackend(ctx context.Context, backend BackendDB) context.Context {
	return context.WithValue(ctx, backendkey, backend)
}

func Close(ctx context.Context) {
	getBackendDB(ctx).Close()
}
func AddUser(ctx context.Context, usr *types.User) error {
	return getBackendDB(ctx).AddUser(usr)
}
func Login(ctx context.Context, name string, pwd string) (map[string]interface{}, error) {
	return getBackendDB(ctx).Login(name, pwd)
}
func Insert(ctx context.Context, todo *types.Todo) (string, error) {
	return getBackendDB(ctx).Insert(todo)
}

func Delete(ctx context.Context, uuid string) error {
	return getBackendDB(ctx).Delete(uuid)
}

func Update(ctx context.Context, uuid string, todo *types.Todo) error {
	return getBackendDB(ctx).Update(uuid, todo)
}

func SetComplete(ctx context.Context, uuid string, status bool) error {
	return getBackendDB(ctx).SetComplete(uuid, status)
}

func GetTodoByStatus(ctx context.Context, status bool) (titles []string, err error) {
	return getBackendDB(ctx).GetTodoByStatus(status)
}

func GetTodoByPriority(ctx context.Context, priority string) (titles []string, err error) {
	return getBackendDB(ctx).GetTodoByPriority(priority)
}

func getBackendDB(ctx context.Context) BackendDB {
	return ctx.Value(backendkey).(BackendDB)
}
