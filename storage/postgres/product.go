package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	pb "product/genproto/ProductService"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	db *sql.DB
}

func NewProduct(db *sql.DB) *Product {
	return &Product{db: db}
}

func (p *Product) AddProductCategori(req *pb.AddProductCategoriRequest) (*pb.ProductCategori, error) {
	query := `insert into product_categories(
				id, name, description, created_at
			)values(
				$1,$2,$3,$4)`

	id := uuid.NewString()
	newtime := time.Now()

	_, err := p.db.Exec(query, id, req.Name, req.Description, newtime)
	if err != nil {
		return nil, err
	}

	return &pb.ProductCategori{Id: id, Name: req.Name, Description: req.Description}, nil
}

func (p *Product) UpdateProductCategori(req *pb.ProductCategori) (*pb.ProductCategori, error) {
	query := `update 
				product_categories
			set
				name=$1, description=$2, created_at=$3
			where 
				id=$4 `

	newtime := time.Now()
	_, err := p.db.Exec(query, req.Name, req.Description, newtime, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.ProductCategori{Id: req.Id, Name: req.Name, Description: req.Description}, nil
}

func (p *Product) DeleteProductCategori(req *pb.DeleteProductCategoriRequest) (*pb.Status, error) {
	query := `delete from 
				product_categories
			  where
				id=$1`

	_, err := p.db.Exec(query, req.Id)
	fmt.Println(req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.Status{Message: "The category you specified has been deleted"}, nil
}

// User id ni tekshirish kerak
func (p *Product) AddProduct(req *pb.AddProductRequest) (*pb.AddProductResponse, error) {
	query := `insert into products(
				id, name, description, price, category_id, artisan_id, quantity, created_at, updated_at
			)values(
				$1,$2,$3,$4,$5,$6,$7,$8,$9)`

	id := uuid.NewString()
	newtime := time.Now()
	_, err := p.db.Exec(query, id, req.Name, req.Description, req.Price, req.CategoryId, req.ArtisanId, req.Quantity, newtime, newtime)
	if err != nil {
		return nil, err
	}

	return &pb.AddProductResponse{
		Id:          id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		CategoryId:  req.CategoryId,
		ArtisanId:   req.ArtisanId,
		Quantity:    req.Quantity,
	}, nil
}

// user id ni tekshirish kerak
func (p *Product) UpdateProduct(req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	fmt.Println(req.ProductId)
	fmt.Println("fifbbfbfbuifdub")
	query := `update
				products
			set
				name=$1,
				description=$2,
				price=$3,
				quantity=$4,
				artisan_id=$5,
				updated_at=$6
			where 
				id=$7 and
				deleted_at is null`

	newtime := time.Now()
	_, err := p.db.Exec(query, req.Name, req.Description, req.Price, req.Quantity, req.ArtisanId, newtime, req.ProductId)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateProductResponse{
		Id:          req.ProductId,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		CategoryId:  req.CategoryId,
		ArtisanId:   req.ArtisanId,
		Quantity:    req.Quantity,
	}, nil
}

func (p *Product) DeleteProduct(req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	query := `update 
				products
			set
				deleted_at=$1
			where
				id=$2`

	newtime := time.Now()
	_, err := p.db.Exec(query, newtime, req.ProductId)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteProductResponse{
		Message: "The product you specified has been removed",
	}, nil
}

func (p *Product) ListProducts(req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	offset := (req.Page - 1) * req.Limit
	query := `
		select 
			id, name, description, price, category_id, quantity
		from 
			products
		order by 
			created_at DESC
		LIMIT $1 OFFSET $2;
	`

	rows, err := p.db.Query(query, req.Limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*pb.Product
	for rows.Next() {
		var product pb.Product
		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.CategoryId,
			&product.Quantity,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	countQuery := "SELECT COUNT(*) FROM products;"
	var total int32
	err = p.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, err
	}

	response := &pb.ListProductsResponse{
		Products: products,
		Total:    total,
		Page:     req.Page,
		Limit:    req.Limit,
	}

	return response, nil
}

func (p *Product) GetProduct(req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	query := `select 
				id, name, description, price, category_id, artisan_id, quantity
			from 
				products
			where
				id=$1`

	product := pb.GetProductResponse{}
	err := p.db.QueryRow(query, req.ProductId).Scan(
		&product.Id,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.CategoryId,
		&product.ArtisanId,
		&product.Quantity,
	)

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (p *Product) SearchProducts(req *pb.SearchProductsRequest) (*pb.SearchProductsResponse, error) {
	query := `
		SELECT
			id, name, description, price, category_id, artisan_id, quantity
		FROM
			products
		WHERE
			id = $1
	`
	rows, err := p.db.Query(query, req.ProductId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*pb.Product
	for rows.Next() {
		var product pb.Product
		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.CategoryId,
			&product.ArtisanId,
			&product.Quantity,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	countQuery := "SELECT COUNT(*) FROM products WHERE id = $1"
	var total int32
	err = p.db.QueryRow(countQuery, req.ProductId).Scan(&total)
	if err != nil {
		return nil, err
	}

	response := &pb.SearchProductsResponse{
		Products: products,
		Total:    total,
	}

	return response, nil
}

func (p *Product) UpdateShippingDetails(req *pb.UpdateShippingDetailsRequest) (*pb.UpdateShippingDetailsResponse, error) {
	shippingAddress, err := json.Marshal(req.ShippingAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal shipping address: %v", err)
	}

	query := `update 
				orders
			set
				shipping_address=$1, updated_at=$2
			where
				id=$3
			returning
				id,shipping_address`

	var orderId, updateShippingAddressJson string
	newtime := time.Now()
	err = p.db.QueryRow(query, shippingAddress, newtime, req.OrderId).Scan(&orderId, &updateShippingAddressJson)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no order found with id %s", req.OrderId)
		}
		return nil, err
	}

	var updateShippingAddress pb.Address
	err = json.Unmarshal([]byte(updateShippingAddressJson), &updateShippingAddress)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateShippingDetailsResponse{
		Id:              orderId,
		ShippingAddress: &updateShippingAddress,
	}, nil
}
