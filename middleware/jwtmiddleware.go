package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"zerorequest/errorx"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type JwtAuthMiddleware struct {
	secret string
}

func NewJwtAuthMiddleware(secret string) *JwtAuthMiddleware {
	return &JwtAuthMiddleware{
		secret: secret,
	}
}

func (m *JwtAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//1.获取token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httpx.Error(w, errorx.NewCodeError(401, "Authorization header is missing"))
		}
		//2.检查Bearer格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			httpx.Error(w, errorx.NewCodeError(401, "Authorization header format must be Bearer {token}"))

			return
		}
		//3.解析和验证token
		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			//验证签名方法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.secret), nil
		})

		if err != nil {
			httpx.Error(w, errorx.NewCodeError(401, "Invalid token"))
			return
		}
	}

}
