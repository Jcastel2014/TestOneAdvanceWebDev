package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/jcastel2014/test1/internal/validator"
)

type Reviews struct {
	ID            int64     `json:"id"`
	Product_Id    string    `json:"pid"`
	Rating        float64   `json:"rating"`
	Helpful_Count int       `json:"helpful_count"`
	Comment       string    `json:"comment"`
	Created_at    time.Time `json:"created_at"`
	Updated_at    time.Time `json:"updated_at"`
}

func (p ProductModel) InsertReview(review *Reviews, id int64) error {
	err := p.DoesProductExists(id)

	if err != nil {
		return err
	}

	query := `
	INSERT INTO reviews(product_id, rating, comment)
	VALUES ($1, $2, $3)
	RETURNING created_at, id
	
	`

	args := []any{id, review.Rating, review.Comment}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = p.DB.QueryRowContext(ctx, query, args...).Scan(&review.Created_at, &review.ID)

	if err != nil {
		return err
	}

	return p.UpdateAverage(id)

}

func (p ProductModel) UpdateAverage(pid int64) error {

	query := `
	UPDATE products
	Set average_rating = (select AVG(rating)::NUMERIC(10,2) from reviews where product_id = $1)
	WHERE id = $1
	RETURNING id
	`

	args := []any{pid}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return p.DB.QueryRowContext(ctx, query, args...).Scan(&pid)
}

func (p ProductModel) GetReview(id int64, rid int64) (*Reviews, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	} else if rid < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT id, product_id, rating, helpful_count, comment, created_at, updated_at 
	FROM reviews
        WHERE id = $1 AND product_id = $2
		
	`

	args := []any{rid, id}

	var review Reviews

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := p.DB.QueryRowContext(ctx, query, args...).Scan(&review.ID, &review.Product_Id, &review.Rating, &review.Helpful_Count, &review.Comment, &review.Created_at, &review.Updated_at)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &review, nil

}

func (p ProductModel) UpdateReview(review *Reviews, id int64) error {
	query := `
	UPDATE reviews
	SET rating =$1, comment=$2, updated_at=$3
	WHERE id = $4
	RETURNING id
	`

	args := []any{review.Rating, review.Comment, time.Now(), id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return p.DB.QueryRowContext(ctx, query, args...).Scan(&review.ID)

}
func (p ProductModel) DoesProductExists(id int64) error {
	query := `
		SELECT COUNT(*)
		FROM products
		WHERE id = $1
	`
	args := []any{id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var count int

	err := p.DB.QueryRowContext(ctx, query, args...).Scan(&count)

	if err != nil {
		return err
	}

	if count < 1 {
		return errors.New("no rows found")
	}

	return nil
}

func ValidateReview(v *validator.Validator, review *Reviews, handlerId int) {

	switch handlerId {
	case 1:
		v.Check(review.Comment != "", "comment", "must be provided")
		v.Check(review.Rating > 0, "rating", "must be greater than 0")
		v.Check(review.Rating <= 5, "rating", "must be less than 5")
		v.Check(len(review.Comment) <= 100, "comment", "must not be more than 100 byte long")

	default:
		log.Println("Unable to locate handler ID: %s", handlerId)
		v.AddError("default", "Handler ID not provided")
	}
}
