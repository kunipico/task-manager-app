import { cookies } from "next/headers";
import { NextResponse } from "next/server";

export async function DELETE(request: Request) {
  try {
    // Cookieからトークンを取得
    const cookieStore = cookies();
    const token = cookieStore.get('token')?.value;
    console.log('token : ',token)

    if (!token) {
      // トークンが存在しない場合は401エラーを返す
      return NextResponse.json(
        { error: 'Unauthorized: No token found' },
        { status: 401 }
      );
    }

    const res = await fetch("http://localhost:8080/logout", {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
        'Cookie': `token=${token}`,
      },
      credentials: "include",
    });

    const data = await res.text();

    if (res.ok) {
      return NextResponse.json({ message: data }, { status: 200 });
    } else {
      return NextResponse.json({ message: data }, { status: res.status });
    }
  } catch (error) {
    return NextResponse.json(
      { message: "Internal Server Error" },
      { status: 500 }
    );
  }
}