// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Blog struct {
	ID        int64     `json:"id"`
	Author    int64     `json:"author"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Order struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	Amount          float64   `json:"amount"`
	Status          string    `json:"status"`
	ShippingAddress string    `json:"shipping_address"`
	CreatedAt       time.Time `json:"created_at"`
}

type OrderItem struct {
	ID        int64   `json:"id"`
	OrderID   int64   `json:"order_id"`
	ProductID int64   `json:"product_id"`
	Color     *string `json:"color"`
	Size      *string `json:"size"`
	Quantity  int32   `json:"quantity"`
}

type Product struct {
	ID          int64   `json:"id"`
	ProductName string  `json:"product_name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int64   `json:"quantity"`
	// admins may have discount. Float of percentage ie 14.5
	Discount *float64 `json:"discount"`
	// calculate when reviews is created. 1-5
	Rating       *int32   `json:"rating"`
	SizeOptions  []string `json:"size_options"`
	ColorOptions []string `json:"color_options"`
	Category     string   `json:"category"`
	Brand        *string  `json:"brand"`
	// list of file paths to the product images
	ImageUrl  []string  `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Review struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	ProductID int64  `json:"product_id"`
	Rating    int32  `json:"rating"`
	Review    string `json:"review"`
}

type Session struct {
	ID           pgtype.UUID `json:"id"`
	UserID       int64       `json:"user_id"`
	RefreshToken string      `json:"refresh_token"`
	IsBlocked    bool        `json:"is_blocked"`
	UserAgent    string      `json:"user_agent"`
	UserIp       string      `json:"user_ip"`
	ExpiresAt    time.Time   `json:"expires_at"`
	CreatedAt    time.Time   `json:"created_at"`
}

type User struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Subscription bool   `json:"subscription"`
	// list of product id in the cart
	UserCart []int64 `json:"user_cart"`
	// user or admin
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
