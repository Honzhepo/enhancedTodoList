// Initialize and fetch tasks on load
window.onload = function() {
    fetchTasks();
}

// Fetch and render tasks
function fetchTasks() {
    getTasks().then(tasks => {
        const taskList = document.getElementById('taskList');
        taskList.innerHTML = '';
        tasks.forEach(task => {
            const taskDiv = document.createElement('div');
            taskDiv.classList.add('task', task.type);
            if (task.done) taskDiv.classList.add('done');

            taskDiv.innerHTML = `
                <span>${task.title} - ${task.deadline.split('T')[0]}</span>
                <button onclick="toggleTask(${task.id})">${task.done ? 'Undo' : 'Complete'}</button>
            `;
            taskList.appendChild(taskDiv);
        });
    });
}

// Add new task
function addNewTask() {
    const title = document.getElementById('taskTitle').value;
    const deadline = document.getElementById('taskDeadline').value;
    const type = document.getElementById('taskType').value;

    if (title && deadline) {
        addTask(title, deadline, type).then(() => fetchTasks());
    } else {
        alert('Please fill out all fields');
    }
}

// Toggle task done/undone
function toggleTask(id) {
    toggleTaskDone(id).then(() => fetchTasks());
}
