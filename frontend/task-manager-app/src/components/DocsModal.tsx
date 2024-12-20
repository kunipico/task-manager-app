'use client';

import { useState, useEffect } from "react";

type TaskModalProps = {
  taskName:string;
  taskId: number;
  isOpen: boolean;
  onClose: () => void;
};

type DocInfo = {
  doc: string;
};

export default function DocsModal({taskName, taskId, isOpen, onClose }: TaskModalProps)
{
  const [documents, setDocuments] = useState<DocInfo[]>([]);
  const [newDoc, setNewDoc] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!isOpen) return;
    fetchDocInfo();
  },[isOpen, taskId]);

  if (!isOpen) return;
  
  const fetchDocInfo = async () =>{
    setLoading(true);
    setError(null);
    setDocuments([]);

    // APIから既存のドキュメントを取得
    try{
      const res = await fetch(`/api/tasks/documents/${taskId}`,{
        method : 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: "include",
      });
      if (!res.ok){
        throw new Error(`Failed to fetch documents: ${res.statusText}`);
      };

      const data = await res.json();
      console.log('Modal !!! afterdata : ',data);
      if (data != null) {
        const contents = data.map((doc: { content: string }) => ({doc: doc.content}));
        console.log('contents : ',contents);
        setDocuments(contents);
      } 
    }catch(err:any){
      setError(err.message || "An error occurred while fetching documents");
    }finally{
      setLoading(false);
    }
  };
 
  const handleAddDocument = async () => {
    if (!newDoc.trim()) {
      setError("Document content cannot be empty");
      return;
    }

    setLoading(true);
    setError(null);

    try{
    // 新しいドキュメントをAPIに送信 
    const res = await fetch(`/api/tasks/documents/${taskId}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
      body: JSON.stringify({ content: newDoc }),
    });
    if(!res.ok){
      throw new Error(`Failed to add document: ${res.statusText}`);
    }
    // const addedDocument = await res.json();
    //   setDocuments((prevDocs) => [...(prevDocs || []), addedDocument]);
      
    } catch (err: any) {
      setError(err.message || "An error occurred while adding the document");
    } finally {
      setLoading(false);
      setNewDoc(""); // フォームをリセット
      // モーダルを閉じる
      onClose();
    }
  };

  const handleModalClick = (e: React.MouseEvent) => {
    e.stopPropagation(); // バブリングを防止
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setNewDoc(e.target.value);
  };

  return (
    <div
      className="fixed inset-0 bg-gray-900 bg-opacity-50 flex justify-center items-center z-50"
      onClick={onClose}
    >
      <div
        className="bg-white rounded-lg shadow-lg p-6 w-11/12 max-w-lg"
        onClick={handleModalClick}
      >
        <h2 className="text-xl font-semibold mb-4 text-center">
          「{taskName}」
        </h2>
        {loading && <p className="text-center">読み込み中...</p>}
        {error && <p className="text-center text-red-500">{error}</p>}
  
        <div className="mb-6">
          <h3 className="text-md font-semibold mb-2">Add New Document</h3>
          <textarea
            className="w-full p-4 border-2 border-blue-500 rounded-lg focus:ring-4 focus:ring-blue-300 bg-blue-50 mb-4 text-gray-800"
            value={newDoc}
            onChange={handleInputChange}
            placeholder="Write your document here..."
            rows={5}
          />
          <div className="flex justify-center mb-6">
            <button
              className="bg-blue-500 text-white p-2 rounded-lg hover:bg-blue-600 transition font-semibold "
              onClick={handleAddDocument}
            >
              Add Document
            </button>
          </div>
        </div>
  
        <div>
          <h3 className="text-md font-semibold mb-2"> Documents</h3>
          {documents && documents.length > 0 ? (
            <div className="max-h-48 overflow-y-auto space-y-2">
              {documents.map((document, index) => (
                <div
                  key={index}
                  className="p-3 border rounded-md shadow-sm bg-gray-50 whitespace-pre-wrap"
                >
                  {document.doc.split('\n').map((line,i)=> (
                    <span key={i}>
                      {line}
                      <br />
                    </span>
                ))}
                </div>
              ))}
            </div>
          ) : (
            <p className="text-gray-500 text-sm text-center">
              No documents available.
            </p>
          )}
        </div>
      </div>
    </div>
  );
  
}