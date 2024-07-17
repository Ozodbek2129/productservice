package postgres

import (
	"fmt"
	pb "product/genproto/ProductService"
	"time"

	"github.com/google/uuid"
)

// User ni tekshirib ketish kerak
func (p *Product) AddProductRating(req *pb.AddProductRatingRequest) (*pb.AddProductRatingResponse, error) {
	query := `insert into ratings(
				id, product_id, user_id, rating, comment, created_at
			)values(
				$1,$2,$3,$4,$5,$6) `

	id := uuid.NewString()
	newtime := time.Now()
	rating := fmt.Sprintf("%.1f", float64(req.Rating))
	
	_, err := p.db.Exec(query, id, req.ProductId, req.UserId, rating, req.Comment, newtime)
	if err != nil {
		return nil, err
	}

	return &pb.AddProductRatingResponse{
		Id:        id,
		ProductId: req.ProductId,
		UserId:    req.UserId,
		Rating:    req.Rating,
		Comment:   req.Comment,
	}, nil
}

func (p *Product) ListProductRatings(req *pb.ListProductRatingsRequest) (*pb.ListProductRatingsResponse, error) {
	query := `
		SELECT 
			user_id, rating, comment 
		FROM 
			ratings 
		WHERE 
			product_id = $1
	`
	rows, err := p.db.Query(query, req.ProductId)
	if err != nil {
		return nil, fmt.Errorf("failed to query product ratings: %v", err)
	}
	defer rows.Close()

	var ratings []*pb.Rating
	var totalRatings int32
	var sumRatings float64

	for rows.Next() {
		var r pb.Rating
		if err := rows.Scan(&r.UserId, &r.Rating, &r.Comment); err != nil {
			return nil, fmt.Errorf("failed to scan rating: %v", err)
		}
		ratings = append(ratings, &r)
		totalRatings++
		sumRatings += r.Rating
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during rows iteration: %v", err)
	}

	averageRating := 0.0
	if totalRatings > 0 {
		averageRating = sumRatings / float64(totalRatings)
	}

	return &pb.ListProductRatingsResponse{
		Ratings:       ratings,
		AverageRating: averageRating,
		TotalRatings:  totalRatings,
	}, nil
}
