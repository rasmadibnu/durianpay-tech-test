package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/durianpay/fullstack-boilerplate/internal/entity"
	"github.com/durianpay/fullstack-boilerplate/internal/module/user/usecase"
	"github.com/durianpay/fullstack-boilerplate/internal/transport"
)

type UserHandler struct {
	userUC usecase.UserUsecase
}

func NewUserHandler(uc usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUC: uc}
}

type createUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type updateUserRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type updatePasswordRequest struct {
	Password string `json:"password"`
}

type userResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func toUserResponse(u *entity.User) userResponse {
	return userResponse{ID: u.ID, Email: u.Email, Role: u.Role}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
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

	var search *string
	if s := r.URL.Query().Get("search"); s != "" {
		search = &s
	}

	result, err := h.userUC.GetAll(search, page, limit)
	if err != nil {
		transport.WriteError(w, err)
		return
	}

	resp := make([]userResponse, len(result.Data))
	for i, u := range result.Data {
		resp[i] = userResponse{ID: u.ID, Email: u.Email, Role: u.Role}
	}

	w.Header().Set("Content-Type", "application/json")
	out := map[string]any{
		"users":       resp,
		"page":        result.Page,
		"limit":       result.Limit,
		"total":       result.Total,
		"total_pages": result.TotalPages,
	}
	if err := json.NewEncoder(w).Encode(out); err != nil {
		transport.WriteAppError(w, entity.ErrorInternal("internal server error"))
	}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		transport.WriteAppError(w, entity.ErrorBadRequest("invalid user id"))
		return
	}

	user, err := h.userUC.GetByID(id)
	if err != nil {
		transport.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(toUserResponse(user)); err != nil {
		transport.WriteAppError(w, entity.ErrorInternal("internal server error"))
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}

	user, err := h.userUC.Create(req.Email, req.Password, req.Role)
	if err != nil {
		transport.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(toUserResponse(user)); err != nil {
		transport.WriteAppError(w, entity.ErrorInternal("internal server error"))
	}
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		transport.WriteAppError(w, entity.ErrorBadRequest("invalid user id"))
		return
	}

	var req updateUserRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}

	user, err := h.userUC.Update(id, req.Email, req.Role)
	if err != nil {
		transport.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(toUserResponse(user)); err != nil {
		transport.WriteAppError(w, entity.ErrorInternal("internal server error"))
	}
}

func (h *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		transport.WriteAppError(w, entity.ErrorBadRequest("invalid user id"))
		return
	}

	var req updatePasswordRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}

	if err := h.userUC.UpdatePassword(id, req.Password); err != nil {
		transport.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "password updated"}); err != nil {
		transport.WriteAppError(w, entity.ErrorInternal("internal server error"))
	}
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		transport.WriteAppError(w, entity.ErrorBadRequest("invalid user id"))
		return
	}

	if err := h.userUC.Delete(id); err != nil {
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
