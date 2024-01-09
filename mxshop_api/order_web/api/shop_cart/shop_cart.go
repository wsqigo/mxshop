package shop_cart

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"mxshop_api/order_web/api"
	"mxshop_api/order_web/forms"
	"mxshop_api/order_web/global"
	"mxshop_api/order_web/global/response"
	"mxshop_api/order_web/proto"
	"net/http"
)

func ListCart(ctx *gin.Context) {
	// 获取购物车商品
	userId := ctx.GetInt("userId")

	resp, err := global.OrderSrvClient.GetCartItemList(context.Background(), &proto.UserInfo{
		Id: int32(userId),
	})
	if err != nil {
		zap.S().Errorw("查询购物车列表失败")
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	if resp.Total == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"total": 0,
		})
		return
	}

	ids := make([]int32, 0, resp.Total)
	for _, data := range resp.Data {
		ids = append(ids, data.GoodsId)
	}

	// 请求商品服务获取商品信息
	goodsResp, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: ids,
	})
	if err != nil {
		zap.S().Errorw("批量查询商品列表失败")
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	/*
		{
			"total":12,
			"data":[
				{
					"id":1,
					"goods_id":421,
					"goods_name":421,
					"goods_price":421,
					"goods_image":421,
					"nums":421,
					"checked": true,
				}
			]
		}
	*/

	goodsList := make([]*response.CartItem, 0, resp.Total)
	for _, cartItem := range resp.Data {
		for _, goods := range goodsResp.Data {
			if cartItem.GoodsId != goods.Id {
				continue
			}
			goodsList = append(goodsList, &response.CartItem{
				Id:         cartItem.Id,
				GoodsId:    goods.Id,
				GoodsName:  goods.Name,
				GoodsImage: goods.GoodsFrontImage,
				GoodsPrice: goods.ShopPrice,
				Nums:       cartItem.Nums,
				Checked:    cartItem.Checked,
			})
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"total":  resp.Total,
		"data":   goodsList,
	})
}

// AddCartItem 添加商品到购物车
func AddCartItem(ctx *gin.Context) {
	form := forms.ShopCartForm{}
	err := ctx.ShouldBind(&form)
	if err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	// 为了严谨性，添加商品到购物车之前，记得检查一下商品是否存在
	_, err = global.GoodsSrvClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{
		Id: form.GoodsId,
	})
	if err != nil {
		zap.S().Errorw("查询商品信息失败")
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	// 如果现在添加到购物车的数量和库存的数量不一致
	invResp, err := global.InventorySrvClient.GetGoodsInvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: form.GoodsId,
	})
	if err != nil {
		zap.S().Errorw("查询库存信息失败")
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	if invResp.Num < form.Nums {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "库存不足",
		})
		return
	}

	userId := ctx.GetInt("userId")
	resp, err := global.OrderSrvClient.CreateCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  int32(userId),
		GoodsId: form.GoodsId,
		Nums:    form.Nums,
	})
	if err != nil {
		zap.S().Errorw("添加到购物车失败")
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"id":     resp.Id,
	})
}

func DeleteCartItem(ctx *gin.Context) {
	id := ctx.Param("id")
	// 登录时，改成从 ctx.GetInt("userId")拿
	userId := ctx.Query("userId")

	_, err := global.OrderSrvClient.DeleteCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  cast.ToInt32(userId),
		GoodsId: cast.ToInt32(id),
	})
	if err != nil {
		zap.S().Errorw("删除购物车记录失败")
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func UpdateCartItem(ctx *gin.Context) {
	// o/v1/421
	id := ctx.Param("id")
	form := forms.ShopCartUpdateForm{}
	err := ctx.ShouldBind(&form)
	if err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	req := &proto.CartItemRequest{
		UserId:  form.UserId,
		GoodsId: cast.ToInt32(id),
		Nums:    form.Nums,
	}

	if form.Checked != nil {
		req.Checked = *form.Checked
	}

	_, err = global.OrderSrvClient.UpdateCartItem(context.Background(), req)
	if err != nil {
		zap.S().Errorw("更新购物车记录失败")
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}
