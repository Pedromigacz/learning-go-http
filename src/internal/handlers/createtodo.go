package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Pedromigacz/learning-go-http/src/internal/store"
)

type CreateTodoHandler struct {
	todos *[]store.Todo
}

type CreateTodoHandlerParams struct {
	Todos *[]store.Todo
}

func NewCreateTodoHandler(params CreateTodoHandlerParams) *CreateTodoHandler {
	return &CreateTodoHandler{
		todos: params.Todos,
	}
}

func (h *CreateTodoHandler) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	todo := store.Todo{}
	err := json.NewDecoder(r.Body).Decode(&todo)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	*h.todos = append(*h.todos, todo)

	w.WriteHeader(http.StatusCreated)
}
