package banner

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"google.golang.org/protobuf/types/known/emptypb"
	"mxshop_api/goods_web/api"
	"mxshop_api/goods_web/forms"
	"mxshop_api/goods_web/global"
	"mxshop_api/goods_web/global/response"
	"mxshop_api/goods_web/proto"
	"net/http"
)

func ListBanner(ctx *gin.Context) {
	resp, err := global.GoodsSrvClient.GetBannerList(context.Background(), &emptypb.Empty{})
	if err != nil {
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	bannerList := make([]*response.BannerInfoResp, 0, len(resp.Data))
	for _, data := range resp.Data {
		bannerList = append(bannerList, &response.BannerInfoResp{
			Id:    data.Id,
			Index: data.Index,
			Image: data.Image,
			Url:   data.Url,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"total":  resp.Total,
		"data":   bannerList,
	})
}

func CreateBanner(ctx *gin.Context) {
	form := &forms.BannerForm{}
	err := ctx.ShouldBind(&form)
	if err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	resp, err := global.GoodsSrvClient.CreateBanner(context.Background(), &proto.BannerInfoRequest{
		Index: form.Index,
		Image: form.Image,
		Url:   form.Url,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	data := &response.BannerInfoResp{
		Id:    resp.Id,
		Index: resp.Index,
		Image: resp.Image,
		Url:   resp.Url,
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   data,
	})
}

func DeleteBanner(ctx *gin.Context) {
	id := ctx.Param("id")

	_, err := global.GoodsSrvClient.DeleteBanner(context.Background(), &proto.BannerInfoRequest{
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

func UpdateBanner(ctx *gin.Context) {
	form := &forms.BannerForm{}
	err := ctx.ShouldBindJSON(form)
	if err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	_, err = global.GoodsSrvClient.UpdateBanner(context.Background(), &proto.BannerInfoRequest{
		Id:    cast.ToInt32(id),
		Index: form.Index,
		Image: form.Image,
		Url:   form.Url,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}
