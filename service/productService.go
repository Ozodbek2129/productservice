package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"product/config"
	userClient "product/genproto/AuthService"
	pb "product/genproto/ProductService"
	"product/pkg/logger"
	"product/storage/postgres"

	"google.golang.org/grpc"
)

type ArtisanConnectService struct {
	pb.UnimplementedProductServiceServer
	Product    *postgres.Product
	Logger     *slog.Logger
	UserClient userClient.AuthUserServiceClient
}

func NewArtisanConnectService(db *sql.DB) *ArtisanConnectService {
	conn, err := grpc.Dial(config.Load().USER_SERVICE, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to user service: %w", err)
	}

	return &ArtisanConnectService{
		Product:    postgres.NewProduct(db),
		Logger:     logger.NewLogger(),
		UserClient: userClient.NewAuthUserServiceClient(conn),
	}
}

func (s *ArtisanConnectService) AddProductCategori(ctx context.Context, req *pb.AddProductCategoriRequest) (*pb.ProductCategori, error) {
	resp, err := s.Product.AddProductCategori(req)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("error adding categories: %v", err))
		return nil, err
	}
	return resp, nil
}

func (s *ArtisanConnectService) UpdateProductCategori(ctx context.Context, req *pb.ProductCategori) (*pb.ProductCategori, error) {
	resp, err := s.Product.UpdateProductCategori(req)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("error updating categories: %v", err))
		return nil, err
	}
	return resp, nil
}

func (s *ArtisanConnectService) DeleteProductCategori(ctx context.Context, req *pb.DeleteProductCategoriRequest) (*pb.Status, error) {
	resp, err := s.Product.DeleteProductCategori(req)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("error deleting categories: %v", err))
		return nil, err
	}
	return resp, nil
}

func (s *ArtisanConnectService) AddProduct(ctx context.Context, req *pb.AddProductRequest) (*pb.AddProductResponse, error) {
	idCheckReq := &userClient.Id{
		Id: req.ArtisanId,
	}
	idCheckResp, err := s.UserClient.IdCheck(ctx, idCheckReq)
	if err != nil || !idCheckResp.B {
		s.Logger.Error(fmt.Sprintf("No such user id: %v", err))
		return nil, fmt.Errorf("no such user id: %v", err)
	}

	resp, err := s.Product.AddProduct(req)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("Error adding products: %v", err))
		return nil, err
	}
	return resp, nil
}

func (s *ArtisanConnectService) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	idCheckReq := &userClient.Id{
		Id: req.ArtisanId,
	}
	idCheckResp, err := s.UserClient.IdCheck(ctx, idCheckReq)
	if err != nil || !idCheckResp.B {
		s.Logger.Error(fmt.Sprintf("No such user id: %v", err))
		return nil, fmt.Errorf("no such user id: %v", err)
	}

	resp, err := s.Product.UpdateProduct(req)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("error updating products: %v", err))
		return nil, err
	}
	return resp, nil
}


func (s *ArtisanConnectService) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest)(*pb.DeleteProductResponse,error){
	resp,err:=s.Product.DeleteProduct(req)
	if err!=nil{
		s.Logger.Error(fmt.Sprintf("The reference on the id you specified has been deleted: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *ArtisanConnectService) ListProducts(ctx context.Context, req *pb.ListProductsRequest)(*pb.ListProductsResponse,error){
	resp,err:=s.Product.ListProducts(req)
	if err!=nil{
		s.Logger.Error(fmt.Sprintf("error retrieving product list: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *ArtisanConnectService) GetProduct(ctx context.Context,req *pb.GetProductRequest)(*pb.GetProductResponse,error){
	resp,err:=s.Product.GetProduct(req)
	if err!=nil{
		s.Logger.Error(fmt.Sprintf("error in receiving products: %v",err))
		return nil,err
	}

	return resp,nil
}

func (s *ArtisanConnectService) SearchProducts(ctx context.Context,req *pb.SearchProductsRequest)(*pb.SearchProductsResponse,error){
	resp,err:=s.Product.SearchProducts(req)
	if err!=nil{
		s.Logger.Error(fmt.Sprintf("product search error: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *ArtisanConnectService) AddProductRating(ctx context.Context,req *pb.AddProductRatingRequest)(*pb.AddProductRatingResponse,error){
	idCheckReq := &userClient.Id{
		Id: req.UserId,
	}
	idCheckResp, err := s.UserClient.IdCheck(ctx, idCheckReq)
	if err != nil || !idCheckResp.B {
		s.Logger.Error(fmt.Sprintf("No such user id: %v", err))
		return nil, fmt.Errorf("no such user id: %v", err)
	}

	resp,err:=s.Product.AddProductRating(req)
	if err!=nil{
		s.Logger.Error(fmt.Sprintf("error adding product rating: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *ArtisanConnectService) ListProductRatings(ctx context.Context,req *pb.ListProductRatingsRequest)(*pb.ListProductRatingsResponse,error){
	resp,err:=s.Product.ListProductRatings(req)
	if err!=nil{
		s.Logger.Error(fmt.Sprintf("error in getting the list of product ratings: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *ArtisanConnectService) PlaceOrder(ctx context.Context,req *pb.PlaceOrderRequest)(*pb.PlaceOrderResponse,error){
	idCheckReq := &userClient.Id{
		Id: req.UserId,
	}
	idCheckResp, err := s.UserClient.IdCheck(ctx, idCheckReq)
	if err != nil || !idCheckResp.B {
		s.Logger.Error(fmt.Sprintf("No such user id: %v", err))
		return nil, fmt.Errorf("no such user id: %v", err)
	}

	resp,err:=s.Product.PlaceOrder(req)
	if err!=nil{
		s.Logger.Error(fmt.Sprintf("ordering error: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *ArtisanConnectService) CancelOrder(ctx context.Context,req *pb.CancelOrderRequest)(*pb.CancelOrderResponse,error){
	resp,err:=s.Product.CancelOrder(req)
	if err!=nil{
		s.Logger.Error(fmt.Sprintf("order cancellation: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *ArtisanConnectService) UpdateOrderStatus(ctx context.Context,req *pb.UpdateOrderStatusRequest)(*pb.UpdateOrderStatusResponse,error){
	resp,err:=s.Product.UpdateOrderStatus(req)
	if err!=nil{
		s.Logger.Error(fmt.Sprintf("error updating order status: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *ArtisanConnectService) ListOrders(ctx context.Context,req *pb.ListOrdersRequest)(*pb.ListOrdersResponse,error){
	resp,err:=s.Product.ListOrders(req)
	if err!=nil{
		s.Logger.Error(fmt.Sprintf("view the list of orders: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *ArtisanConnectService) GetOrder(ctx context.Context,req *pb.GetOrderRequest)(*pb.GetOrderResponse,error){
	resp,err:=s.Product.GetOrder(req)
	if err!=nil{
		s.Logger.Error(fmt.Sprintf("error in receiving the order: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *ArtisanConnectService) PayOrder(ctx context.Context,req *pb.PayOrderRequest)(*pb.PayOrderResponse,error){
	resp,err:=s.Product.PayOrder(req)
	if err!=nil{
		s.Logger.Error(fmt.Sprintf("error when paying for the order: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *ArtisanConnectService) CheckPaymentStatus(ctx context.Context,req *pb.CheckPaymentStatusRequest)(*pb.CheckPaymentStatusResponse,error){
	resp,err:=s.Product.CheckPaymentStatus(req)
	if err!=nil{
		s.Logger.Error(fmt.Sprintf("Error checking payment status: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *ArtisanConnectService) UpdateShippingDetails(ctx context.Context,req *pb.UpdateShippingDetailsRequest)(*pb.UpdateShippingDetailsResponse,error){
	resp,err:=s.Product.UpdateShippingDetails(req)
	if err!=nil{
		s.Logger.Error(fmt.Sprintf("update delivery details: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *ArtisanConnectService) ManageArtisanCategories(ctx context.Context,req *pb.ManageArtisanCategoriesRequest)(*pb.ManageArtisanCategoriesResponse,error){
	idCheckReq := &userClient.Id{
		Id: req.Userid,
	}
	idCheckResp, err := s.UserClient.IdCheck(ctx, idCheckReq)
	if err != nil || !idCheckResp.B {
		s.Logger.Error(fmt.Sprintf("No such user id: %v", err))
		return nil, fmt.Errorf("no such user id: %v", err)
	}

	name,err:=s.UserClient.SearchName(ctx,idCheckReq)
	if err!=nil{
		s.Logger.Error(fmt.Sprintf("Manage Artisan Categories error: %v",err))
		return nil,err
	}

	resp,err:=s.Product.ManageArtisanCategories(name.Name,req)
	if err!=nil{
		s.Logger.Error(fmt.Sprintf("error in retrieving data from the database:  %v",err))
	}

	return resp,nil
}

func (s *ArtisanConnectService) ManageProductCategories(ctx context.Context,req *pb.ManageProductCategoriesRequest)(*pb.ManageProductCategoriesResponse,error){
	resp,err:=s.Product.ManageProductCategories(req)
	if err!=nil{
		s.Logger.Error(fmt.Sprintf("Error managing product categories: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *ArtisanConnectService) GetStatistics(ctx context.Context,req *pb.GetStatisticsRequest)(*pb.GetStatisticsResponse,error){
	resp,err:=s.Product.GetStatistics(req)
	if err!=nil{
		s.Logger.Error(fmt.Sprintf("Error getting statistics: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *ArtisanConnectService) TrackUserActivity(ctx context.Context,req *pb.TrackUserActivityRequest)(*pb.TrackUserActivityResponse,error){
	idCheckReq := &userClient.Id{
		Id: req.UserId,
	}
	idCheckResp, err := s.UserClient.IdCheck(ctx, idCheckReq)
	if err != nil || !idCheckResp.B {
		s.Logger.Error(fmt.Sprintf("No such user id: %v", err))
		return nil, fmt.Errorf("no such user id: %v", err)
	}

	resp,err:=s.Product.TrackUserActivity(req)
	if err!=nil{
		s.Logger.Error(fmt.Sprintf("Error tracking user activity: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *ArtisanConnectService) GetProductRecommendations(ctx context.Context,req *pb.GetRecommendationsRequest)(*pb.GetRecommendationsResponse,error){
	resp,err:=s.Product.GetProductRecommendations(req)
	if err!=nil{
		s.Logger.Error(fmt.Sprintf("Error getting recommendations: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *ArtisanConnectService) GetArtisanRankings(ctx context.Context,req *pb.GetArtisanRankingsRequest)(*pb.GetArtisanRankingsResponse,error){
	resp,err:=s.Product.GetArtisanRankings(req)
	if err!=nil{
		s.Logger.Error(fmt.Sprintf("Error getting crafter rating: %v",err))
		return nil,err
	}
	return resp,nil
}