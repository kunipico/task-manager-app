// export const login = async (username: string, password: string) => {
//   try {
//     const res = await fetch("http://localhost:8080/login", {
//       method: "POST",
//       headers: {
//         "Content-Type": "application/json",
//       },
//       body: JSON.stringify({ username, password }),
//       credentials: "include",
//     });

//     const data = await res.text();
//     return data;
//   } catch (error) {
//     throw new Error("ログインに失敗しました");
//   }
// };

// export async function signup(username: string, email: string, password: string) {
//   console.log('username: ',username);
//   console.log('email: ',email);
//   console.log('password: ',password);

//   try {
//     const res = await fetch("http://localhost:8080/signup", {
//         method: "POST",
//         headers: {
//             "Content-Type": "application/json",
//         },
//         body: JSON.stringify({ username, email, password }),
//     });
//     const data = await res.text();
//     return data;
//   } catch (error) {
//     if (error instanceof Error) {
//       console.log('api-error1: ',error);
//       throw new Error(error.message);  // エラーメッセージをセット
//     } else {
//       console.log('api-error2: ',error);
//       throw new Error("予期しないエラーが発生しました");
//     }
//   }
// }