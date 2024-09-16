// components/Card.tsx

import React from 'react';

interface Task {
  Task_Id: number;
  Task_Name: string;
  Task_Done: boolean;
}

interface CardProps {
  task: Task;
  onToggleTask: (taskId: number) => void;
  onDeleteTask: (taskId: number) => void;
}

const Card: React.FC<CardProps> = ({ task, onToggleTask, onDeleteTask}) => {
  return (
    <div key={task.Task_Id} className="bg-white shadow-md rounded-lg p-6">
      <h2 className="text-xl font-semibold mb-2">{task.Task_Name}</h2>
      <p className={`text-lg ${task.Task_Done ? 'text-green-600' : 'text-red-600'}`}>
        {task.Task_Done ? '完了' : '未完了'}
      </p>
      <div className="flex justify-between mt-4">
        <button
          className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
          onClick={() => onToggleTask(task.Task_Id)}
        >
          状態切り替え
        </button>
        <button
          className="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded"
          onClick={() => onDeleteTask(task.Task_Id)}
        >
          削除
        </button>
      </div>
    </div>
    // <div className="bg-blue-200 shadow-md rounded-lg p-6">
    //   <h2 className="text-blue-900 text-xl font-semibold mb-2">{task.Task_Name}</h2>
    //   <p className={`text-lg ${task.Task_Done ? 'text-green-600' : 'text-red-600'}`}>
    //     {task.Task_Done ? '完了' : '未完了'}
    //   </p>
    // </div>
  );
};

export default Card;