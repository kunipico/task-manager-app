'use client';

import { useEffect, useState } from 'react';

type TaskModalProps = {
  taskName:string;
  taskId: number;
  isOpen: boolean;
  onClose: () => void;
};

type TimeInfo = {
  today: string;
  thisWeek: string;
  thisMonth: string;
};

export default function TaskModal({taskName, taskId, isOpen, onClose }: TaskModalProps) {
  const [timeInfo, setTimeInfo] = useState<TimeInfo | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!isOpen) return;

    const fetchTimeInfo = async () => {
      setLoading(true);
      setError(null);

      console.log('taskId: ',taskId)

      try {
        const res = await fetch(`http://localhost:8080/tasks/time-info/${taskId}`, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: "include",
        });
        if (res.ok) {

        };

        const data = await res.json();
        setTimeInfo(data);
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchTimeInfo();
  }, [isOpen, taskId]);

  if (!isOpen) return null;

  const handleModalClick = (e: React.MouseEvent) => {
    e.stopPropagation(); // バブリングを防止
  };

  return (
    <div 
      className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center"
      onClick={onClose} // 背景クリックで閉じる
    >
      <div 
        className="bg-white p-6 rounded shadow-lg w-96"
        onClick={handleModalClick} // モーダル内クリックではバブリングを防ぐ
      >
        {loading && <p>読み込み中...</p>}
        {error && <p className="text-red-500">{error}</p>}
        {timeInfo && (
          <div className="text-center">
            <h2 className="text-2xl mb-4">
              「{taskName} 」実施時間
            </h2>
            <p className="mb-4">今日 : {timeInfo.today}</p>
            <p className="mb-4">今週 : {timeInfo.thisWeek}</p>
            <p className="mb-4">今月 : {timeInfo.thisMonth}</p>
          </div>
        )}
      </div>
    </div>
  );
}