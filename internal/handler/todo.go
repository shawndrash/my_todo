package handler

import (
	"my_todo/internal/database"
	"my_todo/internal/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TodoHandler 处理 Todo 相关的请求
type TodoHandler struct {
	db database.Database
}

func NewTodoHandler(db database.Database) *TodoHandler {
	return &TodoHandler{db: db}
}

// GetTodos 获取所有的 Todo
func (h *TodoHandler) GetTodos(c *gin.Context) {
	todos, err := h.db.GetTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

// CreateTodo 创建一个新的 Todo
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var todo model.Todo
	if err := c.BindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.CreateTodo(todo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, todo)
}

// UpdateTodo 更新一个 Todo
func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	id := c.Param("id")
	int_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var todo model.Todo
	if err := c.BindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.UpdateTodo(int_id, todo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todo)
}

// DeleteTodo 删除一个 Todo
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	id := c.Param("id")
	int_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.db.DeleteTodo(int_id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}
