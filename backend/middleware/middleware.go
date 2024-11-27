// Cookieの認証をこのミドルウェアにまとめる。
// 2024.11.17まだ利用できない。

package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

var secretKey = []byte("my_secret_key") // JWT署名用のシークレットキー

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CookieからJWTを取得
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Unauthorized: No cookie found", http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value

		// JWTの検証
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// トークンの署名アルゴリズムを検証
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// トークンのカスタムクレームから情報を取得（必要なら）
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			exp := int64(claims["exp"].(float64))
			if time.Now().Unix() > exp {
				http.Error(w, "Unauthorized: Token expired", http.StatusUnauthorized)
				return
			}

			// 必要ならリクエストコンテキストに情報を追加
			r.Header.Set("UserID", claims["sub"].(string))
		} else {
			http.Error(w, "Unauthorized: Invalid claims", http.StatusUnauthorized)
			return
		}

		// 次のハンドラーを呼び出す
		next.ServeHTTP(w, r)
	})
}