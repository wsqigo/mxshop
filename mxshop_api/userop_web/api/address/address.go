package address

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"mxshop_api/userop_web/api"
	"mxshop_api/userop_web/forms"
	"mxshop_api/userop_web/global"
	"mxshop_api/userop_web/global/response"
	"mxshop_api/userop_web/models"
	"mxshop_api/userop_web/proto"
	"net/http"
)

func ListAddress(ctx *gin.Context) {
	req := &proto.AddressRequest{}

	claims, _ := ctx.Get("claims")
	model := claims.(*models.CustomClaims)
	if model.AuthorityId == 1 {
		req.UserId = int32(model.ID)
	}

	resp, err := global.AddressSrvClient.GetAddressList(context.Background(), req)
	if err != nil {
		zap.S().Errorw("获取地址列表失败")
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	data := make([]*response.AddressResp, 0, resp.Total)
	for _, message := range resp.Data {
		data = append(data, &response.AddressResp{
			Id:           message.Id,
			UserId:       message.UserId,
			Province:     message.Province,
			City:         message.City,
			District:     message.District,
			Address:      message.Address,
			SignerName:   message.SignerName,
			SignerMobile: message.SignerMobile,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total": resp.Total,
		"data":  data,
	})
}

func CreateAddress(ctx *gin.Context) {
	form := &forms.AddressForm{}
	err := ctx.ShouldBind(form)
	if err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	userId, _ := ctx.Get("userId")
	resp, err := global.AddressSrvClient.CreateAddress(context.Background(), &proto.AddressRequest{
		UserId:       int32(userId.(uint)),
		Province:     form.Province,
		City:         form.City,
		District:     form.District,
		Address:      form.Address,
		SignerName:   form.SignerName,
		SignerMobile: form.SignerMobile,
	})

	if err != nil {
		zap.S().Errorw("添加地址失败")
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"id":     resp.Id,
	})
}

func DeleteAddress(ctx *gin.Context) {
	id := ctx.Param("id")

	_, err := global.AddressSrvClient.DeleteAddress(context.Background(), &proto.AddressRequest{
		Id: cast.ToInt32(id),
	})
	if err != nil {
		zap.S().Errorw("删除地址失败")
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"msg":    "删除成功",
	})
}

func UpdateAddress(ctx *gin.Context) {
	form := forms.AddressForm{}
	err := ctx.ShouldBindJSON(&form)
	if err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	_, err = global.AddressSrvClient.UpdateAddress(context.Background(), &proto.AddressRequest{
		Id:           cast.ToInt32(id),
		Province:     form.Province,
		City:         form.City,
		District:     form.District,
		Address:      form.Address,
		SignerName:   form.SignerName,
		SignerMobile: form.SignerMobile,
	})
	if err != nil {
		zap.S().Errorw("更新地址失败")
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"msg":    "更新地址成功",
	})
}
