package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"api-doc-example/internal/models"
)

type TodoHandler struct {
	todos  map[int]models.Todo
	nextID int
}

func NewTodoHandler() *TodoHandler {
	return &TodoHandler{
		todos:  make(map[int]models.Todo),
		nextID: 1,
	}
}

// ListTodos godoc
// @Summary      List all tasks
// @Description  Returns a list of all tasks wrapped in SuccessResponse
// @Tags         todos
// @Produce      json
// @Success      200  {object}  models.SuccessResponse{data=[]models.Todo}  "List of tasks"
// @Failure      500  {object}  models.ErrorResponse  "Internal server error"
// @Router       /todos [get]
func (h *TodoHandler) ListTodos(w http.ResponseWriter, r *http.Request) {
	var todos []models.Todo
	for _, todo := range h.todos {
		todos = append(todos, todo)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(models.SuccessResponse{
		Success: true,
		Data:    todos,
	})
}

// GetTodo godoc
// @Summary      Get a task by ID
// @Description  Returns a single task by its ID wrapped in SuccessResponse
// @Tags         todos
// @Produce      json
// @Param        id   path      int   true  "Task ID"
// @Success      200  {object}  models.SuccessResponse{data=models.Todo}  "Task found"
// @Failure      400  {object}  models.ErrorResponse  "Invalid ID"
// @Failure      404  {object}  models.ErrorResponse  "Task not found"
// @Router       /todos/{id} [get]
func (h *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "400",
			Message: "Неверный ID",
		})
		return
	}

	todo, exists := h.todos[id]
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "404",
			Message: "Задача не найдена",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(models.SuccessResponse{
		Success: true,
		Data:    todo,
	})
}

// CreateTodo godoc
// @Summary      Create a new task
// @Description  Creates a new task with the provided data
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        request  body  models.CreateTodoRequest  true  "Task data"
// @Success      201  {object}  models.SuccessResponse{data=models.Todo}  "Task created"
// @Failure      400  {object}  models.ErrorResponse  "Invalid request"
// @Router       /todos [post]
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "400",
			Message: "Невалидный запрос",
		})
		return
	}

	todo := models.Todo{
		ID:          h.nextID,
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		Done:        req.Done,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	h.todos[h.nextID] = todo
	h.nextID++

	log.Printf("Создана новая задача: %+v", todo)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(models.SuccessResponse{
		Success: true,
		Data:    todo,
	})
}

// UpdateTodo godoc
// @Summary      Update an existing task
// @Description  Updates a task by its ID. All fields are optional.
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        id       path  int                      true  "Task ID"
// @Param        request  body  models.UpdateTodoRequest  true  "Updated task data"
// @Success      200  {object}  models.SuccessResponse{data=models.Todo}  "Task updated"
// @Failure      400  {object}  models.ErrorResponse  "Invalid ID or request"
// @Failure      404  {object}  models.ErrorResponse  "Task not found"
// @Router       /todos/{id} [put]
func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "400",
			Message: "Неверный ID",
		})
		return
	}

	var req models.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "400",
			Message: "Невалидный запрос",
		})
		return
	}

	todo, exists := h.todos[id]
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "404",
			Message: "Задача не найдена",
		})
		return
	}

	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Description != nil {
		todo.Description = *req.Description
	}
	if req.Done != nil {
		todo.Done = *req.Done
	}
	todo.UpdatedAt = time.Now()

	h.todos[id] = todo

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(models.SuccessResponse{
		Success: true,
		Data:    todo,
	})
}

// DeleteTodo godoc
// @Summary      Delete a task
// @Description  Deletes a task by its ID
// @Tags         todos
// @Produce      json
// @Param        id   path  int  true  "Task ID"
// @Success      204  "Task deleted"
// @Failure      400  {object}  models.ErrorResponse  "Invalid ID"
// @Failure      404  {object}  models.ErrorResponse  "Task not found"
// @Router       /todos/{id} [delete]
func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "400",
			Message: "Неверный ID",
		})
		return
	}

	if _, exists := h.todos[id]; !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "404",
			Message: "Задача не найдена",
		})
		return
	}

	delete(h.todos, id)

	w.WriteHeader(http.StatusNoContent)
}
