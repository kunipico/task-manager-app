'use client';

import React, { useEffect, useState } from 'react';
import TaskCard from '../components/Cards';
import AddTaskModal from '../components/AddTaskModal'
import TimeAnalysisModal from '../components/TimeAnalysisModal'
import DocsModal from '../components/DocsModal'

import { useRouter } from 'next/navigation';


interface Task {
  Task_ID: number;
  Task_Name: string;
  Task_Details: string;
  Task_Done: string;
}

export default function Home() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [isAddTaskModalOpen, setIsAddTaskModalOpen] = useState(false);

  const router = useRouter(); // useRouterをLoginPage内で定義



  // 時間解析モーダル処理
  const [isTimeAnalysisModalOpen, setIsTimeAnalysisModalOpen] = useState(false);
  const [selectedTaskId, setSelectedTaskId] = useState<number | null>(null);
  const [selectedTaskName, setSelectedTaskName] = useState<string | null>(null);
  // モーダルを開く
  const openTimeAnalysis = (taskId:number, taskName: string) => {
    setSelectedTaskId(taskId);
    setSelectedTaskName(taskName);
    setIsTimeAnalysisModalOpen(true);
  }
  // モーダルを閉じる
  const closeTimeAnalysisModal = () => {
    setSelectedTaskId(null);
    setSelectedTaskName(null);
    setIsTimeAnalysisModalOpen(false);
  }



  // ドキュメントモーダル処理
  const [isDocsModalOpen, setIsDocsModalOpen] = useState(false);
  // モーダルを開く
  const openDocs = (taskId:number, taskName: string) => {
    setSelectedTaskId(taskId);
    setSelectedTaskName(taskName);
    setIsDocsModalOpen(true);
  }
  // モーダルを閉じる
  const closeDocsModal = () => {
    setSelectedTaskId(null);
    setSelectedTaskName(null);
    setIsDocsModalOpen(false);
  }



  useEffect(() => {
    fetchTasks();
  }, []);
  
  // タスクをGoAPIから取得
  const fetchTasks = async () => {
    const res = await fetch('http://localhost:8080/tasks', {
    // const res = await fetch('/api/tasks', {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include"
    });
    if (!res.ok){
      const errorText = await res.json(); // JSONでない場合に備えてテキストとして取得
      console.log('Error fetching tasks:', errorText);
      router.push("/login");
      return;
    }
    const data = await res.json();
    console.log('tasks : ',data);
    setTasks(data);
  };

  // モーダルを開く
  const openAddTaskModal = () => setIsAddTaskModalOpen(true);

  // モーダルを閉じる
  const closeAddTaskModal = () => setIsAddTaskModalOpen(false);

  // タスクを追加
  const addTask = async (taskName: string, taskDetails: string) => {
    const newTask: Task = {
      Task_ID: tasks.length + 1,
      Task_Name: taskName,
      Task_Details: taskDetails,
      Task_Done: "Standby",
    };
    const res = await fetch('http://localhost:8080/tasks', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: "include",
      body: JSON.stringify({ Task_Name: newTask.Task_Name, Task_Details:newTask.Task_Details, Task_Done: 'Standby' }),
    });
    if (res.ok) {
      setTasks([...tasks, newTask]);
      // fetchTasks();   // 再度タスクを取得
    };
  }

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
    const task = tasks.find((t) => t.Task_ID === taskId);
    if (!task) return;
    
    // 状態を切り替えるロジック
    let newStatus;
    if (task.Task_Done === 'Standby') {
      newStatus = 'Inprogress';
    } else if (task.Task_Done === 'Inprogress') {
      newStatus = 'Standby';
    }

    const res = await fetch(`http://localhost:8080/tasks/${taskId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },credentials: "include",
      body: JSON.stringify({ Task_Done: newStatus }),
    });

    if (res.ok) {
      fetchTasks(); // 状態変更後にタスク一覧を再取得
    }
  };

  const logout = async(e:React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    try {
      const res = await fetch("http://localhost:8080/logout", {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
      });
      const data = await res.text();
      console.log("data : ",data);
      if (res.ok) {
          router.push("/login");
        return
      }
    } catch (error) {
      console.log("ログアウトに失敗しました");
      router.push("/login");
    }
  };

  // const openDocs = async (taskId: number) => {

  // }

  return (
    <div className="justify-center items-center min-h-screen py-4 px-10 bg-gray-100">
      <h1 className="text-gray-500 text-2xl font-semibold mb-6 text-center">Task Manager</h1>
      <div className='flex justify-between items-center'>
        <div className="mb-4">
          <button
            className="border border-gray-300  hover:bg-green-500 hover:text-white bg-white text-green-500 bg py-2 px-4 rounded"
            onClick={openAddTaskModal}
          >
            追加
          </button>
        </div>
        <div className="mb-4">
          <button
            className="border border-gray-300 hover:bg-blue-500 hover:text-white bg-white text-blue-500 bg py-2 px-4 rounded"
            onClick={logout}
          >
            Logout
          </button>
        </div>
      </div>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {tasks ? (
          tasks.map((task,index) => (
            <TaskCard
              // key={`${task.Task_Id}-${index}`}
              key={`${index}`}
              task={task}
              onToggleTask={toggleTask}
              onDeleteTask={deleteTask}
              onOpenDocs={openDocs}
              onOpenTimeAnalysis={openTimeAnalysis}
            />
          ))
        ) : (
          <p>タスクがありません。</p>
        )}
      </div>
       {/* モーダルコンポーネントを使用 */}
       <AddTaskModal
        isOpen={isAddTaskModalOpen}
        onClose={closeAddTaskModal}
        onSubmit={addTask}
      />
      <TimeAnalysisModal
        taskName={selectedTaskName!}
        taskId={selectedTaskId!}
        isOpen={isTimeAnalysisModalOpen}
        onClose={closeTimeAnalysisModal}
        // onSubmit={task.Task_ID}
      />
      <DocsModal
        taskName={selectedTaskName!}
        taskId={selectedTaskId!}
        isOpen={isDocsModalOpen}
        onClose={closeDocsModal}
      />
    </div>
  );
}