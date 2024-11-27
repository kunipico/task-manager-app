import React from 'react';

interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (taskName: string, taskDetails: string) => void;
}

export default function Modal({ isOpen, onClose, onSubmit }: ModalProps) {
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
    // onSubmit(taskName,taskDetails);
    // setTaskName(''); // 入力内容をリセット
    // onClose(); // モーダルを閉じる
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
        {/* <input
          type="text"
          value={taskName}
          onChange={(e) => setTaskName(e.target.value)}
          placeholder="新しいタスクを入力"
          className="w-full border border-gray-300 p-2 mb-4 rounded"
        />
        <div className="flex justify-end gap-2">
          <button
            className="bg-gray-300 text-gray-700 px-4 py-2 rounded hover:bg-gray-400"
            onClick={onClose}
          >
            キャンセル
          </button>
          <button
            className="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600"
            onClick={handleSubmit}
          >
            追加
          </button> */}
        </div>
      </div>
    </div>
  );
}