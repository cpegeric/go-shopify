package goshopify

import (
//	"encoding/json"
	"fmt"
	"net/http"
//	"time"

	"github.com/shopspring/decimal"
)

const payoutsBasePath = "payouts"
const payoutsResourceName = "payouts"

// PayoutService is an interface for interfacing with the payouts endpoints of
// the Shopify API.
// See: https://help.shopify.com/api/reference/shopify_payments
type PayoutService interface {
	List(interface{}) ([]Payout, error)
	ListWithPagination(interface{}) ([]Payout, *Pagination, error)
	Get(int64, interface{}) (*Payout, error)
}

// PayoutServiceOp handles communication with the payout related methods of the
// Shopify API.
type PayoutServiceOp struct {
	client *Client
}

// A struct for all available payout list options.
// See: https://help.shopify.com/api/reference/payout#index
type PayoutListOptions struct {
	SinceId 	  int64     `url:"since_id,omitempty"`
	LastId            int64     `url:"last_id,omitempty"`
	DateMin           string    `url:"date_min,omitempty"`
	DateMax           string    `url:"date_max,omitempty"`
	Date              string    `url:"date,omitempty"`
	Status            string    `url:"status,omitempty"`
}

// Payout represents a Shopify payout
type Payout struct {
	ID                    int64            `json:"id,omitempty"`
	Status                string           `json:"status,omitempty"`
	Date                  string           `json:"date,omitempty"`
	Currency              string           `json:"currency,omitempty"`
	Amount                *decimal.Decimal `json:"amount,omitempty"`
	Summary               *Summary         `json:"summary,omitempty"`
}

type Summary struct {
	AdjustmentsFeeAmount       *decimal.Decimal  `json:"adjustments_fee_amount,omitempty"`
	AdjustmentsGrossAmount     *decimal.Decimal  `json:"adjustments_gross_amount,omitempty"`
	ChargesFeeAmount           *decimal.Decimal  `json:"charges_fee_amount,omitempty"`
	ChargesGrossAmount         *decimal.Decimal  `json:"charges_gross_amount,omitempty"`
	RefundsFeeAmount           *decimal.Decimal  `json:"refunds_fee_amount,omitempty"`
	RefundsGrossAmount         *decimal.Decimal  `json:"refunds_gross_amount,omitempty"`
	ReservedFundsFeeAmount     *decimal.Decimal  `json:"reserved_funds_fee_amount,omitempty"`
	ReservedFundsGrossAmount   *decimal.Decimal  `json:"reserved_funds_gross_amount,omitempty"`
	RetriedPayoutsFeeAmount    *decimal.Decimal  `json:"retried_payouts_fee_amount,omitempty"`
	RetriedPayoutsGrossAmount  *decimal.Decimal  `json:"retried_payouts_gross_amount,omitempty"`
}


// Represents the result from the shopify_payments/payouts/X.json endpoint
type PayoutResource struct {
	Payout *Payout `json:"payout"`
}

// Represents the result from the shopify_payments/payouts.json endpoint
type PayoutsResource struct {
	Payouts []Payout `json:"payouts"`
}


// List payouts
func (s *PayoutServiceOp) List(options interface{}) ([]Payout, error) {
	payouts, _, err := s.ListWithPagination(options)
	if err != nil {
		return nil, err
	}
	return payouts, nil
}

func (s *PayoutServiceOp) ListWithPagination(options interface{}) ([]Payout, *Pagination, error) {
	path := fmt.Sprintf("shopify_payments/%s.json", payoutsBasePath)
	resource := new(PayoutsResource)
	headers := http.Header{}

	headers, err := s.client.createAndDoGetHeaders("GET", path, nil, options, resource)
	if err != nil {
		return nil, nil, err
	}

	// Extract pagination info from header
	linkHeader := headers.Get("Link")

	pagination, err := extractPagination(linkHeader)
	if err != nil {
		return nil, nil, err
	}

	return resource.Payouts, pagination, nil
}


// Get individual payout
func (s *PayoutServiceOp) Get(payoutID int64, options interface{}) (*Payout, error) {
	path := fmt.Sprintf("shopify_payments/%s/%d.json", payoutsBasePath, payoutID)
	resource := new(PayoutResource)
	err := s.client.Get(path, resource, options)
	return resource.Payout, err
}

