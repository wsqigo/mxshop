package main

import (
	"fmt"

	"mxshop_api/user_web/global"
	"mxshop_api/user_web/initialize"
	mvalidator "mxshop_api/user_web/validator"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func main() {
	// 1. 初始化logger
	initialize.InitLogger()

	// 2. 初始化配置文件
	initialize.InitConfig()

	// 3. 初始化routers
	r := initialize.Routers()

	// 4. 初始化翻译
	err := initialize.InitTrans("zh")
	if err != nil {
		zap.S().Panic("init translator failed", err)
	}

	// 5. 初始化srv的链接
	initialize.InitSrvConn()

	// 注册验证器
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		_ = v.RegisterValidation("mobile", mvalidator.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码!", true) // see universal-translation
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	zap.S().Debugf("启动服务器, 端口:%d", global.ServerConfig.Port)
	if err := r.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动失败", err)
	}
}
