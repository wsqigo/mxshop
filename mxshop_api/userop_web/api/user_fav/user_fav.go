package user_fav

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"mxshop_api/userop_web/api"
	"mxshop_api/userop_web/forms"
	"mxshop_api/userop_web/global"
	"mxshop_api/userop_web/global/response"
	"mxshop_api/userop_web/proto"
	"net/http"
)

func ListUserFav(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")
	userFavRsp, err := global.UserFavSrvClient.GetFavList(context.Background(), &proto.UserFavRequest{
		UserId: cast.ToInt32(userId),
	})
	if err != nil {
		zap.S().Errorw("获取收藏列表失败")
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	ids := make([]int32, 0, len(userFavRsp.Data))
	for _, item := range userFavRsp.Data {
		ids = append(ids, item.GoodsId)
	}

	//请求商品服务
	goods, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: ids,
	})
	if err != nil {
		zap.S().Errorw("[List] 批量查询【商品列表】失败")
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	goodsList := make([]*response.UserFavResp, 0, userFavRsp.Total)
	for _, item := range userFavRsp.Data {
		goodsItem := &response.UserFavResp{
			Id: item.GoodsId,
		}

		for _, good := range goods.Data {
			if item.GoodsId == good.Id {
				goodsItem.Name = good.Name
				goodsItem.ShopPrice = good.ShopPrice
			}
		}

		goodsList = append(goodsList, goodsItem)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   goodsList,
	})
}

func CreateUserFav(ctx *gin.Context) {
	form := forms.UserFavForm{}
	err := ctx.ShouldBindJSON(&form)
	if err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	// 查询商品id是否存在
	_, err = global.GoodsSrvClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{
		Id: form.GoodsId,
	})
	if err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	userId, _ := ctx.Get("userId")
	_, err = global.UserFavSrvClient.AddUserFav(context.Background(), &proto.UserFavRequest{
		UserId:  cast.ToInt32(userId),
		GoodsId: form.GoodsId,
	})
	if err != nil {
		zap.S().Errorw("添加收藏记录失败")
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func DeleteUserFav(ctx *gin.Context) {
	id := ctx.Param("id")

	userId, _ := ctx.Get("userId")
	_, err := global.UserFavSrvClient.DeleteUserFav(context.Background(), &proto.UserFavRequest{
		UserId:  cast.ToInt32(userId),
		GoodsId: cast.ToInt32(id),
	})
	if err != nil {
		zap.S().Errorw("删除收藏记录失败")
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"msg":    "删除成功",
	})
}

func GetUserFavDetail(ctx *gin.Context) {
	goodsId := ctx.Param("id")

	userId, _ := ctx.Get("userId")
	_, err := global.UserFavSrvClient.GetUserFavDetail(context.Background(), &proto.UserFavRequest{
		UserId:  cast.ToInt32(userId),
		GoodsId: cast.ToInt32(goodsId),
	})
	if err != nil {
		zap.S().Errorw("查询收藏状态失败")
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}
