package main

import (
	"fmt"
	uuid "github.com/satori/go.uuid"

	"mxshop_api/user_web/global"
	"mxshop_api/user_web/initialize"
	"mxshop_api/user_web/utils"
	"mxshop_api/user_web/utils/register"
	mvalidator "mxshop_api/user_web/validator"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
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

	viper.AutomaticEnv()
	// 如果是本地开发环境端口号固定，线上环境启动获取端口号
	isDebug := viper.GetBool("MXSHOP_DEBUG")
	if !isDebug {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}

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

	// 服务注册
	registerClient := register.NewRegistryClient(
		global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port,
	)
	err = registerClient.Register(
		uuid.NewV4().String(),
		global.ServerConfig.Name,
		global.ServerConfig.Host,
		global.ServerConfig.Port,
		global.ServerConfig.Tags,
	)
	if err != nil {
		zap.S().Panicf("服务注册失败：%v", err)
	}

	zap.S().Debugf("启动服务器, 端口:%d", global.ServerConfig.Port)
	if err := r.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动服务失败", err)
	}
}
