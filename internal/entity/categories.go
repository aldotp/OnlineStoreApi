package entity

import "time"

type Category struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Products  []Product `json:"products,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
