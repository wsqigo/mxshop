package handler

import (
	"mxshop_srvs/goods_srv/proto"
)

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

//func (g GoodsServer) GoodsList(ctx context.Context, request *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (g GoodsServer) BatchGetGoods(ctx context.Context, info *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (g GoodsServer) CreateGoods(ctx context.Context, info *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (g GoodsServer) DeleteGoods(ctx context.Context, info *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (g GoodsServer) UpdateGoods(ctx context.Context, info *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (g GoodsServer) GetGoodsDetail(ctx context.Context, request *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (g GoodsServer) GetAllCategoryList(ctx context.Context, empty *emptypb.Empty) (*proto.CategoryListResponse, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (g GoodsServer) GetSubCategory(ctx context.Context, request *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (g GoodsServer) CreateCategory(ctx context.Context, request *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (g GoodsServer) DeleteCategory(ctx context.Context, request *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (g GoodsServer) UpdateCategory(ctx context.Context, request *proto.CategoryInfoRequest) (*emptypb.Empty, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (g GoodsServer) BrandList(ctx context.Context, request *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (g GoodsServer) CreateBrand(ctx context.Context, request *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (g GoodsServer) DeleteBrand(ctx context.Context, request *proto.BrandRequest) (*emptypb.Empty, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (g GoodsServer) UpdateBrand(ctx context.Context, request *proto.BrandRequest) (*emptypb.Empty, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (g GoodsServer) BannerList(ctx context.Context, empty *emptypb.Empty) (*proto.BannerListResponse, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (g GoodsServer) CreateBanner(ctx context.Context, request *proto.BannerRequest) (*proto.BannerResponse, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (g GoodsServer) DeleteBanner(ctx context.Context, request *proto.BannerRequest) (*emptypb.Empty, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (g GoodsServer) UpdateBanner(ctx context.Context, request *proto.BannerRequest) (*emptypb.Empty, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (g GoodsServer) CategoryBrandList(ctx context.Context, request *proto.CategoryBrandFilterRequest) (*proto.CategoryBrandListResponse, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (g GoodsServer) GetCategoryBrandList(ctx context.Context, request *proto.CategoryInfoRequest) (*proto.BrandListResponse, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (g GoodsServer) CreateCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*proto.CategoryBrandResponse, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (g GoodsServer) DeleteCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (g GoodsServer) UpdateCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
//	//TODO implement me
//	panic("implement me")
//}
