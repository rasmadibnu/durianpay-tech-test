package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/durianpay/fullstack-boilerplate/internal/entity"
	"github.com/durianpay/fullstack-boilerplate/internal/module/payment/usecase"
	"github.com/durianpay/fullstack-boilerplate/internal/openapigen"
	"github.com/durianpay/fullstack-boilerplate/internal/transport"
)

type PaymentHandler struct {
	paymentUC usecase.PaymentUsecase
}

func NewPaymentHandler(uc usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{paymentUC: uc}
}

type createPaymentRequest struct {
	MerchantID int    `json:"merchant_id"`
	Amount     string `json:"amount"`
	Status     string `json:"status"`
}

type updatePaymentRequest struct {
	MerchantID int    `json:"merchant_id"`
	Amount     string `json:"amount"`
	Status     string `json:"status"`
}

type updatePaymentStatusRequest struct {
	Status string `json:"status"`
}

func (h *PaymentHandler) GetDashboardV1Payments(w http.ResponseWriter, r *http.Request, params openapigen.GetDashboardV1PaymentsParams) {
	page := 1
	limit := 10
	if params.Page != nil && *params.Page > 0 {
		page = *params.Page
	}
	if params.Limit != nil && *params.Limit > 0 {
		limit = *params.Limit
		if limit > 100 {
			limit = 100
		}
	}

	var search *string
	if s := r.URL.Query().Get("search"); s != "" {
		search = &s
	}

	result, err := h.paymentUC.GetPayments(params.Status, params.Id, params.Sort, search, page, limit)
	if err != nil {
		transport.WriteError(w, err)
		return
	}

	oaPayments := make([]openapigen.Payment, len(result.Data))
	for i, p := range result.Data {
		id := p.ID
		merchantId := p.MerchantID
		merchantName := p.MerchantName
		amount := p.Amount
		status := p.Status
		createdAt := p.CreatedAt
		oaPayments[i] = openapigen.Payment{
			Id:           &id,
			MerchantId:   &merchantId,
			MerchantName: &merchantName,
			Amount:       &amount,
			Status:       &status,
			CreatedAt:    &createdAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	resp := map[string]any{
		"payments":    oaPayments,
		"page":        result.Page,
		"limit":       result.Limit,
		"total":       result.Total,
		"total_pages": result.TotalPages,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		transport.WriteAppError(w, entity.ErrorInternal("internal server error"))
	}
}

func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var req createPaymentRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}

	payment, err := h.paymentUC.CreatePayment(req.MerchantID, req.Amount, req.Status)
	if err != nil {
		transport.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(payment); err != nil {
		transport.WriteAppError(w, entity.ErrorInternal("internal server error"))
	}
}

func (h *PaymentHandler) UpdatePayment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req updatePaymentRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}

	payment, err := h.paymentUC.UpdatePayment(id, req.MerchantID, req.Amount, req.Status)
	if err != nil {
		transport.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payment); err != nil {
		transport.WriteAppError(w, entity.ErrorInternal("internal server error"))
	}
}

func (h *PaymentHandler) UpdatePaymentStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req updatePaymentStatusRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}

	if err := h.paymentUC.UpdatePaymentStatus(id, req.Status); err != nil {
		transport.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "status updated"}); err != nil {
		transport.WriteAppError(w, entity.ErrorInternal("internal server error"))
	}
}

func (h *PaymentHandler) DeletePayment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.paymentUC.DeletePayment(id); err != nil {
		transport.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) bool {
	if r.Body == nil {
		transport.WriteAppError(w, entity.ErrorBadRequest("empty body"))
		return false
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		transport.WriteAppError(w, entity.ErrorBadRequest("failed to read body"))
		return false
	}

	if err := json.Unmarshal(body, dst); err != nil {
		transport.WriteAppError(w, entity.ErrorBadRequest("invalid json: "+err.Error()))
		return false
	}
	return true
}
