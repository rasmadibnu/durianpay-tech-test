package api

import (
	"net/http"

	ah "github.com/durianpay/fullstack-boilerplate/internal/module/auth/handler"
	mh "github.com/durianpay/fullstack-boilerplate/internal/module/merchant/handler"
	"github.com/durianpay/fullstack-boilerplate/internal/openapigen"
)

type APIHandler struct {
	Auth     *ah.AuthHandler
	Merchant *mh.MerchantHandler
}

var _ openapigen.ServerInterface = (*APIHandler)(nil)

func (h *APIHandler) PostDashboardV1AuthLogin(w http.ResponseWriter, r *http.Request) {
	h.Auth.PostDashboardV1AuthLogin(w, r)
}

// Merchants
func (h *APIHandler) GetDashboardV1Merchants(w http.ResponseWriter, r *http.Request) {
	h.Merchant.GetMerchants(w, r)
}

func (h *APIHandler) PostDashboardV1Merchants(w http.ResponseWriter, r *http.Request) {
	h.Merchant.CreateMerchant(w, r)
}

func (h *APIHandler) GetDashboardV1MerchantsId(w http.ResponseWriter, r *http.Request, id int) {
	h.Merchant.GetMerchant(w, r)
}

func (h *APIHandler) PutDashboardV1MerchantsId(w http.ResponseWriter, r *http.Request, id int) {
	h.Merchant.UpdateMerchant(w, r)
}

func (h *APIHandler) DeleteDashboardV1MerchantsId(w http.ResponseWriter, r *http.Request, id int) {
	h.Merchant.DeleteMerchant(w, r)
}

func (h *APIHandler) GetDashboardV1Payments(w http.ResponseWriter, r *http.Request, body openapigen.GetDashboardV1PaymentsParams) {
	// TODO add implementations
}
