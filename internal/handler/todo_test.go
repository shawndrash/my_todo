package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"my_todo/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDatabase 是一个模拟的数据库实现
type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) GetTodos() ([]model.Todo, error) {
	args := m.Called()
	return args.Get(0).([]model.Todo), args.Error(1)
}

func (m *MockDatabase) CreateTodo(todo model.Todo) error {
	args := m.Called(todo)
	return args.Error(0)
}

func (m *MockDatabase) UpdateTodo(id int64, todo model.Todo) error {
	args := m.Called(id, todo)
	return args.Error(0)
}

func (m *MockDatabase) DeleteTodo(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDatabase) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestGetTodos(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockTodos      []model.Todo
		mockError      error
		expectedStatus int
	}{
		{
			name: "成功获取待办事项",
			mockTodos: []model.Todo{
				{Id: 1, Text: "测试任务1", Completed: false},
				{Id: 2, Text: "测试任务2", Completed: true},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "数据库错误",
			mockTodos:      nil,
			mockError:      errors.New("数据库错误"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDatabase)
			mockDB.On("GetTodos").Return(tt.mockTodos, tt.mockError)

			handler := NewTodoHandler(mockDB)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			handler.GetTodos(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.mockError == nil {
				var response []model.Todo
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.mockTodos, response)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func TestCreateTodo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		input          model.Todo
		mockError      error
		expectedStatus int
	}{
		{
			name:           "成功创建待办事项",
			input:          model.Todo{Text: "新任务", Completed: false},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "数据库错误",
			input:          model.Todo{Text: "新任务", Completed: false},
			mockError:      errors.New("数据库错误"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDatabase)
			mockDB.On("CreateTodo", mock.AnythingOfType("model.Todo")).Return(tt.mockError)

			handler := NewTodoHandler(mockDB)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonData, _ := json.Marshal(tt.input)
			c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer(jsonData))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.CreateTodo(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockDB.AssertExpectations(t)
		})
	}
}

func TestUpdateTodo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		todoID         string
		input          model.Todo
		mockError      error
		expectedStatus int
	}{
		{
			name:           "成功更新待办事项",
			todoID:         "1",
			input:          model.Todo{Text: "更新的任务", Completed: true},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "无效的ID",
			todoID:         "invalid",
			input:          model.Todo{Text: "更新的任务", Completed: true},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDatabase)
			if tt.todoID != "invalid" {
				mockDB.On("UpdateTodo", int64(1), tt.input).Return(tt.mockError)
			}

			handler := NewTodoHandler(mockDB)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: tt.todoID}}

			jsonData, _ := json.Marshal(tt.input)
			c.Request = httptest.NewRequest("PUT", "/", bytes.NewBuffer(jsonData))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.UpdateTodo(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.todoID != "invalid" {
				mockDB.AssertExpectations(t)
			}
		})
	}
}

func TestDeleteTodo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		todoID         string
		mockError      error
		expectedStatus int
	}{
		{
			name:           "成功删除待办事项",
			todoID:         "1",
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "无效的ID",
			todoID:         "invalid",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "数据库错误",
			todoID:         "1",
			mockError:      errors.New("数据库错误"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDatabase)
			if tt.todoID != "invalid" {
				mockDB.On("DeleteTodo", int64(1)).Return(tt.mockError)
			}

			handler := NewTodoHandler(mockDB)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: tt.todoID}}

			handler.DeleteTodo(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.todoID != "invalid" {
				mockDB.AssertExpectations(t)
			}
		})
	}
}
