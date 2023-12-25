package goods

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"mxshop_api/goods_web/api"
	"mxshop_api/goods_web/forms"
	"mxshop_api/goods_web/global"
	"mxshop_api/goods_web/global/response"
	"mxshop_api/goods_web/proto"
	"mxshop_api/goods_web/utils"
	"net/http"
)

type ListGoodsReq struct {
	PriceMin    int32  `json:"price_min"`
	PriceMax    int32  `json:"price_max"`
	IsHot       bool   `json:"is_hot"`
	IsNew       bool   `json:"is_new"`
	IsTab       bool   `json:"is_tab"`
	TopCategory int32  `json:"category_id"`
	Pages       int32  `json:"page_num"`
	PagePerNums int32  `json:"page_size"`
	KeyWords    string `json:"keywords"`
	Brand       int32  `json:"brand"`
}

func ConvertGoodsInfo2Response(info *proto.GoodsInfoResponse) *response.GoodsInfoResp {
	res := &response.GoodsInfoResp{
		Id:              info.Id,
		Name:            info.Name,
		ShopPrice:       info.ShopPrice,
		GoodsBrief:      info.GoodsBrief,
		GoodsDesc:       info.GoodsDesc,
		ShipFree:        info.ShipFree,
		Images:          info.Images,
		DescImages:      info.DescImages,
		GoodsFrontImage: info.GoodsFrontImage,
		IsNew:           info.IsHot,
		IsHot:           info.IsHot,
		OnSale:          info.OnSale,
		Category: &response.CategoryInfoResp{
			Id:   info.Category.Id,
			Name: info.Category.Name,
		},
		Brand: &response.BrandInfoResp{
			Id:   info.Brand.Id,
			Name: info.Brand.Name,
			Logo: info.Brand.Logo,
		},
	}

	return res
}

// ListGoods 商品列表
func ListGoods(ctx *gin.Context) {
	// todo: query也用结构体
	req := &ListGoodsReq{}
	err := utils.ShouldQueryParam(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	resp, err := global.GoodsSrvClient.GetGoodsList(context.Background(), &proto.GoodsFilterRequest{
		PriceMin:    req.PriceMin,
		PriceMax:    req.PriceMax,
		IsHot:       req.IsHot,
		IsNew:       req.IsNew,
		IsTab:       req.IsTab,
		TopCategory: req.TopCategory,
		Pages:       req.Pages,
		PagePerNums: req.PagePerNums,
		KeyWords:    req.KeyWords,
		Brand:       req.Brand,
	})
	if err != nil {
		zap.S().Errorf("查询商品列表失败: %v", err)
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	goodsList := make([]*response.GoodsInfoResp, 0, len(resp.Data))
	for _, data := range resp.Data {
		goodsList = append(goodsList, ConvertGoodsInfo2Response(data))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total": resp.Total,
		"data":  resp.Data,
	})
}

func CreateGoods(ctx *gin.Context) {
	form := &forms.GoodsForm{}
	if err := ctx.ShouldBind(&form); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	client := global.GoodsSrvClient
	resp, err := client.CreateGoods(context.Background(), &proto.CreateGoodsInfo{
		Name:            form.Name,
		GoodsSn:         form.GoodsSn,
		Stocks:          form.Stocks,
		MarketPrice:     form.MarketPrice,
		ShopPrice:       form.ShopPrice,
		GoodsBrief:      form.GoodsBrief,
		ShipFree:        form.ShipFree,
		Images:          form.Images,
		DescImages:      form.DescImages,
		GoodsFrontImage: form.GoodsSn,
		CategoryId:      form.CategoryId,
		BrandId:         form.Brand,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	//如何设置库存
	//TODO 商品的库存 - 分布式事务

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   resp,
	})
}

func Detail(ctx *gin.Context) {
	id := ctx.Param("id")

	resp, err := global.GoodsSrvClient.GetGoodsDetail(
		context.Background(),
		&proto.GoodInfoRequest{
			Id: cast.ToInt32(id),
		},
	)
	if err != nil {
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	res := ConvertGoodsInfo2Response(resp)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   res,
	})
}

func DeleteGoods(ctx *gin.Context) {
	id := ctx.Param("id")

	_, err := global.GoodsSrvClient.DeleteGoods(context.Background(), &proto.DeleteGoodsInfo{
		Id: cast.ToInt32(id),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func GetGoodsStocks(ctx *gin.Context) {
	_ = ctx.Param("id")

	// todo: 商品的库存
}

func UpdateStatus(ctx *gin.Context) {
	form := &forms.GoodsStatusForm{}
	err := ctx.ShouldBind(form)
	if err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	_, err = global.GoodsSrvClient.UpdateGoods(context.Background(), &proto.CreateGoodsInfo{
		Id:     cast.ToInt32(id),
		IsHot:  form.IsHot,
		IsNew:  form.IsNew,
		OnSale: form.OnSale,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"msg":    "修改成功",
	})
}

func UpdateGoods(ctx *gin.Context) {
	form := &forms.GoodsForm{}
	err := ctx.ShouldBind(form)
	if err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param(":id")
	_, err = global.GoodsSrvClient.UpdateGoods(context.Background(), &proto.CreateGoodsInfo{
		Id:              cast.ToInt32(id),
		Name:            form.Name,
		GoodsSn:         form.GoodsSn,
		Stocks:          form.Stocks,
		MarketPrice:     form.MarketPrice,
		ShopPrice:       form.ShopPrice,
		GoodsBrief:      form.GoodsBrief,
		ShipFree:        form.ShipFree,
		Images:          form.Images,
		DescImages:      form.DescImages,
		GoodsFrontImage: form.FrontImage,
		CategoryId:      form.CategoryId,
		BrandId:         form.Brand,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"msg":    "更新成功",
	})
}
