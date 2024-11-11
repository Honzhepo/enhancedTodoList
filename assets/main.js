window.onload = function() {
    fetchTasks();
}
const dateInput = document.getElementById("taskDeadline");
dateInput.valueAsDate = new Date();

// Fetch and render tasks
function fetchTasks() {
    getTasks().then(tasks => {
        // Sort tasks by deadline, with the earliest deadlines first
        tasks.sort((a, b) => new Date(a.deadline) - new Date(b.deadline));

        const taskList = document.getElementById('taskList');
        taskList.innerHTML = ''; // Clear the current list

        tasks.forEach(task => {
            const taskDiv = document.createElement('div');
            taskDiv.classList.add('task', task.type);

            // Toggle task completion on click
            taskDiv.onclick = function() {
                toggleTask(task.id);
            };

            if (task.done) taskDiv.classList.add('done');

            // Format deadline as day-month-year
            const deadlineDate = new Date(task.deadline);
            const formattedDate = `${deadlineDate.getDate().toString().padStart(2, '0')}-${(deadlineDate.getMonth() + 1).toString().padStart(2, '0')}-${deadlineDate.getFullYear()}`;

            // Display task deadline and title with a delete button
            taskDiv.innerHTML = `
                <span>${formattedDate}</span> <span>${task.title}</span>
                <button class="deleteBtn" onclick="deleteTask(${task.id})">🗑</button>
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

// Delete task
function deleteTask(id) {
    deleteTaskGo(id).then(() => fetchTasks());
}