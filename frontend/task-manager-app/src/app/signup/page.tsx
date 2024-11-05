"use client";

import { useState } from "react";
// import { signup } from "@/lib/api";
import Link from "next/link";

export default function SignupPage() {
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [message, setMessage] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const res = await fetch("http://localhost:8080/signup", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, email, password }),
      });
      const data = await res.text();
      setMessage(data);
    } catch (error) {
      if (error instanceof Error) {
        setMessage(error.message);  // エラーメッセージをセット
      }
    }
  }
  // try {
  //     const response = await signup(username, email, password);
  //     setMessage(response);
  //     console.log('message1: ',message);
  //   } catch (error) {
  //     if (error instanceof Error) {
  //       setMessage(error.message);  // エラーメッセージをセット
  //       console.log('error1: ',error);
  //     } else {
  //       setMessage("予期しないエラーが発生しました");
  //       console.log('error2: ',error);
  //     }
  //     console.log('message2: ',message);
  //   }
  // };

  return (
    <div className="flex justify-center items-center min-h-screen bg-gray-100">
      <div className="bg-white p-8 rounded-lg shadow-md w-full max-w-md">
        <h2 className="text-2xl font-bold text-center mb-6">ユーザー登録</h2>
        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label htmlFor="username" className="block text-gray-700 mb-2">ユーザー名:</label>
            <input
              type="text"
              id="username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            </div>
            <div>
              <label htmlFor="username" className="block text-gray-700 mb-2">メールアドレス:</label>
              <input
                type="text"
                id="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div>
              <label htmlFor="password" className="block text-gray-700 mb-2">パスワード:</label>
            <input
              type="password"
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <button type="submit" className="w-full bg-blue-500 text-white py-2 rounded-lg hover:bg-blue-600 transition">
          登録
          </button>
        </form>
        {message && <p className="text-center text-red-500 mt-4">{message}</p>}
        {!message && <p className="invisible text-center text-red-500 mt-4">{message}</p>}
        <p className="text-center text-gray-600 mt-4">
          <Link href="/login" className="text-blue-500 hover:underline">ログインページに戻る</Link>
        </p>
      </div>
    </div>
  );
}