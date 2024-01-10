package message

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop_api/userop_web/api"
	"mxshop_api/userop_web/forms"
	"mxshop_api/userop_web/global"
	"mxshop_api/userop_web/global/response"
	"mxshop_api/userop_web/models"
	"mxshop_api/userop_web/proto"
	"net/http"
)

func ListMessage(ctx *gin.Context) {
	req := &proto.MessageRequest{}

	claims, _ := ctx.Get("claims")
	model := claims.(*models.CustomClaims)
	if model.AuthorityId == 1 {
		req.UserId = int32(model.ID)
	}

	resp, err := global.MessageSrvClient.GetMessageList(context.Background(), req)
	if err != nil {
		zap.S().Errorw("获取留言失败")
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	data := make([]*response.MessageResp, 0, resp.Total)
	for _, message := range resp.Data {
		data = append(data, &response.MessageResp{
			Id:      message.Id,
			UserId:  message.UserId,
			Type:    message.MessageType,
			Subject: message.Subject,
			Message: message.Message,
			File:    message.File,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total": resp.Total,
		"data":  data,
	})
}

func CreateMessage(ctx *gin.Context) {
	form := &forms.MessageForm{}
	err := ctx.ShouldBind(form)
	if err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	userId, _ := ctx.Get("userId")
	resp, err := global.MessageSrvClient.CreateMessage(context.Background(), &proto.MessageRequest{
		UserId:      int32(userId.(uint)),
		MessageType: form.MessageType,
		Subject:     form.Subject,
		Message:     form.Message,
		File:        form.File,
	})

	if err != nil {
		zap.S().Errorw("添加留言失败")
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"id":     resp.Id,
	})
}
