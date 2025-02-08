package router

import (
	"github.com/gin-gonic/gin"
	"my_todo/internal/handler"
)

func SetupRoutes(r *gin.Engine, todoHandler *handler.TodoHandler) {
	v1 := r.Group("/v1")
	{
		todos := v1.Group("/todos")
		{
			todos.GET("/", todoHandler.GetTodos)
			todos.POST("/", todoHandler.CreateTodo)
			todos.PUT("/:id", todoHandler.UpdateTodo)
			todos.DELETE("/:id", todoHandler.DeleteTodo)
		}
	}
}
