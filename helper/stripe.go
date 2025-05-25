package helper

import (
	"context"
	"os"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/paymentintent"
	"github.com/stripe/stripe-go/v82/subscription"
)

type StripeClient struct {
	client *stripe.Client
}

func NewStripeClient() *StripeClient {
	key := os.Getenv("STRIPE_SECRET_KEY")
	stripe.Key = key
	return &StripeClient{client: stripe.NewClient(key)}
}

// CreateSubscriptionSession inicia una sesión de Checkout para un plan periódico
func (s *StripeClient) CreateSubscriptionSession(
	ctx context.Context,
	customerID, priceID,
	successURL,
	cancelURL string,
	metadata map[string]string,
	persist bool,
) (*stripe.CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{
		Customer:           stripe.String(customerID),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Mode:               stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{{
			Price:    stripe.String(priceID),
			Quantity: stripe.Int64(1),
		}},
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
	}

	if persist {
		params.SubscriptionData = &stripe.CheckoutSessionSubscriptionDataParams{
			Metadata: metadata,
		}
	} else {
		params.Metadata = metadata
	}

	return session.New(params)
}

// helper/stripe.go
func (s *StripeClient) CreateOneTimeSession(
	ctx context.Context,
	customerID, priceID,
	successURL, cancelURL string,
	metadata map[string]string,
) (*stripe.CheckoutSession, error) {

	params := &stripe.CheckoutSessionParams{
		Customer:           stripe.String(customerID),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)), // 🔑
		LineItems: []*stripe.CheckoutSessionLineItemParams{{
			Price:    stripe.String(priceID), // debe ser price_xxx de tipo one_time
			Quantity: stripe.Int64(1),
		}},
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),

		Metadata: metadata, // viaja en la sesión
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			Metadata: metadata, // viaja también en el PaymentIntent
		},
	}
	return session.New(params)
}

func (s *StripeClient) CreateCustomer(ctx context.Context, email string) (*stripe.Customer, error) {
	params := &stripe.CustomerCreateParams{
		Email: stripe.String(email),
	}
	return s.client.V1Customers.Create(ctx, params)
}

func (s *StripeClient) GetSubscription(ctx context.Context, subscriptionID string) (*stripe.Subscription, error) {
	return subscription.Get(subscriptionID, &stripe.SubscriptionParams{})
}

func GetPaymentIntent(id string) (*stripe.PaymentIntent, error) {
	return paymentintent.Get(id, nil)
}
