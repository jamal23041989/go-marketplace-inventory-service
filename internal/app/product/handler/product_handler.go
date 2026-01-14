package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jamal23041989/go-marketplace-inventory-service/internal/app/product/service"
	ers "github.com/jamal23041989/go-marketplace-inventory-service/internal/core/errors"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	products, err := h.service.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.respondWithJSON(w, http.StatusOK, products)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := h.getID(r, "/products/")
	if err != nil {
		h.respondWithError(w, err)
		return
	}

	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
			log.Printf("error encoding response: %v", err)
		}
	}
}

func (h *ProductHandler) respondWithError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError

	if errors.Is(err, ers.ErrInvalidInput) {
		status = http.StatusBadRequest
	} else if errors.Is(err, ers.ErrProductNotFound) {
		status = http.StatusNotFound
	}

	http.Error(w, err.Error(), status)
}

func (h *ProductHandler) getID(r *http.Request, path string) (uuid.UUID, error) {
	idStr := strings.TrimPrefix(r.URL.Path, path)
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: invalid uuid format", ers.ErrInvalidInput)
	}
	return id, nil
}
