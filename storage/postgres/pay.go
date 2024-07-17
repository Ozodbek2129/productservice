package postgres

import (
	"database/sql"
	"fmt"
	pb "product/genproto/ProductService"
	"time"

	"github.com/google/uuid"
)

func (p *Product) PayOrder(req *pb.PayOrderRequest) (*pb.PayOrderResponse, error) {
	id := uuid.NewString()

	paymentDetails := req.PaymentDetails
	if len(paymentDetails.CardNumber) != 16 {
		return nil, fmt.Errorf("sizning karta raqamingiz 16 ta raqam emas?")
	}

	if len(paymentDetails.CardExpiry) != 5 {
		return nil, fmt.Errorf("sizning karta muddatingiz xato kirtildi?")
	}

	query := `insert into payments(
				id, order_id, amount, status, card_number, card_expiry, transaction_id, payment_method, created_at
			)values(
				$1,$2,(select total_amount from orders where id=$3),$4,$5,$6,$7,$8,$9)`

	newtime := time.Now()
	_, err := p.db.Exec(query, id, req.OrderId, req.OrderId, "Completed", paymentDetails.CardNumber, paymentDetails.CardExpiry, "nimadur", req.PaymentMethod, newtime)
	if err != nil {
		return nil, err
	}

	return &pb.PayOrderResponse{
		Id:     id,
		Status: "Your action has been accepted",
	}, nil
}

func (p *Product) CheckPaymentStatus(req *pb.CheckPaymentStatusRequest) (*pb.CheckPaymentStatusResponse, error) {
	query := `select 
				status
			from
				payments
			where
				order_id=$1`

	var status string
	err := p.db.QueryRow(query, req.OrderId).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no payment found for order_id %s", req.OrderId)
		}
		return nil, err
	}
	return &pb.CheckPaymentStatusResponse{
		Status: status,
	}, nil
}
