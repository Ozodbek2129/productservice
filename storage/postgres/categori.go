package postgres

import (
	pb "product/genproto/ProductService"
)

func (p *Product) ManageArtisanCategories(name string, req *pb.ManageArtisanCategoriesRequest) (*pb.ManageArtisanCategoriesResponse, error) {
	query := `SELECT 
				p.id, p.name, o.quantity, p.price, COALESCE(AVG(r.rating), 0) as rating, COALESCE(r.comment, '')
			FROM
				products as p
			LEFT JOIN
				order_items as o ON p.id = o.product_id
			LEFT JOIN
				ratings as r ON p.id = r.product_id
			WHERE
				p.name = $1
			GROUP BY
				p.id, p.name, o.quantity, p.price, r.comment`

	rows, err := p.db.Query(query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resp pb.ManageArtisanCategoriesResponse
	for rows.Next() {
		var productID string
		var productName string
		var quantity int32
		var price float32
		var rating float32
		var comment string
		err := rows.Scan(&productID, &productName, &quantity, &price, &rating, &comment)

		if err != nil {
			return nil, err
		}

		resp.Categories = append(resp.Categories, &pb.Category{
			Productid: productID,
			Name:      productName,
			Quantity:  quantity,
			Price:     price,
			Rating:    rating,
			Comment:   comment,
		})
	}
	return &resp, nil
}

func (p *Product) ManageProductCategories(req *pb.ManageProductCategoriesRequest) (*pb.ManageProductCategoriesResponse, error) {
	query := `select 
				c.name, c.description, p.id, p.name, p.description, p.price, p.quantity, r.rating
			from
				products as p
			join
				product_categories as c
			on 
				c.id=p.category_id
			left join
				ratings as r
			on
				p.id=r.product_id
			where
				p.id=$1`

	var productId string
	var categoriName string
	var categoriDescription string
	var productName string
	var productDescription string
	var quantity int32
	var price float32
	var rating float32
	err := p.db.QueryRow(query, req.ProductId).Scan(&categoriName, &categoriDescription, &productId, &productName, &productDescription, &price, &quantity, &rating)
	if err != nil {
		return nil, err
	}

	return &pb.ManageProductCategoriesResponse{
		ProductId:           productId,
		CategoriName:        categoriName,
		CategoriDescription: categoriDescription,
		ProductName:         productName,
		ProductDescription:  productDescription,
		Quantity:            quantity,
		Price:               price,
		Rating:              rating,
	}, nil
}

func (p *Product) GetStatistics(req *pb.GetStatisticsRequest)(*pb.GetStatisticsResponse,error){
	query:=`select 
				total_amount
			from
				orders
			where
				id=$1`

	var totalamount float32
	err:=p.db.QueryRow(query,req.OrderId).Scan(&totalamount)
	if err!=nil{
		return nil,err
	}
	
	statistika:=totalamount*float32(req.Statistik)

	return &pb.GetStatisticsResponse{
		Message: "Total amount",
		Value: float64(statistika),
	},nil
}

func (p *Product) TrackUserActivity(req *pb.TrackUserActivityRequest)(*pb.TrackUserActivityResponse,error){
	query:=`select
				total_amount
			from 
				orders
			where
				user_id=$1`

	var totalamount float32
	err:=p.db.QueryRow(query,req.UserId).Scan(&totalamount)
	if err!=nil{
		return nil,err
	}
	
	return &pb.TrackUserActivityResponse{
		UserId: req.UserId,
		TotalAmount: totalamount,
		Activity: "Activity",
	},nil
}

func (p *Product) GetProductRecommendations(req *pb.GetRecommendationsRequest)(*pb.GetRecommendationsResponse,error){
	query:=`select
				p.id, p.name, p.price, c.name
			from 
				products as p
			join
				product_categories as c
			on
				p.category_id=c.id
			where
				p.id=$1`

	resp:=pb.GetRecommendationsResponse{}
	err:=p.db.QueryRow(query,req.Productid).Scan(&resp.Id,&resp.Productname,&resp.Price,&resp.Categoriname)
	if err!=nil{
		return nil,err
	}
	return &resp,nil
}

func (p *Product) GetArtisanRankings(req *pb.GetArtisanRankingsRequest)(*pb.GetArtisanRankingsResponse,error){
	query:=` SELECT
				p.artisan_id,
            	p.name AS name,
            	COALESCE(AVG(r.rating), 0) AS rating,
            	COUNT(p.id) AS total_products
       	 	FROM
            	products as p
        	JOIN
            	ratings as r ON r.product_id = p.id
        	GROUP BY
            	p.id, p.name
        	ORDER BY
            	rating DESC;
    		`

	rows,err:=p.db.Query(query)
	if err!=nil{
		return nil,err
	}
	defer rows.Close()

	var ranking []*pb.ArtisanRanking
	for rows.Next(){
		var rank pb.ArtisanRanking
		err:=rows.Scan(&rank.ArtisanId,&rank.Name,&rank.Rating,&rank.TotalProducts)
		if err!=nil{
			return nil,err
		}
		ranking=append(ranking, &rank)
	}

	return &pb.GetArtisanRankingsResponse{
		Rankings: ranking,
	},nil
}