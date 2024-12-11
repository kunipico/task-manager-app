package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"task-manager-api/db"
	"time"

	"github.com/golang-jwt/jwt"
	// "github.com/labstack/gommon/email"
	// "github.com/pelletier/go-toml/query"
	"golang.org/x/crypto/bcrypt"
)

// ユーザー構造体
type User struct {
	ID   int    `json:"User_ID"`
	Name string `json:"User_Name"`
  Email string `json:"Emailaddress"`
  Password string `json:"Password"`
}

//認証情報
type Credentials struct {
	ID string `json:"ID"`
  Username string `json:"username"`
	Password string `json:"password"`
}

// JWTシークレットキー
var jwtKey = []byte("my_secret_key")

// JWTのペイロードを定義
type Claims struct {
	Userinfo string `json:"userinfo"`
	jwt.StandardClaims
}

// ユーザ登録処理
func Signup(w http.ResponseWriter, r *http.Request) {
	// リクエストボディからJSONデータを読み込む
	var requestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
  }

  // JSONをデコード
  err := json.NewDecoder(r.Body).Decode(&requestBody)
  if err != nil {
      http.Error(w, "無効なリクエストボディです", http.StatusBadRequest)
      log.Println("JSON Decode Error:", err)
      return
  }

  // ユーザー名とパスワードをフォームから取得
  username := requestBody.Username
  password := requestBody.Password
  email    := requestBody.Email
  fmt.Println("username: ",username)
  fmt.Println("password: ",password)
  fmt.Println("email: ",email)

  // バリデーションチェック
  if username == "" || email == "" || password == "" {
    http.Error(w, "ユーザー名とメールアドレス、パスワードは必須です", http.StatusBadRequest)
    return
  }
  // セキュアなバリデーション条件を追加する
  // メールアドレスのバリデーションを追加する

  // パスワードをハッシュ化
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
  if err != nil {
    http.Error(w, "サーバーエラーが発生しました", http.StatusInternalServerError)
    log.Println("Password Hashing Error:", err)
    return
  }

  // DBから情報取得
	var userID int
	err = db.DB.QueryRow("SELECT User_ID, Password FROM Users WHERE Emailaddress = ?", email).Scan(&userID)
	fmt.Println("DBからUser情報取得 err: ",err)
	if err == sql.ErrNoRows {
			// 重複するユーザーがいない場合の処理
			// ユーザー登録処理を進める
			fmt.Println("メールアドレスが未登録です。登録を進めます。")

			// ユーザー情報をデータベースに挿入
			query := "INSERT INTO Users (User_Name, Emailaddress, Password) VALUES (?, ?, ?)"
			result, err := db.DB.Exec(query, username, email, string(hashedPassword))
			if err != nil {
				http.Error(w, "ユーザーの登録に失敗しました", http.StatusInternalServerError)
				log.Println("DB Insert Error:", err)
				return
			}
			// 成功メッセージを返す
			// fmt.Fprintf(w, "ユーザー登録が完了しました。")

			// 新しいユーザーIDを取得
			newUserID, err := result.LastInsertId()
			fmt.Println("newUserID: ",newUserID)
			if err != nil {
				http.Error(w, "ユーザーIDの取得に失敗しました", http.StatusInternalServerError)
				log.Println("LastInsertId Error:", err)
				return
			}

			// JWTを生成し、クライアントに返す処理をここで実行
			signupToken, err := GenerateJWT(strconv.FormatInt(newUserID,10))
			if err != nil {
				fmt.Println("JWT生成エラー:", err)
				return
			}
			// クライアントにJWTを返す（例として標準出力）
			fmt.Println("JWT:", signupToken)
			// トークンをレスポンスとして返す
			// JWTをCookieにセット
			http.SetCookie(w, &http.Cookie{
				Name:     "token",
				Value:    signupToken,
				HttpOnly: true,
				Secure:   false, // HTTPSでのみ送信
				Path:     "/",
				SameSite: http.SameSiteLaxMode,
				MaxAge: 0, //Maxage:0でセッションが閉じられるまで。Maxage:NでN秒後まで。
			})
			// ログイン成功メッセージを送信
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("SignUp successful"))
	} else {
			http.Error(w, "既に登録済みのメールアドレスです", http.StatusBadRequest)
	}
}

//ユーザー登録確認用
func Users(w http.ResponseWriter, r *http.Request){
  rows, err := db.DB.Query("SELECT User_ID, User_Name, Emailaddress, Password FROM Users")
	if err != nil {
		http.Error(w, "データベースエラー", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		if err != nil {
			http.Error(w, "データ取得エラー", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

//ログイン処理
func Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
  // ユーザー名とパスワードをフォームから取得
  username_input := creds.Username
  password_input := creds.Password
  fmt.Println("username: ",username_input)
  fmt.Println("password: ",password_input)

  // DBから情報取得
  rows, err := db.DB.Query("SELECT User_ID, Password FROM Users WHERE User_Name=?",username_input)
	if err != nil {
		http.Error(w, "データベースエラー", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user_db User
		err := rows.Scan(&user_db.ID, &user_db.Password)
		if err != nil {
			http.Error(w, "データ取得エラー", http.StatusInternalServerError)
			return
		}
    fmt.Println("user.Password: ",user_db.Password)
    if CheckPasswordHash(password_input, user_db.Password){
      // 認証成功
      fmt.Println("Correct Password!!!")

      // JWTを生成し、クライアントに返す処理をここで実行
      token, err := GenerateJWT(strconv.Itoa(user_db.ID))
      if err != nil {
        fmt.Println("JWT生成エラー:", err)
        return
      }
      // クライアントにJWTを返す（例として標準出力）
      fmt.Println("JWT:", token)
      // トークンをレスポンスとして返す
      // JWTをCookieにセット
      http.SetCookie(w, &http.Cookie{
        Name:     "token",
        Value:    token,
        HttpOnly: true,
        Secure:   false, // HTTPSでのみ送信
        Path:     "/",
        SameSite: http.SameSiteLaxMode,
				MaxAge: 0, //Maxage:0でセッションが閉じられるまで。Maxage:NでN秒後まで。
      })
      // ログイン成功メッセージを送信
      w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
      w.Header().Set("Access-Control-Allow-Credentials", "true")
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(http.StatusOK)
      // w.Write([]byte("Login successful"))
      
      return
	  } else {
      // 認証失敗
      fmt.Println("Unauthorized")
      http.Error(w, "Unauthorized", http.StatusUnauthorized)
      return
	  }
  }
  fmt.Println("ユーザーが見つかりません")
  http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

// パスワードの比較関数
func CheckPasswordHash(password, hash string) bool {
  // ソルトを含めてハッシュ化しているため、下記のbcryptメソッドが必要
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// JWTを生成する関数
func GenerateJWT(userinfo string) (string, error) {
  expirationTime := time.Now().Add(7* 24 * time.Hour)
	claims := &Claims{
		Userinfo : userinfo,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	// JWTトークンを生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

// ログアウト処理
func Logout(w http.ResponseWriter, r *http.Request) {
	// Cookieの内容を確認
	cookie, err := r.Cookie("token")
	fmt.Println("Cookie : ", cookie)
  if err != nil {
    http.Error(w, `{"error":"Nothing Cookie!"}`, http.StatusUnauthorized)
		fmt.Println("Nothing Cookie!!!")
    return
  }
  tokenStr := cookie.Value
  // JWTの検証
  claims := &Claims{}
  token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
    return jwtKey, nil
  })
  if err != nil || !token.Valid {
    http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
    return
  }	

	// JWTをCookieにセット
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    cookie.Value,
		HttpOnly: true,
		Secure:   false, // HTTPSでのみ送信
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		MaxAge: -1, //Maxage:-Nで即時削除
	})
	// ログイン成功メッセージを送信
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	// w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000/login")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout successful"))
}
