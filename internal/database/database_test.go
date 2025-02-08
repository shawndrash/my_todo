package database

import (
    "testing"
    "my_todo/internal/model"
    "github.com/stretchr/testify/assert"
)

func TestSQLiteDB(t *testing.T) {
    // 使用内存数据库进行测试
    db, err := NewDatabase(":memory:")
    assert.NoError(t, err)
    defer db.Close()

    // 测试创建待办事项
    t.Run("创建待办事项", func(t *testing.T) {
        todo := model.Todo{
            Text:      "测试任务",
            Completed: false,
        }
        err := db.CreateTodo(todo)
        assert.NoError(t, err)

        // 验证创建是否成功
        todos, err := db.GetTodos()
        assert.NoError(t, err)
        assert.Len(t, todos, 1)
        assert.Equal(t, "测试任务", todos[0].Text)
        assert.False(t, todos[0].Completed)
    })

    // 测试获取待办事项列表
    t.Run("获取待办事项列表", func(t *testing.T) {
        todos, err := db.GetTodos()
        assert.NoError(t, err)
        assert.NotEmpty(t, todos)
    })

    // 测试更新待办事项
    t.Run("更新待办事项", func(t *testing.T) {
        todos, err := db.GetTodos()
        assert.NoError(t, err)
        assert.NotEmpty(t, todos)

        firstTodo := todos[0]
        updatedTodo := model.Todo{
            Text:      "更新后的任务",
            Completed: true,
        }

        err = db.UpdateTodo(firstTodo.Id, updatedTodo)
        assert.NoError(t, err)

        // 验证更新是否成功
        todos, err = db.GetTodos()
        assert.NoError(t, err)
        assert.Equal(t, "更新后的任务", todos[0].Text)
        assert.True(t, todos[0].Completed)
    })

    // 测试删除待办事项
    t.Run("删除待办事项", func(t *testing.T) {
        todos, err := db.GetTodos()
        assert.NoError(t, err)
        assert.NotEmpty(t, todos)

        firstTodo := todos[0]
        err = db.DeleteTodo(firstTodo.Id)
        assert.NoError(t, err)

        // 验证删除是否成功
        todos, err = db.GetTodos()
        assert.NoError(t, err)
        assert.Empty(t, todos)
    })

    // 测试错误情况
    t.Run("无效ID更新", func(t *testing.T) {
        todo := model.Todo{
            Text:      "测试任务",
            Completed: false,
        }
        err := db.UpdateTodo(999, todo)
        assert.NoError(t, err) // SQLite不会为不存在的ID返回错误
    })

    t.Run("无效ID删除", func(t *testing.T) {
        err := db.DeleteTodo(999)
        assert.NoError(t, err) // SQLite不会为不存在的ID返回错误
    })
}
