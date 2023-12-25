package middlewares

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"mxshop_api/goods_web/global"
	"mxshop_api/goods_web/models"
	"net/http"
	"time"
)

func JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 我们这里 jwt 鉴权取头部信息 x-token
		// 登录时会返回 token 信息
		// 这里前端需要把 token 存储到 cookie 或者本地 localStorage 中，不过需要跟后端协商过期时间
		// 可以约定刷新令牌或者重新登录
		token := ctx.Request.Header.Get("x-token")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"msg": "请登录",
			})
			ctx.Abort()
			return
		}
		j := NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"msg": "授权已过期",
				})
				ctx.Abort()
				return
			}
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"msg": "请登录",
			})
			ctx.Abort()
			return
		}

		ctx.Set("claims", claims)
		ctx.Set("userId", claims.ID)
		ctx.Next()
	}
}

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")
)

func NewJWT() *JWT {
	return &JWT{
		[]byte(global.ServerConfig.JWTInfo.SigningKey),
	}
}

// CreateToken 创建一个token
func (j *JWT) CreateToken(claims models.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// ParseToken 解析 token
func (j *JWT) ParseToken(tokenString string) (*models.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}

	if token != nil {
		if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
			return claims, nil
		}
	}
	return nil, TokenInvalid
}

// RefreshToken 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}
