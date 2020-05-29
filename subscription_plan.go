package paypal

import (
	"fmt"
	"net/http"
	"time"
)

type (
	// SubscriptionDetailResp struct
	SubscriptionPlan struct {
		ID        string                 `json:"id,omitempty"`
		ProductId string                 `json:"product_id"`
		Name      string                 `json:"name"`
		Status    SubscriptionPlanStatus `json:"status"`
		Description    SubscriptionPlanStatus `json:"description,omitempty"`
		BillingCycles []BillingCycle `json:"billing_cycles"`
		PaymentPreferences PaymentPreferences `json:"payment_preferences"`
		Taxes Taxes `json:"taxes"`
		QuantitySupported bool`json:"quantity_supported"` //Indicates whether you can subscribe to this plan by providing a quantity for the goods or service.
	}

	// Doc https://developer.paypal.com/docs/api/subscriptions/v1/#definition-billing_cycle
	BillingCycle struct{
		PricingScheme PricingScheme `json:"pricing_scheme"` // The active pricing scheme for this billing cycle. A free trial billing cycle does not require a pricing scheme.
		Frequency Frequency `json:"frequency"` // The frequency details for this billing cycle.
		TenureType TenureType `json:"tenure_type"` // The tenure type of the billing cycle. In case of a plan having trial cycle, only 2 trial cycles are allowed per plan. The possible values are:
		Sequence int `json:"sequence"` // The order in which this cycle is to run among other billing cycles. For example, a trial billing cycle has a sequence of 1 while a regular billing cycle has a sequence of 2, so that trial cycle runs before the regular cycle.
		TotalCycles int `json:"total_cycles"` // The number of times this billing cycle gets executed. Trial billing cycles can only be executed a finite number of times (value between 1 and 999 for total_cycles). Regular billing cycles can be executed infinite times (value of 0 for total_cycles) or a finite number of times (value between 1 and 999 for total_cycles).
	}

	// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#definition-payment_preferences
	PaymentPreferences struct{
		AutoBillOutstanding     bool                  `json:"auto_bill_outstanding"`
		SetupFee                Money                 `json:"setup_fee"`
		SetupFeeFailureAction   SetupFeeFailureAction `json:"setup_fee_failure_action"`
		PaymentFailureThreshold int                   `json:"payment_failure_threshold"`
	}

	PricingScheme struct {
		Version int `json:"version"`
		FixedPrice Money `json:"fixed_price"`
		CreateTime time.Time `json:"create_time"`
		UpdateTime time.Time `json:"update_time"`
	}

	Frequency struct {
		IntervalUnit IntervalUnit `json:"interval_unit"`
		IntervalCount int `json:"interval_count"`
	}

	Taxes struct {
		Percentage string `json:"percentage"`
		Inclusive bool `json:"inclusive"`
	}

	CreateSubscriptionPlanResponse struct {
		SubscriptionPlan
		SharedResponse
	}

	SubscriptionPlanListParameters struct {
		ProductId string `json:"product_id"`
		PlanIds string `json:"plan_ids"` // Filters the response by list of plan IDs. Filter supports upto 10 plan IDs.
		ListParams
	}

	ListSubscriptionPlansResponse struct {
		Plans []SubscriptionPlan `json:"plans"`
		ListResponse
}
)

func (self *SubscriptionPlan) GetUpdatePatch() []Patch {
	return []Patch{
		{
			Operation: "replace",
			Path:      "/description",
			Values:    self.Description,
		},
		{
			Operation: "replace",
			Path:      "/payment_preferences/auto_bill_outstanding",
			Values:    self.PaymentPreferences.AutoBillOutstanding,
		},
		{
			Operation: "replace",
			Path:      "/payment_preferences/payment_failure_threshold",
			Values:    self.PaymentPreferences.PaymentFailureThreshold,
		},
		{
			Operation: "replace",
			Path:      "/payment_preferences/setup_fee",
			Values:    self.PaymentPreferences.SetupFee,
		},
		{
			Operation: "replace",
			Path:      "/payment_preferences/setup_fee_failure_action",
			Values:    self.PaymentPreferences.SetupFeeFailureAction,
		},
		{
			Operation: "replace",
			Path:      "/taxes/percentage",
			Values:    self.Taxes.Percentage,
		},
	}
}

type SubscriptionPlanStatus string
const (
	SUBSCRIPTION_PLAN_STATUS_CREATED SubscriptionPlanStatus = "CREATED"
	SUBSCRIPTION_PLAN_STATUS_INACTIVE SubscriptionPlanStatus = "INACTIVE"
	SUBSCRIPTION_PLAN_STATUS_ACTIVE SubscriptionPlanStatus = "ACTIVE"
)

type IntervalUnit string
const (
	INTERVAL_UNIT_DAY IntervalUnit = "DAY"
	INTERVAL_UNIT_WEEK IntervalUnit = "WEEK"
	INTERVAL_UNIT_MONTH IntervalUnit = "MONTH"
	INTERVAL_UNIT_YEAR IntervalUnit = "YEAR"
)

type TenureType string
const (
	TENURE_TYPE_REGULAR TenureType = "REGULAR"
	TENURE_TYPE_TRIAL TenureType = "TRIAL"
)

type SetupFeeFailureAction string
const (
	SETUP_FEE_FAILURE_ACTION_CONTINUE SetupFeeFailureAction = "CONTINUE"
	SETUP_FEE_FAILURE_ACTION_CANCEL SetupFeeFailureAction = "CANCEL"
)

// CreateSubscriptionPlan creates a subscriptionPlan
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#plans_create
// Endpoint: POST /v1/billing/plans
func (c *Client) CreateSubscriptionPlan(newPlan SubscriptionPlan) (*CreateSubscriptionPlanResponse, error) {
	req, err := c.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/plans"), newPlan)
	response := &CreateSubscriptionPlanResponse{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)
	return response, err
}

// UpdateSubscriptionPlan. updates a plan
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#plans_patch
// Endpoint: PATCH /v1/billing/plans/:plan_id
func (c *Client) UpdateSubscriptionPlan(updatedPlan SubscriptionPlan) error {
	req, err := c.NewRequest(http.MethodPatch, fmt.Sprintf("%s%s%s", c.APIBase, "/v1/billing/plans/", updatedPlan.ID), updatedPlan.GetUpdatePatch())
	if err != nil {
		return err
	}
	err = c.SendWithAuth(req, nil)
	return err
}

// UpdateSubscriptionPlan. updates a plan
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#plans_get
// Endpoint: GET /v1/billing/plans/:plan_id
func (c *Client) GetSubscriptionPlan(planId string) (*SubscriptionPlan, error) {
	req, err := c.NewRequest(http.MethodGet, fmt.Sprintf("%s%s%s", c.APIBase, "/v1/billing/plans/", planId), nil)
	response := &SubscriptionPlan{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)
	return response, err
}

// List all plans
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#plans_list
// Endpoint: GET /v1/billing/plans
func (c *Client) ListSubscriptionPlans(params *SubscriptionPlanListParameters) (*ListProductsResponse, error) {
	req, err := c.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/plans"), nil)
	response := &ListProductsResponse{}
	if err != nil {
		return response, err
	}

	if params != nil {
		q := req.URL.Query()
		q.Add("page", params.Page)
		q.Add("page_size", params.PageSize)
		q.Add("total_required", params.TotalRequired)
		q.Add("product_id", params.ProductId)
		q.Add("plan_ids", params.PlanIds)
		req.URL.RawQuery = q.Encode()
	}

	err = c.SendWithAuth(req, response)
	return response, err
}