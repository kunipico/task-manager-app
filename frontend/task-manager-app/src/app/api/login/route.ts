import { NextResponse } from "next/server";

export async function POST(request: Request) {
  try {
    const body = await request.json();

    const res = await fetch("http://localhost:8080/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
      credentials: "include",
    });

    // Go API のレスポンスから Set-Cookie ヘッダーを取得
    const setCookieHeader = res.headers.get("Set-Cookie");

    const data = await res.text();
 
    if (res.ok && setCookieHeader) {
      // Set-Cookie をレスポンスヘッダーとして設定
      const response = NextResponse.json({ message: data }, { status: 200 });
      response.headers.set("Set-Cookie", setCookieHeader); // 返却用のSet-Cookieを設定
      return response;
    } else {
      return NextResponse.json({ message: data }, { status: res.status });
    }
  } catch (error) {
    console.error("Error during login:", error);
    return NextResponse.json(
      { message: "Internal Server Error" },
      { status: 500 }
    );
  }
}