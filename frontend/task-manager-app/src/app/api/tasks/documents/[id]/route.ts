import { cookies } from 'next/headers';
import { NextRequest, NextResponse } from 'next/server';

export async function POST(req: NextRequest, { params }: { params: { id: string } }) {
  const apiUrl = `http://localhost:8080/tasks/documents/${params.id}`;
  try {
    // Cookieからトークンを取得
    const cookieStore = cookies();
    const token = cookieStore.get('token')?.value;
    console.log("RouteHandler docs!!!")

    if (!token) {
      // トークンが存在しない場合は401エラーを返す
      return NextResponse.json(
        { error: 'Unauthorized: No token found' },
        { status: 401 }
      );
    }

    const body = await req.json();
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
      return NextResponse.json({ error: 'Failed to create task documents' }, { status: res.status });
    }

    const createTaskDoc = res.json();
    console.log('createTaskDoc : ',createTaskDoc);
    return NextResponse.json(createTaskDoc, { status: 200 });
  } catch (error) {
    console.log('RouteHandler Internal Server Error');
    return NextResponse.json({ error: 'Internal Server Error' }, { status: 500 });
  }
}

export async function GET(req: NextRequest, { params }: { params: { id: string } }) {
  const apiUrl = `http://localhost:8080/tasks/documents/${params.id}`;
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
      return NextResponse.json({ error: 'Failed to get task documents' }, { status: res.status });
    }
    const taskDoc = await res.json();
    return NextResponse.json(taskDoc, { status: 200 });
  } catch (error) {
    return NextResponse.json({ error: 'Internal Server Error' }, { status: 500 });
  }
}