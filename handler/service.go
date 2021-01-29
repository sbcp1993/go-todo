package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"go-todo/db"
	"go-todo/runtime"
	"go-todo/types"
	"go-todo/utils"
)

type todoHandler struct {
	postgres *db.DBConnection
}

func (handler *todoHandler) addUser(w http.ResponseWriter, r *http.Request) {
	ctx := db.SetBackend(r.Context(), handler.postgres)

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var user types.User

	if err := json.Unmarshal(b, &user); err != nil {
		responseError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := runtime.InsertUser(ctx, &user); err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (handler *todoHandler) getToken(w http.ResponseWriter, r *http.Request) {
	ctx := db.SetBackend(r.Context(), handler.postgres)
	var res map[string]interface{}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var user types.User

	if err := json.Unmarshal(b, &user); err != nil {
		responseError(w, http.StatusBadRequest, err.Error())
		return
	}

	if res, err = runtime.Login(ctx, user.Name, user.PasswordHash); err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseOk(w, res)
}

func (handler *todoHandler) createTodo(w http.ResponseWriter, r *http.Request) {
	ctx := db.SetBackend(r.Context(), handler.postgres)

	if err := utils.TokenValid(r); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var todo types.Todo
	if err := json.Unmarshal(b, &todo); err != nil {
		responseError(w, http.StatusBadRequest, err.Error())
		return
	}

	uuid, err := runtime.Insert(ctx, &todo)
	if err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseOk(w, uuid)
}

func (handler *todoHandler) logout(w http.ResponseWriter, r *http.Request) {
	os.Setenv("ACCESS_SECRET", "")

	w.WriteHeader(http.StatusOK)
}

func (handler *todoHandler) deleteTodo(w http.ResponseWriter, r *http.Request) {
	ctx := db.SetBackend(r.Context(), handler.postgres)

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var req struct {
		UUID string `json:"id"`
	}
	if err := json.Unmarshal(b, &req); err != nil {
		responseError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := runtime.Delete(ctx, req.UUID); err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *todoHandler) updateTodo(w http.ResponseWriter, r *http.Request) {
	ctx := db.SetBackend(r.Context(), handler.postgres)

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var req struct {
		UUID string     `json:"id"`
		Todo types.Todo `json:"todo"`
	}
	if err := json.Unmarshal(b, &req); err != nil {
		responseError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := runtime.Update(ctx, req.UUID, &req.Todo); err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *todoHandler) completeTodo(w http.ResponseWriter, r *http.Request) {
	ctx := db.SetBackend(r.Context(), handler.postgres)

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var req struct {
		UUID   string `json:"id"`
		Status bool   `json:"status"`
	}
	if err := json.Unmarshal(b, &req); err != nil {
		responseError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := runtime.SetComplete(ctx, req.UUID, req.Status); err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *todoHandler) getTodoByPriority(w http.ResponseWriter, r *http.Request) {
	var todoTitles []string
	var err error

	ctx := db.SetBackend(r.Context(), handler.postgres)

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var req struct {
		Priority string `json:"priority"`
	}

	if err := json.Unmarshal(b, &req); err != nil {
		responseError(w, http.StatusBadRequest, err.Error())
		return
	}

	if todoTitles, err = runtime.GetTodoByPriority(ctx, req.Priority); err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseOk(w, todoTitles)
}

func (handler *todoHandler) getTodoByStatus(w http.ResponseWriter, r *http.Request) {
	var todoTitles []string
	var err error

	ctx := db.SetBackend(r.Context(), handler.postgres)

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var req struct {
		Status bool `json:"status"`
	}

	if err := json.Unmarshal(b, &req); err != nil {
		responseError(w, http.StatusBadRequest, err.Error())
		return
	}

	if todoTitles, err = runtime.GetTodoByStatus(ctx, req.Status); err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseOk(w, todoTitles)
}

func responseOk(w http.ResponseWriter, body interface{}) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(body)
}

func responseError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")

	body := map[string]string{
		"error": message,
	}
	json.NewEncoder(w).Encode(body)
}
