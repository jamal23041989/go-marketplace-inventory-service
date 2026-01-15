package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jamal23041989/go-marketplace-inventory-service/internal/app/product/service"
	ers "github.com/jamal23041989/go-marketplace-inventory-service/internal/core/errors"
	"github.com/jamal23041989/go-marketplace-inventory-service/internal/core/logger"
)

type ProductHandler struct {
	service service.ProductService
	logger  logger.Logger
}

func NewProductHandler(service service.ProductService, logger logger.Logger) *ProductHandler {
	return &ProductHandler{
		service: service,
		logger:  logger,
	}
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	if !h.checkMethod(w, r, http.MethodPost) {
		return
	}

	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, fmt.Errorf("%w: invalid body", ers.ErrInvalidInput))
		return
	}

	createProduct, err := h.service.Create(r.Context(), req.ToDomain())
	if err != nil {
		h.respondWithError(w, err)
		return
	}

	h.respondWithJSON(w, http.StatusCreated, createProduct)
}

func (h *ProductHandler) GetById(w http.ResponseWriter, r *http.Request) {
	if !h.checkMethod(w, r, http.MethodGet) {
		return
	}

	id, err := h.getID(r, "/products/")
	if err != nil {
		h.respondWithError(w, err)
		return
	}

	product, err := h.service.GetById(r.Context(), id)
	if err != nil {
		h.respondWithError(w, err)
		return
	}

	h.respondWithJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if !h.checkMethod(w, r, http.MethodGet) {
		return
	}

	products, err := h.service.GetAll(r.Context())
	if err != nil {
		h.respondWithError(w, err)
		return
	}

	h.respondWithJSON(w, http.StatusOK, products)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !h.checkMethod(w, r, http.MethodPatch) {
		return
	}

	id, err := h.getID(r, "/products/")
	if err != nil {
		h.respondWithError(w, err)
		return
	}

	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, err)
		return
	}

	product := req.ToUpdateDTO()
	updateProduct, err := h.service.Update(r.Context(), id, product)
	if err != nil {
		h.respondWithError(w, err)
		return
	}

	h.respondWithJSON(w, http.StatusOK, updateProduct)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if !h.checkMethod(w, r, http.MethodDelete) {
		return
	}

	id, err := h.getID(r, "/products/")
	if err != nil {
		h.respondWithError(w, err)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		h.respondWithError(w, err)
		return
	}

	h.respondWithJSON(w, http.StatusNoContent, nil)
}

func (h *ProductHandler) respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			h.logger.Error("error encoding response: %v", err)
		}
	}
}

func (h *ProductHandler) respondWithError(w http.ResponseWriter, err error) {
	var valErr *ers.ValidationError

	if errors.As(err, &valErr) {
		h.logger.Warn("validation error: %v", valErr.Error())
		http.Error(w, valErr.Error(), http.StatusBadRequest)
		return
	} else if errors.Is(err, ers.ErrProductNotFound) {
		h.logger.Warn("product not found: %v", ers.ErrProductNotFound)
		http.Error(w, fmt.Errorf("%w: not found error", ers.ErrProductNotFound).Error(), http.StatusNotFound)
		return
	}

	h.logger.Error("internal error: %v", err)
	http.Error(w, fmt.Errorf("%w: internal error", ers.ErrInternalServerError).Error(), http.StatusInternalServerError)
	return
}

func (h *ProductHandler) getID(r *http.Request, path string) (uuid.UUID, error) {
	idStr := strings.TrimPrefix(r.URL.Path, path)
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Warn("invalid id: %v", idStr)
		return uuid.Nil, fmt.Errorf("%w: invalid uuid format", ers.ErrInvalidInput)
	}
	return id, nil
}

func (h *ProductHandler) checkMethod(w http.ResponseWriter, r *http.Request, expectedMethod string) bool {
	if r.Method != expectedMethod {
		h.logger.Warn("method not allowed: %v", expectedMethod)
		http.Error(w, fmt.Errorf("%w: method not allowed error", ers.ErrMethodNotAllowed).Error(), http.StatusMethodNotAllowed)
		return false
	}
	return true
}
