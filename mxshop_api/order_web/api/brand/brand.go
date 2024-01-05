package brand

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"mxshop_api/order_web/api"
	"mxshop_api/order_web/api/category"
	"mxshop_api/order_web/forms"
	"mxshop_api/order_web/global"
	"mxshop_api/order_web/global/response"
	"mxshop_api/order_web/proto"
	"net/http"
)

func ConvertBrandInfo2Response(data *proto.BrandInfoResponse) *response.BrandInfoResp {
	return &response.BrandInfoResp{
		Id:   data.Id,
		Name: data.Name,
		Logo: data.Logo,
	}
}

func GetBrandList(ctx *gin.Context) {
	// todo: query也用结构体
	pNum := ctx.DefaultQuery("page_num", "0")
	pSize := ctx.DefaultQuery("page_size", "0")

	resp, err := global.GoodsSrvClient.GetBrandList(context.Background(), &proto.BrandFilterRequest{
		Pages:       cast.ToInt32(pNum),
		PagePerNums: cast.ToInt32(pSize),
	})
	if err != nil {
		zap.S().Errorf("查询品牌列表失败: %v", err)
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	brandList := make([]*response.BrandInfoResp, 0, len(resp.Data))
	for _, data := range resp.Data {
		brandList = append(brandList, ConvertBrandInfo2Response(data))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"total":  resp.Total,
		"data":   resp.Data,
	})
}

func CreateBrand(ctx *gin.Context) {
	form := &forms.BrandForm{}
	err := ctx.ShouldBind(&form)
	if err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	resp, err := global.GoodsSrvClient.CreateBrand(context.Background(), &proto.BrandInfoRequest{
		Name: form.Name,
		Logo: form.Logo,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	data := ConvertBrandInfo2Response(resp)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   data,
	})
}

func DeleteBrand(ctx *gin.Context) {
	id := ctx.Param("id")

	_, err := global.GoodsSrvClient.DeleteBrand(context.Background(), &proto.BrandInfoRequest{
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

func UpdateBrand(ctx *gin.Context) {
	form := &forms.BrandForm{}
	err := ctx.ShouldBindJSON(form)
	if err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	_, err = global.GoodsSrvClient.UpdateBrand(context.Background(), &proto.BrandInfoRequest{
		Id:   cast.ToInt32(id),
		Name: form.Name,
		Logo: form.Logo,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

//////// 品牌分类 //////

// GetCategoryBrandList 获取分类下的品牌数据
func GetCategoryBrandList(ctx *gin.Context) {
	id := ctx.Param("id")
	resp, err := global.GoodsSrvClient.GetCategoryBrandList(
		context.Background(),
		&proto.CategoryInfoRequest{Id: cast.ToInt32(id)},
	)
	if err != nil {
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	brandInfoList := make([]*response.BrandInfoResp, 0, len(resp.Data))
	for _, data := range resp.Data {
		brandInfoList = append(brandInfoList, ConvertBrandInfo2Response(data))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   brandInfoList,
	})
}

func CategoryBrandList(ctx *gin.Context) {
	//所有的list返回的数据结构
	/*
		{
			"total": 100,
			"data":[{},{}]
		}
	*/
	// todo: query也用结构体
	pNum := ctx.DefaultQuery("page_num", "0")
	pSize := ctx.DefaultQuery("page_size", "0")
	resp, err := global.GoodsSrvClient.CategoryBrandList(context.Background(), &proto.CategoryBrandFilterRequest{
		Pages:       cast.ToInt32(pNum),
		PagePerNums: cast.ToInt32(pSize),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	categoryBrandList := make([]*response.CategoryBrandInfoResp, 0, len(resp.Data))
	for _, data := range resp.Data {
		categoryBrandList = append(categoryBrandList, &response.CategoryBrandInfoResp{
			Id:       data.Id,
			Category: category.ConvertCategoryInfo2Resp(data.Category),
			Brand:    ConvertBrandInfo2Response(data.Brand),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"total":  resp.Total,
		"data":   categoryBrandList,
	})
}

func UpdateCategoryBrand(ctx *gin.Context) {
	form := &forms.CategoryBrandForm{}
	if err := ctx.ShouldBindJSON(form); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	_, err := global.GoodsSrvClient.UpdateCategoryBrand(context.Background(), &proto.CategoryBrandRequest{
		Id:         cast.ToInt32(id),
		CategoryId: form.CategoryId,
		BrandId:    form.BrandId,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}
	ctx.Status(http.StatusOK)
}

func DeleteCategoryBrand(ctx *gin.Context) {
	id := ctx.Param("id")
	_, err := global.GoodsSrvClient.DeleteCategoryBrand(context.Background(), &proto.CategoryBrandRequest{
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
