import React from 'react';

interface AddTaskModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (taskName: string, taskDetails: string) => void;
}

export default function AddTaskModal({ isOpen, onClose, onSubmit }: AddTaskModalProps) {
  const [taskName, setTaskName] = React.useState('');
  const [taskDetails, setTaskDetails] = React.useState('');

  if (!isOpen) return null; // モーダルが開かれていない場合は何も描画しない

  const handleSubmit = () => {
    if (taskName.trim() && taskDetails.trim()) {
      onSubmit(taskName, taskDetails); // タスクを追加
      setTaskName(''); // 入力内容をリセット
      setTaskDetails('');
      onClose(); // モーダルを閉じる
    } else {
      alert('タスク名と詳細を入力してください');
    }
  };

  return (
    <div className="fixed inset-0 flex justify-center items-center bg-black bg-opacity-50">
      <div className="bg-white p-6 rounded shadow-lg w-96">
        <h2 className="text-xl font-semibold mb-4">タスクを追加</h2>
        <input
          type="text"
          placeholder="タスク名"
          value={taskName}
          onChange={(e) => setTaskName(e.target.value)}
          className="w-full mb-4 p-2 border rounded"
        />
        <textarea
          placeholder="タスク詳細"
          value={taskDetails}
          onChange={(e) => setTaskDetails(e.target.value)}
          className="w-full mb-4 p-2 border rounded"
        />
        <div className="flex justify-end">
          <button
            onClick={onClose}
            className="bg-gray-300 hover:bg-gray-400 text-gray-800 py-2 px-4 rounded mr-2"
          >
            キャンセル
          </button>
          <button
            onClick={handleSubmit}
            className="bg-blue-500 hover:bg-blue-600 text-white py-2 px-4 rounded"
          >
            追加
          </button>
        </div>
      </div>
    </div>
  );
}