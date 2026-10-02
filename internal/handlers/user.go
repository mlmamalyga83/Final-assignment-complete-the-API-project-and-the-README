package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"

	"api-doc-example/internal/models"
)

type UserHandler struct {
	mu     sync.RWMutex
	users  map[int]models.User
	nextID int
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		users:  make(map[int]models.User),
		nextID: 1,
	}
}

// ListUsers godoc
// @Summary      List all users
// @Description  Returns a list of all registered users wrapped in SuccessResponse
// @Tags         users
// @Produce      json
// @Success      200  {object}  models.SuccessResponse{data=[]models.User}  "List of users"
// @Failure      500  {object}  models.ErrorResponse  "Internal server error"
// @Router       /users [get]
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var users []models.User
	for _, user := range h.users {
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(models.SuccessResponse{
		Success: true,
		Data:    users,
	})
}

// GetUser godoc
// @Summary      Get a user by ID
// @Description  Returns a single user by their ID wrapped in SuccessResponse
// @Tags         users
// @Produce      json
// @Param        id   path      int   true  "User ID"
// @Success      200  {object}  models.SuccessResponse{data=models.User}  "User found"
// @Failure      400  {object}  models.ErrorResponse  "Invalid ID"
// @Failure      404  {object}  models.ErrorResponse  "User not found"
// @Router       /users/{id} [get]
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
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

	h.mu.RLock()
	user, exists := h.users[id]
	h.mu.RUnlock()

	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "404",
			Message: "Пользователь не найден",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(models.SuccessResponse{
		Success: true,
		Data:    user,
	})
}

// CreateUser godoc
// @Summary      Create a new user
// @Description  Creates a new user with the provided email and name
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body  models.CreateUserRequest  true  "User data"
// @Success      201  {object}  models.SuccessResponse{data=models.User}  "User created"
// @Failure      400  {object}  models.ErrorResponse  "Invalid request"
// @Router       /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "400",
			Message: "Невалидный запрос",
		})
		return
	}

	if req.Email.len == 0 || req.Email.find("@") == -1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "400",
			Message: "Некорректный email",
		})
		return
	}
	if req.Name.len == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "400",
			Message: "Имя пользователя не может быть пустым",
		})
		return
	}
	if req.Name.len > 100 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "400",
			Message: "Имя пользователя не может превышать 100 символов",
		})
		return
	}

	h.mu.Lock()
	user := models.User{
		ID:        h.nextID,
		Email:     req.Email,
		Name:      req.Name,
		CreatedAt: time.Now(),
	}
	h.users[h.nextID] = user
	h.nextID++
	h.mu.Unlock()

	log.Printf("Создан новый пользователь: %+v", user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(models.SuccessResponse{
		Success: true,
		Data:    user,
	})
}

// UpdateUser godoc
// @Summary      Update an existing user
// @Description  Updates a user by their ID. Both email and name are optional.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id       path  int                      true  "User ID"
// @Param        request  body  models.UpdateUserRequest  true  "Updated user data"
// @Success      200  {object}  models.SuccessResponse{data=models.User}  "User updated"
// @Failure      400  {object}  models.ErrorResponse  "Invalid ID or request"
// @Failure      404  {object}  models.ErrorResponse  "User not found"
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
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

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "400",
			Message: "Невалидный запрос",
		})
		return
	}

	if req.Email != nil && req.Email.find("@") == -1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "400",
			Message: "Некорректный email",
		})
		return
	}
	if req.Name != nil && req.Name.len > 100 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "400",
			Message: "Имя пользователя не может превышать 100 символов",
		})
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	user, exists := h.users[id]
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "404",
			Message: "Пользователь не найден",
		})
		return
	}

	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Name != nil {
		user.Name = *req.Name
	}

	h.users[id] = user

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(models.SuccessResponse{
		Success: true,
		Data:    user,
	})
}

// DeleteUser godoc
// @Summary      Delete a user
// @Description  Deletes a user by their ID
// @Tags         users
// @Produce      json
// @Param        id   path  int  true  "User ID"
// @Success      204  "User deleted"
// @Failure      400  {object}  models.ErrorResponse  "Invalid ID"
// @Failure      404  {object}  models.ErrorResponse  "User not found"
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
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

	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.users[id]; !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "404",
			Message: "Пользователь не найден",
		})
		return
	}

	delete(h.users, id)

	w.WriteHeader(http.StatusNoContent)
}
