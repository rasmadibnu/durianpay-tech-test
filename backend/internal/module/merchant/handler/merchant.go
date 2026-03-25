package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/durianpay/fullstack-boilerplate/internal/entity"
	"github.com/durianpay/fullstack-boilerplate/internal/module/merchant/usecase"
	"github.com/durianpay/fullstack-boilerplate/internal/transport"
)

type MerchantHandler struct {
	merchantUC usecase.MerchantUsecase
}

func NewMerchantHandler(uc usecase.MerchantUsecase) *MerchantHandler {
	return &MerchantHandler{merchantUC: uc}
}

type createMerchantRequest struct {
	Name string `json:"name"`
}

type merchantListResponse struct {
	Merchants []entity.Merchant `json:"merchants"`
}

func (h *MerchantHandler) CreateMerchant(w http.ResponseWriter, r *http.Request) {
	var req createMerchantRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}

	merchant, err := h.merchantUC.Create(req.Name)
	if err != nil {
		transport.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(merchant); err != nil {
		transport.WriteAppError(w, entity.ErrorInternal("internal server error"))
	}
}

func (h *MerchantHandler) GetMerchants(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	result, err := h.merchantUC.GetAll(page, limit)
	if err != nil {
		transport.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := map[string]any{
		"merchants":   result.Data,
		"page":        result.Page,
		"limit":       result.Limit,
		"total":       result.Total,
		"total_pages": result.TotalPages,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		transport.WriteAppError(w, entity.ErrorInternal("internal server error"))
	}
}

func (h *MerchantHandler) GetMerchant(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		transport.WriteAppError(w, entity.ErrorBadRequest("invalid merchant id"))
		return
	}

	merchant, err := h.merchantUC.GetByID(id)
	if err != nil {
		transport.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(merchant); err != nil {
		transport.WriteAppError(w, entity.ErrorInternal("internal server error"))
	}
}

func (h *MerchantHandler) UpdateMerchant(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		transport.WriteAppError(w, entity.ErrorBadRequest("invalid merchant id"))
		return
	}

	var req createMerchantRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}

	merchant, err := h.merchantUC.Update(id, req.Name)
	if err != nil {
		transport.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(merchant); err != nil {
		transport.WriteAppError(w, entity.ErrorInternal("internal server error"))
	}
}

func (h *MerchantHandler) DeleteMerchant(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		transport.WriteAppError(w, entity.ErrorBadRequest("invalid merchant id"))
		return
	}

	if err := h.merchantUC.Delete(id); err != nil {
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
