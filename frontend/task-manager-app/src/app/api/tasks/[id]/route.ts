import { cookies } from 'next/headers';
import { NextRequest, NextResponse } from 'next/server';

export async function PUT(req: NextRequest, { params }: { params: { id: string } }) {
  const apiUrl = `http://localhost:8080/tasks/${params.id}`;
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

    const body = await req.json();
    const res = await fetch(apiUrl, {
      method: 'PUT',
      headers: { 
        'Content-Type': 'application/json',
        'Cookie': `token=${token}`,
      },
      credentials: 'include',
      body: JSON.stringify(body),
    });

    if (!res.ok) {
      return NextResponse.json({ error: 'Failed to update task' }, { status: res.status });
    }

    return NextResponse.json({message:"task status changed"}, { status: 200 });
  } catch (error) {
    return NextResponse.json({ error: 'Internal Server Error' }, { status: 500 });
  }
}

export async function DELETE(req: NextRequest, { params }: { params: { id: string } }) {
  const apiUrl = `http://localhost:8080/tasks/${params.id}`;
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
      method: 'DELETE',
      headers: { 
        'Content-Type': 'application/json',
        'Cookie': `token=${token}`,
      },
      credentials: 'include',
    });

    if (!res.ok) {
      return NextResponse.json({ error: 'Failed to delete task' }, { status: res.status });
    }

    return NextResponse.json({ message: 'Task deleted' }, { status: 200 });
  } catch (error) {
    return NextResponse.json({ error: 'Internal Server Error' }, { status: 500 });
  }
}