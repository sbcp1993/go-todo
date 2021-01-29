package runtime

import (
	"context"

	"go-todo/db"
	"go-todo/types"
)

func Close(ctx context.Context) {
	db.Close(ctx)
}

func InsertUser(ctx context.Context, user *types.User) error {
	return db.AddUser(ctx, user)
}
func Login(ctx context.Context, name string, pwd string) (map[string]interface{}, error) {
	return db.Login(ctx, name, pwd)
}
func Insert(ctx context.Context, todo *types.Todo) (string, error) {
	return db.Insert(ctx, todo)
}

func Delete(ctx context.Context, uuid string) error {
	return db.Delete(ctx, uuid)
}

func Update(ctx context.Context, uuid string, todo *types.Todo) error {
	return db.Update(ctx, uuid, todo)
}

func SetComplete(ctx context.Context, uuid string, status bool) error {
	return db.SetComplete(ctx, uuid, status)
}

func GetTodoByStatus(ctx context.Context, status bool) (titles []string, err error) {
	return db.GetTodoByStatus(ctx, status)
}

func GetTodoByPriority(ctx context.Context, priority string) (titles []string, err error) {
	return db.GetTodoByPriority(ctx, priority)
}
