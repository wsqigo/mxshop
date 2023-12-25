package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"mxshop_api/user_web/forms"
	"mxshop_api/user_web/global"
	"mxshop_api/user_web/global/response"
	"mxshop_api/user_web/middlewares"
	"mxshop_api/user_web/models"
	"mxshop_api/user_web/proto"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleGrpcErrorToHttp(ctx *gin.Context, err error) {
	// 将grpc的code转换成http的状态码
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误",
				})
			}
		}
	}
}

func HandleValidatorError(ctx *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
}

func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}

func GetUserList(ctx *gin.Context) {
	pNum := ctx.DefaultQuery("pageNum", "1")
	pSize := ctx.DefaultQuery("pageSize", "10")

	rsp, err := global.UserSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		PNum:  cast.ToUint32(pNum),
		PSize: cast.ToUint32(pSize),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询 【用户列表】 失败")
		HandleGrpcErrorToHttp(ctx, err)
		return
	}

	result := make([]response.UserInfoResp, 0)
	for _, value := range rsp.Data {
		data := response.UserInfoResp{
			ID:       value.Id,
			NickName: value.NickName,
			Birthday: response.JsonTime(time.Unix(value.Birthday, 0)),
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		}

		result = append(result, data)
	}

	ctx.JSON(http.StatusOK, result)
}

func PasswordLogin(ctx *gin.Context) {
	// 表单验证
	form := forms.PasswordLoginForm{}
	if err := ctx.ShouldBindJSON(&form); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	// 验证码验证
	if !store.Verify(form.CaptchaId, form.Captcha, true) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "captcha err",
		})
		return
	}

	// 登录逻辑
	rsp, err := global.UserSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: form.Mobile,
	})
	if err != nil {
		HandleGrpcErrorToHttp(ctx, err)
		return
	}

	// 只是查询到了用户了而已，并没有检查密码
	check, err := global.UserSrvClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
		Password:          form.Password,
		EncryptedPassword: rsp.Password, // 通过手机号查询的密码
	})
	if err != nil {
		HandleGrpcErrorToHttp(ctx, err)
		return
	}

	if !check.Success {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "密码错误",
		})
		return
	}

	// 生成 token
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint(rsp.Id),
		NickName:    rsp.NickName,
		AuthorityId: rsp.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 60*60*24*30,
			Issuer:    "wsqigo",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":         rsp.Id,
		"nick_name":  rsp.NickName,
		"token":      token,
		"expired_at": claims.ExpiresAt * 1000,
	})
}

func Register(ctx *gin.Context) {
	form := forms.RegisterForm{}
	if err := ctx.ShouldBind(&form); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	// 验证码校验
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host,
			global.ServerConfig.RedisInfo.Port),
	})
	val, err := rdb.Get(context.Background(), form.Mobile).Result()
	if err != nil {
		if err == redis.Nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code": "验证码错误",
			})
			return
		}
		HandleValidatorError(ctx, err)
		return
	}

	if val != form.Code {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": "验证码错误",
		})
		return
	}

	user, err := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		NickName: form.Mobile,
		Password: form.Password,
		Mobile:   form.Mobile,
	})
	if err != nil {
		zap.S().Errorw("create user failed", "err", err)
		HandleGrpcErrorToHttp(ctx, err)
		return
	}

	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint(user.Id),
		NickName:    user.NickName,
		AuthorityId: user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(30 * 24 * time.Hour).UnixMilli(),
			Issuer:    "wsqigo",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":         user.Id,
		"nick_name":  user.NickName,
		"token":      token,
		"expired_at": time.Now().Add(30 * 24 * time.Hour).UnixMilli(),
	})
}
