package order

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"mxshop_api/order_web/api"
	"mxshop_api/order_web/forms"
	"mxshop_api/order_web/global"
	"mxshop_api/order_web/global/response"
	"mxshop_api/order_web/models"
	"mxshop_api/order_web/proto"
	"net/http"
	"strconv"
)

func ListOrder(ctx *gin.Context) {
	// 订单的列表
	userId, _ := ctx.Get("userId")
	claims, _ := ctx.Get("claims")

	req := &proto.OrderFilterRequest{}
	// 如果是管理员用户则返回所有的订单
	model := claims.(*models.CustomClaims)
	if model.AuthorityId == 1 {
		req.UserId = cast.ToInt32(userId)
	}

	pages := ctx.DefaultQuery("pNum", "0")
	pageSize := ctx.DefaultQuery("pSize", "0")
	req.Pages = cast.ToInt32(pages)
	req.PagePerNums = cast.ToInt32(pageSize)

	resp, err := global.OrderSrvClient.GetOrderList(context.Background(), req)
	if err != nil {
		zap.S().Errorw("获取订单列表失败")
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	/*
		{
			"total":100,
			"data":[
				{
					"
				}
			]
		}
	*/

	orderList := make([]*response.OrderItem, 0, resp.Total)
	for _, data := range resp.Data {
		orderList = append(orderList, &response.OrderItem{
			Id:      data.Id,
			Status:  data.Status,
			PayType: data.PayType,
			User:    data.UserId,
			Post:    data.Post,
			Total:   data.Total,
			Address: data.Address,
			Name:    data.Name,
			Mobile:  data.Mobile,
			OrderSn: data.OrderSn,
			AddTime: data.AddTime,
		})
	}

	// 生成支付宝的支付url
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   orderList,
	})
}

func CreateOrder(ctx *gin.Context) {
	form := &forms.OrderForm{}
	err := ctx.ShouldBind(form)
	if err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	userId, _ := ctx.Get("userId")
	resp, err := global.OrderSrvClient.CreateOrder(context.Background(), &proto.OrderRequest{
		UserId:  cast.ToInt32(userId),
		Address: form.Address,
		Name:    form.Name,
		Mobile:  form.Mobile,
		Post:    form.Post,
	})
	if err != nil {
		zap.S().Errorw("新建订单失败")
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	// 生成支付宝的支付url
	alipayInfo := global.ServerConfig.AlipayInfo
	client, err := alipay.New(alipayInfo.AppId, alipayInfo.PrivateKey, false)
	if err != nil {
		zap.S().Errorw("实例化支付宝失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}

	err = client.LoadAliPayPublicKey(alipayInfo.AliPublicKey)
	if err != nil {
		zap.S().Errorw("加载支付宝的公钥失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}

	p := alipay.TradePagePay{}
	p.NotifyURL = alipayInfo.NotifyUrl
	p.ReturnURL = alipayInfo.ReturnUrl
	p.Subject = "极客生鲜订单-" + resp.OrderSn
	p.OutTradeNo = resp.OrderSn
	p.TotalAmount = strconv.FormatFloat(resp.Total, 'f', 2, 64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url, err := client.TradePagePay(p)
	if err != nil {
		zap.S().Errorw("生成支付url失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":     "ok",
		"id":         resp.Id,
		"alipay_url": url.String(),
	})
}

func GetCartDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	userId, _ := ctx.Get("userId")

	// 如果是管理员用户则返回所有的订单
	request := &proto.OrderRequest{
		Id: cast.ToInt32(id),
	}

	claims, _ := ctx.Get("claims")
	model := claims.(*models.CustomClaims)
	if model.AuthorityId == 1 {
		request.UserId = cast.ToInt32(userId)
	}

	resp, err := global.OrderSrvClient.GetOrderDetail(context.Background(), request)
	if err != nil {
		zap.S().Errorw("获取订单详情失败")
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	orderDetail := &response.OrderDetailItem{
		OrderItem: response.OrderItem{
			Id:      resp.OrderInfo.Id,
			Status:  resp.OrderInfo.Status,
			PayType: resp.OrderInfo.PayType,
			User:    resp.OrderInfo.UserId,
			Post:    resp.OrderInfo.Post,
			Total:   resp.OrderInfo.Total,
			Address: resp.OrderInfo.Address,
			Name:    resp.OrderInfo.Name,
			Mobile:  resp.OrderInfo.Mobile,
			OrderSn: resp.OrderInfo.OrderSn,
			AddTime: resp.OrderInfo.AddTime,
		},
	}

	goodsList := make([]*response.GoodsItem, 0, len(resp.GoodsItems))
	for _, item := range resp.GoodsItems {
		goodsList = append(goodsList, &response.GoodsItem{
			Id:    item.GoodsId,
			Name:  item.GoodsName,
			Image: item.GoodsImage,
			Price: item.GoodsPrice,
			Nums:  item.Nums,
		})
	}
	orderDetail.GoodsItems = goodsList

	// 生成支付宝的支付url
	alipayInfo := global.ServerConfig.AlipayInfo
	client, err := alipay.New(alipayInfo.AppId, alipayInfo.PrivateKey, false)
	if err != nil {
		zap.S().Errorw("实例化支付宝失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}

	err = client.LoadAliPayPublicKey(alipayInfo.AliPublicKey)
	if err != nil {
		zap.S().Errorw("加载支付宝的公钥失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}

	p := alipay.TradePagePay{}
	p.NotifyURL = alipayInfo.NotifyUrl
	p.ReturnURL = alipayInfo.ReturnUrl
	p.Subject = "极客生鲜订单-" + resp.OrderInfo.OrderSn
	p.OutTradeNo = resp.OrderInfo.OrderSn
	p.TotalAmount = strconv.FormatFloat(resp.OrderInfo.Total, 'f', 2, 64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url, err := client.TradePagePay(p)
	if err != nil {
		zap.S().Errorw("生成支付url失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":     "ok",
		"data":       orderDetail,
		"alipay_url": url.String(),
	})
}
