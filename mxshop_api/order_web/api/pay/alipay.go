package pay

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"
	"mxshop_api/order_web/api"
	"mxshop_api/order_web/global"
	"mxshop_api/order_web/proto"
	"net/http"
)

func Notify(ctx *gin.Context) {
	// 支付宝回调通知
	aliPayInfo := global.ServerConfig.AlipayInfo
	client, err := alipay.New(aliPayInfo.AppId, aliPayInfo.PrivateKey, false)
	if err != nil {
		zap.S().Errorw("实例化支付宝失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}

	err = client.LoadAliPayPublicKey(aliPayInfo.AliPublicKey)
	if err != nil {
		zap.S().Errorw("加载支付宝的公钥失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}

	// DecodeNotification 内部已调用 VerifySign 方法验证签名
	noti, err := client.DecodeNotification(ctx.Request.Form)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	_, err = global.OrderSrvClient.UpdateOrderStatus(context.Background(), &proto.OrderStatus{
		OrderSn: noti.OutTradeNo, // 平台订单号，自己生成的
		Status:  string(noti.TradeStatus),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	ctx.String(http.StatusOK, "success")
}
