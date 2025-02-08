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
            todoItem.dataset.id = todo.id;
            if (todo.completed) {
                todoItem.classList.add('completed');
            }

            // 创建文本元素
            const todoText = document.createElement('span');
            todoText.className = 'todo-text';
            todoText.textContent = todo.text;
            if (todo.completed) {
                todoText.classList.add('completed');
            }
            todoItem.appendChild(todoText);

            // 创建按钮容器
            const buttonContainer = document.createElement('div');
            buttonContainer.className = 'todo-buttons';

            // 添加完成按钮
            const completeButton = document.createElement('button');
            completeButton.className = 'complete-btn';
            completeButton.textContent = todo.completed ? '取消完成' : '完成';
            completeButton.addEventListener('click', () => toggleComplete(todo.id, !todo.completed));
            buttonContainer.appendChild(completeButton);

            // 添加编辑按钮
            const editButton = document.createElement('button');
            editButton.className = 'edit-btn';
            editButton.textContent = '编辑';
            editButton.addEventListener('click', () => editTodo(todo.id, todoText));
            buttonContainer.appendChild(editButton);

            // 添加删除按钮
            const deleteButton = document.createElement('button');
            deleteButton.className = 'delete-btn';
            deleteButton.textContent = '删除';
            deleteButton.addEventListener('click', () => deleteTodo(todo.id));
            buttonContainer.appendChild(deleteButton);

            todoItem.appendChild(buttonContainer);
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

// 编辑待办事项
async function editTodo(id, todoTextElement) {
    const currentText = todoTextElement.textContent;
    const newText = prompt('编辑待办事项:', currentText);

    if (newText !== null && newText.trim() !== '') {
        try {
            const response = await fetch(`/v1/todos/${id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    text: newText.trim(),
                    completed: todoTextElement.classList.contains('completed')
                })
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            // 更新前端显示
            todoTextElement.textContent = newText.trim();
        } catch (error) {
            console.error('Error updating todo:', error);
        }
    }
}

// 切换完成状态
async function toggleComplete(id, completed) {
    try {
        const todoItem = document.querySelector(`li[data-id="${id}"]`);
        const todoText = todoItem.querySelector('.todo-text');
        const completeButton = todoItem.querySelector('.complete-btn');

        const response = await fetch(`/v1/todos/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                text: todoText.textContent,
                completed: completed
            })
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();

        // 更新UI
        const li = todoItem;
        if (data.completed) {
            todoText.classList.add('completed');
            li.classList.add('completed');
            completeButton.textContent = '取消完成';
        } else {
            todoText.classList.remove('completed');
            li.classList.remove('completed');
            completeButton.textContent = '完成';
        }

        // 重新绑定点击事件
        completeButton.onclick = () => toggleComplete(id, !data.completed);

    } catch (error) {
        console.error('Error updating todo:', error);
    }
}

// 更新待办事项
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

        // 更新前端显示通过重新加载来确保状态一致
        getTodos();
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
