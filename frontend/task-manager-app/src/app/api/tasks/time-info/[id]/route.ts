import { cookies } from 'next/headers';
import { NextRequest, NextResponse } from 'next/server';

export async function GET(req: NextRequest, { params }: { params: { id: string } }) {
  const apiUrl = `http://localhost:8080/tasks/time-info/${params.id}`;
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
      return NextResponse.json({ error: 'Failed to get task time-info' }, { status: res.status });
    }
    const timeInfo = await res.json();
    return NextResponse.json(timeInfo, { status: 200 });
  } catch (error) {
    return NextResponse.json({ error: 'Internal Server Error' }, { status: 500 });
  }
}