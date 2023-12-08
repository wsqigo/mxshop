package api

import (
	"context"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"math/rand"
	"mxshop_api/user_web/forms"
	"mxshop_api/user_web/global"
	"net/http"
	"strings"
	"time"
)

// GenerateSmsCode 生成width长度的短信验证码
func GenerateSmsCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}

	return sb.String()
}

func SendSms(ctx *gin.Context) {
	form := forms.SendSmsForm{}
	if err := ctx.ShouldBind(&form); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	client, err := dysmsapi.NewClientWithAccessKey("cn-beijing",
		global.ServerConfig.AliSmsInfo.ApiKey, global.ServerConfig.AliSmsInfo.ApiSecret)
	if err != nil {
		zap.S().Infow("generate sms client failed", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "获取sms服务失败",
		})
		return
	}
	smsCode := GenerateSmsCode(6)
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https"
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionID"] = "cn-beijing"
	request.QueryParams["PhoneNumbers"] = form.Mobile // 手机号
	request.QueryParams["SignName"] = "极客在线"
	request.QueryParams["TemplateCode"] = "SMS_464041494"
	request.QueryParams["TemplateParam"] = fmt.Sprintf(`{"code":%s}`, smsCode)
	response, err := client.ProcessCommonRequest(request)
	fmt.Println(client.DoAction(request, response))
	if err != nil {
		zap.S().Errorw("generate sms code failed", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成验证码失败",
		})
		return
	}

	// 后面注册的时候会将短信验证码带回来
	// 需要将验证码保存起来

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host,
			global.ServerConfig.RedisInfo.Port),
	})
	rdb.Set(context.Background(), form.Mobile, smsCode, time.Duration(global.ServerConfig.RedisInfo.Expire)*time.Second)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})

}
