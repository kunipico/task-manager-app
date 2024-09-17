'use client';

import React, { useEffect, useState } from 'react';
import TaskCard from '../components/Cards';

interface Task {
  Task_Id: number;
  Task_Name: string;
  Task_Done: boolean;
}

export default function Home() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [newTask, setNewTask] = useState<string>('');

  useEffect(() => {
    fetchTasks();
  }, []);
  
  // タスクをGoAPIから取得
  const fetchTasks = async () => {
    const res = await fetch('http://localhost:8080/tasks');
    const data = await res.json();
    console.log('tasks : ',data);
    setTasks(data);
  };

  // タスク追加
  const addTask = async () => {
    if (!newTask) return;

    const res = await fetch('http://localhost:8080/tasks', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ Task_Name: newTask, Task_Done: false }),
    });

    if (res.ok) {
      setNewTask(''); // 入力欄をクリア
      fetchTasks();   // 再度タスクを取得
    }
  };

  // タスク削除
  const deleteTask = async (taskId: number) => {
    const res = await fetch(`http://localhost:8080/tasks/${taskId}`, {
      method: 'DELETE',
    });

    if (res.ok) {
      fetchTasks(); // 削除後にタスク一覧を再取得
    }
  };

  // タスク状態の切り替え
  const toggleTask = async (taskId: number) => {
    const task = tasks.find((t) => t.Task_Id === taskId);
    if (!task) return;

    const res = await fetch(`http://localhost:8080/tasks/${taskId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ Task_Done: !task.Task_Done }),
    });

    if (res.ok) {
      fetchTasks(); // 状態変更後にタスク一覧を再取得
    }
  };

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-gray-500 text-2xl font-semibold mb-6 text-center">Task Manager</h1>
      <div className="mb-4">
        <input
          type="text"
          value={newTask}
          onChange={(e) => setNewTask(e.target.value)}
          placeholder="新しいタスクを追加"
          className="border border-gray-300 p-2 mr-2 rounded"
        />
        <button
          className="border border-gray-300 bg-white hover:bg-green-500 text-green-300 bg py-2 px-4 rounded"
          onClick={addTask}
        >
          追加
        </button>
      </div>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {tasks ? (
          tasks.map((task) => (
            <TaskCard
              key={task.Task_Id}
              task={task}
              onToggleTask={toggleTask}
              onDeleteTask={deleteTask}
            />
          ))
        ) : (
          <p>タスクがありません。</p>
        )}
      </div>
    </div>
  );
}

