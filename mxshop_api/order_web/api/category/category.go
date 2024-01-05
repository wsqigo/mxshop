package category

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"mxshop_api/order_web/api"
	"mxshop_api/order_web/forms"
	"mxshop_api/order_web/global"
	"mxshop_api/order_web/global/response"
	"mxshop_api/order_web/proto"
	"net/http"
)

func ConvertCategoryInfo2Resp(data *proto.CategoryInfoResponse) *response.CategoryInfoResp {
	return &response.CategoryInfoResp{
		Id:               data.Id,
		Name:             data.Name,
		ParentCategoryId: data.ParentCategory,
		Level:            data.Level,
		IsTab:            data.IsTab,
	}
}

func ListCategory(ctx *gin.Context) {
	resp, err := global.GoodsSrvClient.GetAllCategoryList(context.Background(), &emptypb.Empty{})
	if err != nil {
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	var data []*response.CategoryInfoResp
	err = json.Unmarshal([]byte(resp.JsonData), &data)
	if err != nil {
		zap.S().Errorf("查询商品分类失败: %v", err.Error())
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   data,
	})
}

func CreateCategory(ctx *gin.Context) {
	form := &forms.CategoryForm{}
	err := ctx.ShouldBind(&form)
	if err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	resp, err := global.GoodsSrvClient.CreateCategory(context.Background(), &proto.CategoryInfoRequest{
		Name:  form.Name,
		IsTab: form.IsTab,
		Level: form.Level,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	data := ConvertCategoryInfo2Resp(resp)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   data,
	})
}

func GetCategoryDetail(ctx *gin.Context) {
	id := ctx.Param("id")

	resp, err := global.GoodsSrvClient.GetSubCategory(context.Background(), &proto.CategoryListRequest{
		Id: cast.ToInt32(id),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	subCategoryList := make([]*response.CategoryInfoResp, 0, len(resp.SubCategoryList))
	for _, category := range resp.SubCategoryList {
		subCategoryList = append(subCategoryList, &response.CategoryInfoResp{
			Id:               category.Id,
			Name:             category.Name,
			Level:            category.Level,
			ParentCategoryId: category.ParentCategory,
			IsTab:            category.IsTab,
		})
	}
	res := response.CategoryInfoResp{
		Id:               resp.Info.Id,
		Name:             resp.Info.Name,
		Level:            resp.Info.Level,
		IsTab:            resp.Info.IsTab,
		ParentCategoryId: resp.Info.ParentCategory,
		SubCategoryList:  subCategoryList,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   res,
	})
}

func DeleteCategory(ctx *gin.Context) {
	id := ctx.Param("id")

	// todo:
	// 1. 先查询出该分类写的所有子分类
	// 2. 将所有的分类全部逻辑删除
	// 3. 将该分类下的所有的商品逻辑删除
	_, err := global.GoodsSrvClient.DeleteCategory(context.Background(), &proto.DeleteCategoryRequest{
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

func UpdateCategory(ctx *gin.Context) {
	form := &forms.UpdateCategoryForm{}
	err := ctx.ShouldBindJSON(form)
	if err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	_, err = global.GoodsSrvClient.UpdateCategory(context.Background(), &proto.CategoryInfoRequest{
		Id:   cast.ToInt32(id),
		Name: form.Name,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}
