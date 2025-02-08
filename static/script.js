const todoList = document.getElementById('todo-list');
const newTodoInput = document.getElementById('new-todo-input');
const addTodoButton = document.getElementById('add-todo-button');

// 获取所有待办事项
async function getTodos() {
    try {
        const response = await fetch('/v1/todos');
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const todos = await response.json();
        todoList.innerHTML = '';
        todos.forEach(todo => {
            const todoItem = document.createElement('li');
            todoItem.textContent = todo.text;
            // 添加id
            todoItem.dataset.id = todo.id;

            // 添加删除按钮
            const deleteButton = document.createElement('button');
            deleteButton.textContent = '删除';
            deleteButton.addEventListener('click', () => deleteTodo(todo.id));
            todoItem.appendChild(deleteButton);


            todoList.appendChild(todoItem);
        });
    } catch (error) {
        console.error('Error fetching todos:', error);
    }
}

// 添加新的待办事项
async function addTodo() {
    const newTodoText = newTodoInput.value.trim();
    if (newTodoText !== '') {
        try {
            const response = await fetch('/v1/todos', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ text: newTodoText })
            });
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const newTodo = await response.json();
            // 重新加载所有待办事项以确保显示最新数据
            getTodos();
            newTodoInput.value = '';
        } catch (error) {
            console.error('Error adding todo:', error);
        }
    }
}

async function updateTodo(id, updatedTodo) {
    try {
        const response = await fetch(`/v1/todos/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(updatedTodo)
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        // 更新前端显示
        const todoItem = document.querySelector(`li[data-id="${id}"]`);
        if (todoItem) {
            todoItem.textContent = updatedTodo.text;

            // 重新添加删除按钮
            const deleteButton = document.createElement('button');
            deleteButton.textContent = '删除';
            deleteButton.addEventListener('click', () => deleteTodo(id));
            todoItem.appendChild(deleteButton);
        }

    } catch (error) {
        console.error('Error updating todo:', error);
    }
}

// 删除待办事项
async function deleteTodo(id) {
    try {
        const response = await fetch(`/v1/todos/${id}`, {
            method: 'DELETE'
        });
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        // 从列表中移除待办事项
        const todoItem = document.querySelector(`li[data-id="${id}"]`);
        if (todoItem) {
            todoList.removeChild(todoItem);
        }
    } catch (error) {
        console.error('Error deleting todo:', error);
    }
}


addTodoButton.addEventListener('click', addTodo);

// 页面加载时获取所有待办事项
getTodos();
