// RouteHandler サーバーコンポーネント　クライアントコンポーネントのAPIとして機能する。
// RouteHandler経由でデータfetchが一般的らしい。
// import { cookies } from 'next/headers';
// import {NextRequest, NextResponse } from 'next/server';

// const BACKEND_URL = 'http://go-api:8080/tasks';

// // GET: タスク一覧を取得
// export async function GET(request: NextRequest) {
//   console.log('router.ts!!!')
//   const cookieToken = cookies().get("token")
//   if(cookieToken){
//     console.log("Token from cookies: ", cookieToken.value); // クッキーの値をログに出力
//     cookies().set("token",cookieToken.value)
//   } else {
//     console.log("No token found in cookies.");
//   }
//   const res = await fetch(BACKEND_URL, {
//     method: 'GET',
//     headers: {
//       'Content-Type': 'application/json',
//     },
//     credentials: 'include', // Cookie付きでバックエンドへリクエスト
//   });

//   const data = await res.json();
//   console.log('data : ',data)
//   return NextResponse.json(data);
// }
