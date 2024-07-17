package testing

import (
	"database/sql"
	"fmt"
	"log"
	"product/storage/postgres"
	"testing"

	pb "product/genproto/ProductService"
)

func Connection() *sql.DB {
	db, err := postgres.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func TestManageArtisanCategories(t *testing.T) {
	dbpool := Connection()
	defer dbpool.Close()

	product := postgres.NewProduct(dbpool)
	req := &pb.ManageArtisanCategoriesRequest{
		Userid: "27b8e30f-9823-487a-8dcd-845013abe2b1",
	}
	name := "test-category"

	resp, err := product.ManageArtisanCategories(name, req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
}

func TestManageProductCategories(t *testing.T) {
	dbpool := Connection()
	defer dbpool.Close()

	product := postgres.NewProduct(dbpool)
	req := &pb.ManageProductCategoriesRequest{
		ProductId: "87c94888-6f18-4a98-8dbc-d3b69f360726",
	}

	resp, err := product.ManageProductCategories(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("test-product-id", resp.ProductId)
}

func TestGetStatistics(t *testing.T) {
	dbpool := Connection()
	defer dbpool.Close()

	product := postgres.NewProduct(dbpool)
	req := &pb.GetStatisticsRequest{
		OrderId:   "d8886836-36f0-41dd-b4b4-fb781895a048",
		Statistik: 2,
	}

	resp, err := product.GetStatistics(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
}

func TestTrackUserActivity(t *testing.T) {
	dbpool := Connection()
	defer dbpool.Close()

	product := postgres.NewProduct(dbpool)
	req := &pb.TrackUserActivityRequest{
		UserId: "9dbd532d-b442-44ca-857a-4eb04f75fa9b",
	}

	resp, err := product.TrackUserActivity(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
}

func TestGetProductRecommendations(t *testing.T) {
	dbpool := Connection()
	defer dbpool.Close()

	product := postgres.NewProduct(dbpool)
	req := &pb.GetRecommendationsRequest{
		Productid: "87c94888-6f18-4a98-8dbc-d3b69f360726",
	}

	resp, err := product.GetProductRecommendations(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("test-product-id", resp.Id)
}

func TestGetArtisanRankings(t *testing.T) {
	dbpool := Connection()
	defer dbpool.Close()

	product := postgres.NewProduct(dbpool)

	resp, err := product.GetArtisanRankings(&pb.GetArtisanRankingsRequest{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
}

func TestPlaceOrder(t *testing.T) {
	dbpool := Connection()
	defer dbpool.Close()

	product := postgres.NewProduct(dbpool)
	shippingAddress := &pb.Address{
		Street:  "Test Street",
		City:    "Test City",
		Country: "Test Country",
		ZipCode: "Test ZipCode",
	}

	items := []*pb.OrderItem{
		{
			ProductId: "87c94888-6f18-4a98-8dbc-d3b69f360726",
			Quantity:  2,
			Price:     10.0,
		},
		{
			ProductId: "87c94888-6f18-4a98-8dbc-d3b69f360726",
			Quantity:  1,
			Price:     15.0,
		},
	}

	req := &pb.PlaceOrderRequest{
		UserId:          "9dbd532d-b442-44ca-857a-4eb04f75fa9b",
		Items:           items,
		ShippingAddress: shippingAddress,
	}

	resp, err := product.PlaceOrder(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
}

func TestCancelOrder(t *testing.T) {
	dbpool := Connection()
	defer dbpool.Close()

	product := postgres.NewProduct(dbpool)

	req := &pb.CancelOrderRequest{
		OrderId: "14952d98-05fd-4d88-8ca4-9b4cd2865aa4",
	}

	resp, err := product.CancelOrder(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
}

func TestUpdateOrderStatus(t *testing.T) {
	dbpool := Connection()
	defer dbpool.Close()

	product := postgres.NewProduct(dbpool)

	req := &pb.UpdateOrderStatusRequest{
		OrderId: "14952d98-05fd-4d88-8ca4-9b4cd2865aa4",
	}

	resp, err := product.UpdateOrderStatus(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
}

func TestListOrders(t *testing.T) {
	dbpool := Connection()
	defer dbpool.Close()

	product := postgres.NewProduct(dbpool)

	req := &pb.ListOrdersRequest{
		Page:  1,
		Limit: 10,
	}

	resp, err := product.ListOrders(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
}

func TestGetOrder(t *testing.T) {
	dbpool := Connection()
	defer dbpool.Close()

	product := postgres.NewProduct(dbpool)

	req := &pb.GetOrderRequest{
		OrderId: "14952d98-05fd-4d88-8ca4-9b4cd2865aa4",
	}

	resp, err := product.GetOrder(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
}

func TestPayOrder(t *testing.T) {
	dbpool := Connection()
	defer dbpool.Close()

	product := postgres.NewProduct(dbpool)

	paymentDetails := &pb.PaymentDetails{
		CardNumber: "1234567812345678",
		CardExpiry: "12/24",
	}

	req := &pb.PayOrderRequest{
		OrderId:        "14952d98-05fd-4d88-8ca4-9b4cd2865aa4",
		PaymentDetails: paymentDetails,
		PaymentMethod:  "Credit Card",
	}

	resp, err := product.PayOrder(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Your action has been accepted", resp.Status)
}

func TestCheckPaymentStatus(t *testing.T) {
	dbpool := Connection()
	defer dbpool.Close()

	product := postgres.NewProduct(dbpool)

	req := &pb.CheckPaymentStatusRequest{
		OrderId: "14952d98-05fd-4d88-8ca4-9b4cd2865aa4",
	}

	resp, err := product.CheckPaymentStatus(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
}

func TestAddProductCategori(t *testing.T) {
	db := Connection()
	defer db.Close()

	product := postgres.NewProduct(db)

	req := &pb.AddProductCategoriRequest{
		Name:        "Test Category",
		Description: "This is a test category",
	}

	resp, err := product.AddProductCategori(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
}

func TestUpdateProductCategori(t *testing.T) {
	db := Connection()
	defer db.Close()

	product := postgres.NewProduct(db)

	req := &pb.ProductCategori{
		Id:          "ad6931d6-1252-46fb-889a-15365c716a12",
		Name:        "Updated Category",
		Description: "This is an updated category",
	}

	resp, err := product.UpdateProductCategori(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
}

func TestDeleteProductCategori(t *testing.T) {
	db := Connection()
	defer db.Close()

	product := postgres.NewProduct(db)

	req := &pb.DeleteProductCategoriRequest{
		Id: "ad6931d6-1252-46fb-889a-15365c716a12",
	}

	resp, err := product.DeleteProductCategori(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
}

func TestAddProduct(t *testing.T) {
	db := Connection()
	defer db.Close()

	product := postgres.NewProduct(db)

	req := &pb.AddProductRequest{
		Name:        "Test Product",
		Description: "This is a test product",
		Price:       100.0,
		CategoryId:  "c3d3971f-4929-46e4-ae45-9df7e02dfd77",
		ArtisanId:   "dbb73cf0-5316-4122-8d75-a7da28e64047",
		Quantity:    10,
	}

	resp, err := product.AddProduct(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
}

func TestUpdateProduct(t *testing.T) {
	db := Connection()
	defer db.Close()

	product := postgres.NewProduct(db)

	req := &pb.UpdateProductRequest{
		ProductId:   "ff3b805d-556b-4844-b041-2ccd7ca19bbf",
		Name:        "Updated Product",
		Description: "Updated description",
		Price:       150.0,
		CategoryId:  "c3d3971f-4929-46e4-ae45-9df7e02dfd77",
		ArtisanId:   "dbb73cf0-5316-4122-8d75-a7da28e64047",
		Quantity:    20,
	}

	resp, err := product.UpdateProduct(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
}

func TestDeleteProduct(t *testing.T) {
	db := Connection()
	defer db.Close()

	product := postgres.NewProduct(db)

	req := &pb.DeleteProductRequest{
		ProductId: "ff3b805d-556b-4844-b041-2ccd7ca19bbf",
	}

	resp, err := product.DeleteProduct(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
}

func TestListProducts(t *testing.T) {
	db := Connection()
	defer db.Close()

	product := postgres.NewProduct(db)

	req := &pb.ListProductsRequest{
		Page:  1,
		Limit: 10,
	}

	resp, err := product.ListProducts(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
}

func TestGetProduct(t *testing.T) {
	db := Connection()
	defer db.Close()

	product := postgres.NewProduct(db)

	req := &pb.GetProductRequest{
		ProductId: "91c8a0dd-88b4-4cd8-a19f-349b27959b16",
	}

	resp, err := product.GetProduct(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
}

func TestSearchProducts(t *testing.T) {
	db := Connection()
	defer db.Close()

	product := postgres.NewProduct(db)

	req := &pb.SearchProductsRequest{
		ProductId: "91c8a0dd-88b4-4cd8-a19f-349b27959b16",
	}

	resp, err := product.SearchProducts(req)
	if err!=nil{
		t.Fatal(err)
	}
	fmt.Println(resp)
}

func TestUpdateShippingDetails(t *testing.T) {
	db := Connection()
	defer db.Close()

	product := postgres.NewProduct(db)

	shippingAddress := &pb.Address{
		Street:  "123 Test St",
		City:    "Test City",
		ZipCode: "12345",
		Country: "Test Country",
	}

	req := &pb.UpdateShippingDetailsRequest{
		OrderId:        "412c037e-c134-4c05-bc7a-f813a462e753",
		ShippingAddress: shippingAddress,
	}

	resp, err := product.UpdateShippingDetails(req)
	if err!=nil{
		t.Fatal(err)
	}
	fmt.Println(resp)
}

func TestAddProductRating(t *testing.T) {
	db := Connection()
	defer db.Close()

	product := postgres.NewProduct(db)

	req := &pb.AddProductRatingRequest{
		ProductId: "42937cab-0787-47b3-ba31-90de9eb5aa19",
		UserId:    "3f647e51-3b39-4113-98bb-071feb11ef6f",
		Rating:    4.5,
		Comment:   "Great product!",
	}

	resp, err := product.AddProductRating(req)
	if err!=nil{
		t.Fatal(err)
	}
	fmt.Println(resp)
}