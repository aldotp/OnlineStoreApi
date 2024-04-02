package services

import "context"

type PaymentService interface {
	ProcessPayment(ctx context.Context, amount float64) string
}

type payment struct {
}

func NewPayment() PaymentService {
	return &payment{}
}

// ProcessPayment simulates the payment process with a payment gateway.
func (p *payment) ProcessPayment(ctx context.Context, amount float64) string {
	// Simulate payment processing with a payment gateway
	// Replace this with your actual payment gateway integration logic
	// For demonstration purposes, return "success" if the payment amount is greater than zero
	if amount > 0 {
		return "success"
	}
	return "failure"
}
