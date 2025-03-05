package utils

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// 配置项
var (
	jwtSecret     = []byte("d2f5b8e3a9c7d1e5f6b9a8c7d2e5f3b1") // 建议从配置文件读取
	tokenExpire   = 24 * time.Hour                             // Token有效期
	issuer        = "govue3"                                   // 签发者
	ErrEmptyToken = errors.New("token不能为空")
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

// GenerateToken 生成JWT Token
func GenerateToken(userID uint) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(tokenExpire)

	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  nowTime.Unix(),
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析和验证Token
func ParseToken(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, ErrEmptyToken
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("token格式错误")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("token已过期")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("token尚未生效")
			} else {
				return nil, errors.New("token验证失败")
			}
		}
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的token")
}

// RefreshToken 刷新Token
func RefreshToken(tokenString string) (string, error) {
	// 分割 "Authorization" 头部的值，提取出 token
	authHeader := strings.Split(tokenString, "Bearer ")
	if len(authHeader) != 2 {
		return "", errors.New("无效的Token格式")
	}
	token := authHeader[1]

	claims, err := ParseToken(token)
	if err != nil {
		return "", err
	}
	return GenerateToken(claims.UserID)
}
