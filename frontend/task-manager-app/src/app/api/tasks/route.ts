import { NextRequest, NextResponse } from 'next/server';
import { cookies } from 'next/headers';

// タスク一覧取得 (GETリクエスト)
export async function GET(req: NextRequest) {
  const apiUrl = 'http://localhost:8080/tasks';
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

    const res = await fetch(apiUrl, {
      method: 'GET',
      headers: { 
        'Content-Type': 'application/json',
        'Cookie': `token=${token}`,
      },
      credentials: 'include',
    });

    if (!res.ok) {
      return NextResponse.json({ error:res.body }, { status: res.status });
    }

    const tasks = await res.json();
    return NextResponse.json(tasks, { status: 200 });
  } catch (error) {
    return NextResponse.json({ error: 'Internal Server Error' }, { status: 500 });
  }
}

// タスク追加 (POSTリクエスト)
export async function POST(req: NextRequest) {
  const apiUrl = 'http://localhost:8080/tasks';
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

    const body = await req.json(); // リクエストボディ取得

    const res = await fetch(apiUrl, {
      method: 'POST',
      headers: { 
        'Content-Type': 'application/json',
        'Cookie': `token=${token}`,
      },
      credentials: 'include',
      body: JSON.stringify(body),
    });

    if (!res.ok) {
      return NextResponse.json({ error: 'Failed to add task' }, { status: res.status });
    }

    const newTask = await res.json();
    return NextResponse.json(newTask, { status: 201 });
  } catch (error) {
    return NextResponse.json({ error: 'Internal Server Error' }, { status: 500 });
  }
}