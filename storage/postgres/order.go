package postgres

import (
	"encoding/json"
	"fmt"
	"time"
	pb "product/genproto/ProductService"

	"github.com/google/uuid"
)

// User id ni tekshirish kerak
func (p *Product) PlaceOrder(req *pb.PlaceOrderRequest) (*pb.PlaceOrderResponse, error) {
	id := uuid.New().String()

	var totalAmount float64
	items := req.Items
	for _, item := range items {
		totalAmount += float64(item.Quantity) * item.Price
	}

	status := "Pending"

	shippingAddressJSON, err := json.Marshal(req.ShippingAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal shipping address: %v", err)
	}

	query := `
		INSERT INTO orders(
			id, user_id, total_amount, status, shipping_address, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
	`

	_, err = p.db.Exec(query, id, req.UserId, totalAmount, status, shippingAddressJSON, time.Now(), time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to insert order into database: %v", err)
	}

	for _, item := range items {
		orderId := uuid.New().String()
		orderQueryItems := `
			INSERT INTO order_items(
				id, order_id, product_id, quantity, price
			) VALUES (
				$1, $2, $3, $4, $5
			)
		`

		_, err := p.db.Exec(orderQueryItems, orderId, id, item.ProductId, item.Quantity, item.Price)
		if err != nil {
			return nil, err
		}
	}

	resp := &pb.PlaceOrderResponse{
		Id:              id,
		UserId:          req.UserId,
		Items:           req.Items,
		TotalAmount:     totalAmount,
		Status:          status,
		ShippingAddress: req.ShippingAddress,
	}

	return resp, nil
}

func (p *Product) CancelOrder(req *pb.CancelOrderRequest) (*pb.CancelOrderResponse, error) {
	query := `update 
				orders
			set
				status=$1, deleted_at=$2
			where
				id=$3`

	status := "canceled"
	newtime := time.Now()
	_, err := p.db.Exec(query, status, newtime, req.OrderId)
	if err != nil {
		return nil, err
	}
	return &pb.CancelOrderResponse{Message: "Your order has been cancelled"}, nil
}

func (p *Product) UpdateOrderStatus(req *pb.UpdateOrderStatusRequest) (*pb.UpdateOrderStatusResponse, error) {
	query := `update 
				orders
			set
				status=$1, deleted_at=$2
			where
				id=$3`

	status := "delivered"
	newtime := time.Now()
	_, err := p.db.Exec(query, status, newtime, req.OrderId)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateOrderStatusResponse{
		Id:     req.OrderId,
		Status: "Product with this id has been delivered",
	}, nil
}

func (p *Product) ListOrders(req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	offset := (req.Page - 1) * req.Limit

	query := `select 
				id, user_id, total_amount, status, shipping_address
			from
				orders
			order by
				created_at desc
			LIMIT $1 OFFSET $2`

	rows, err := p.db.Query(query, req.Limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*pb.Order
	for rows.Next() {
		var order pb.Order
		var shippingAddress string
		err := rows.Scan(
			&order.Id,
			&order.UserId,
			&order.TotalAmount,
			&order.Status,
			&shippingAddress,
		)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(shippingAddress), &order.ShippingAddress)
		if err != nil {
			return nil, err
		}

		orderItemQuery := `select 
							product_id, quantity, price
						from
							order_items
						where
							order_id = $1`

		orderItemRows, err := p.db.Query(orderItemQuery, order.Id)
		if err != nil {
			return nil, err
		}
		defer orderItemRows.Close()

		var items []*pb.OrderItem
		for orderItemRows.Next() {
			var item pb.OrderItem
			err := orderItemRows.Scan(&item.ProductId, &item.Quantity, &item.Price)
			if err != nil {
				return nil, err
			}
			items = append(items, &item)
		}

		order.Items = items
		orders = append(orders, &order)
	}

	countQuery := "select count(*) from orders"
	var total int32
	err = p.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, err
	}

	response := &pb.ListOrdersResponse{
		Orders: orders,
		Total:  total,
		Page:   req.Page,
		Limit:  req.Limit,
	}

	return response, nil
}

func (p *Product) GetOrder(req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	orderQuery := `
		SELECT 
			id, user_id, total_amount, status, shipping_address
		FROM 
			orders
		WHERE 
			id = $1
	`

	var order pb.GetOrderResponse
	var shippingAddress string

	err := p.db.QueryRow(orderQuery, req.OrderId).Scan(
		&order.Id,
		&order.UserId,
		&order.TotalAmount,
		&order.Status,
		&shippingAddress,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get order from database: %v", err)
	}

	err = json.Unmarshal([]byte(shippingAddress), &order.ShippingAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal shipping address: %v", err)
	}

	orderItemsQuery := `
		SELECT 
			product_id, quantity, price
		FROM 
			order_items
		WHERE 
			order_id = $1
	`

	orderItemsRows, err := p.db.Query(orderItemsQuery, req.OrderId)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items from database: %v", err)
	}
	defer orderItemsRows.Close()

	var items []*pb.OrderItem

	for orderItemsRows.Next() {
		var item pb.OrderItem
		err := orderItemsRows.Scan(&item.ProductId, &item.Quantity, &item.Price)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %v", err)
		}
		items = append(items, &item)
	}

	order.Items = items

	return &order, nil
}