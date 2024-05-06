package db

import "context"

func (q *Queries) StoreCalculateProductRating(ctx context.Context, productID int64) (float64, error) {
	row := q.db.QueryRow(ctx, calculateProductRating, productID)
	var average_rating float64
	err := row.Scan(&average_rating)
	return average_rating, err
}
