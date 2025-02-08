package database

import (
	"database/sql"
	"fmt"

	"my_todo/internal/model"

	_ "github.com/mattn/go-sqlite3"
)

// Database 接口定义了数据库操作的方法
type Database interface {
	GetTodos() ([]model.Todo, error)
	CreateTodo(todo model.Todo) error
	UpdateTodo(id int64, todo model.Todo) error
	DeleteTodo(id int64) error
	Close() error
}

type SQLiteDB struct {
	db *sql.DB
}

func NewDatabase(connString string) (Database, error) {
	db, err := sql.Open("sqlite3", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 测试数据库连接
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// 创建表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			text TEXT NOT NULL,
			completed BOOLEAN NOT NULL CHECK (completed IN (0, 1))
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &SQLiteDB{db: db}, nil
}

func (s *SQLiteDB) GetTodos() ([]model.Todo, error) {
	rows, err := s.db.Query("SELECT id, text, completed FROM todos")
	if err != nil {
		return nil, fmt.Errorf("failed to query todos: %w", err)
	}
	defer rows.Close()

	var todos []model.Todo
	for rows.Next() {
		var todo model.Todo
		err := rows.Scan(&todo.Id, &todo.Text, &todo.Completed)
		if err != nil {
			return nil, fmt.Errorf("failed to scan todo: %w", err)
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func (s *SQLiteDB) CreateTodo(todo model.Todo) error {
	_, err := s.db.Exec("INSERT INTO todos (text, completed) VALUES (?, ?)", todo.Text, todo.Completed)
	if err != nil {
		return fmt.Errorf("failed to insert todo: %w", err)
	}
	return nil
}

func (s *SQLiteDB) UpdateTodo(id int64, todo model.Todo) error {
	_, err := s.db.Exec("UPDATE todos SET text = ?, completed = ? WHERE id = ?", todo.Text, todo.Completed, id)
	if err != nil {
		return fmt.Errorf("failed to update todo: %w", err)
	}
	return nil
}

func (s *SQLiteDB) DeleteTodo(id int64) error {
	_, err := s.db.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}
	return nil
}

func (s *SQLiteDB) Close() error {
	return s.db.Close()
}
