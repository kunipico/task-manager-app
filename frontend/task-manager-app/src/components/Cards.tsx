// components/Card.tsx

// import React, { useState } from 'react';

interface Task {
  Task_ID: number;
  Task_Name: string;
  Task_Details:string;
  Task_Done: string;
}

interface CardProps {
  task: Task;
  onToggleTask: (taskId: number) => void;
  onDeleteTask: (taskId: number) => void;
  onOpenDocs: (taskId: number, taskName: string) => void;
  onOpenTimeAnalysis: (taskId: number, taskName:string) => void;
}

const Card: React.FC<CardProps> = ({ task, onToggleTask, onDeleteTask, onOpenDocs, onOpenTimeAnalysis}) => {
  // 状態に応じた色とテキストを決定
  const getStatusInfo = (status: string) => {
    switch (status) {
      case 'Standby':
        return { task_status: 'Standby', colorClass: 'text-blue-400', borderColor: 'border-blue-400' };
      case 'Inprogress':
        return { task_status: 'Inprogress', colorClass: 'text-green-400', borderColor: 'border-green-400 scale-105' };
      case 'Done':
        return { task_status: 'Done', colorClass: 'text-red-400', borderColor: 'border-red-400' };
      default:
        return { task_status: 'Unknown', colorClass: 'text-gray-400', borderColor: 'border-gray-400' };
    }
  };

  const { task_status, colorClass, borderColor } = getStatusInfo(task.Task_Done);

  return (
    <div className={`bg-white shadow-md rounded-lg p-6 transform transition-transform duration-200 hover:scale-105 border ${borderColor}`}
    onClick={() => {
      onToggleTask(task.Task_ID);
    }}
    >
      <div className="flex items-center justify-between mb-2">
        <h2 className="text-gray-700 text-xl underline underline-offset-8 mb-2">{task.Task_Name}</h2>
        <p className={`text-base font-bold ${colorClass}`}>{task_status}</p>
      </div>
      <div className="text-gray-700 text-base mb-2">{task.Task_Details}</div>
      
      <div className="flex justify-between mt-4">
        <button
          className="border border-blue-400 text-blue-400 hover:bg-blue-400 hover:text-white rounded px-1"
          onClick={(e) => {
            e.stopPropagation(); // 親要素への伝播を防ぐ
            onOpenDocs(task.Task_ID, task.Task_Name)
          }}
        >
          Docs
        </button>
        <button
          className="border border-gray-300  hover:bg-green-500 hover:text-white bg-white text-green-500 bg py-2 px-1 rounded"
          onClick={(e) => {
            e.stopPropagation(); // 親要素への伝播を防ぐ
            onOpenTimeAnalysis(task.Task_ID,task.Task_Name)
          }}
        >
          TimeAnalysis
        </button>
        <button
          className="text-red-400 hover:bg-red-400 hover:text-white rounded px-1 border border-red-400"
          onClick={() => onDeleteTask(task.Task_ID)}
        >
          Done
        </button>
      </div>
    </div>
  );
};

export default Card;