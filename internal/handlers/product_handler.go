package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/habtamubaye-devsec/shoe-shop/internal/dto"
	"github.com/habtamubaye-devsec/shoe-shop/internal/services"
)

type ProductHandler struct {
	service services.ProductService
}

func NewProductHandler(service services.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	response APIResponse,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}

func writeError(
	w http.ResponseWriter,
	status int,
	message string,
) {
	writeJSON(w, status, APIResponse{
		Success: false,
		Message: message,
	})
}

func decodeJSON(w http.ResponseWriter, r *http.Request, destination interface{}) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(destination); err != nil {
		return err
	}

	var extra interface{}
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return errors.New("request body must contain only one JSON object")
		}
		return err
	}

	return nil
}

func parseID(r *http.Request) (uint, error) {
	rawID := strings.TrimPrefix(r.URL.Path, "/api/v1/products/")

	if rawID == "" || strings.Contains(rawID, "/") {
		return 0, errors.New("invalid product ID")
	}

	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id == 0 {
		return 0, errors.New("invalid product ID")
	}

	return uint(id), nil
}

func (h *ProductHandler) Products(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listProducts(w, r)
	case http.MethodPost:
		h.createProduct(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *ProductHandler) ProductByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getProduct(w, r, id)
	case http.MethodPut:
		h.updateProduct(w, r, id)
	case http.MethodDelete:
		h.deleteProduct(w, r, id)
	default:
		w.Header().Set("Allow", "GET, PUT, DELETE")
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *ProductHandler) createProduct(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateProductRequest

	if err := decodeJSON(w, r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON request body")
		return
	}

	product, err := h.service.Create(r.Context(), input)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Message: "product created successfully",
		Data:    product,
	})
}

func (h *ProductHandler) getProduct(
	w http.ResponseWriter,
	r *http.Request,
	id uint,
) {
	product, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "product retrieved successfully",
		Data:    product,
	})
}

func (h *ProductHandler) listProducts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	page, err := parsePositiveInt(query.Get("page"), 1)
	if err != nil {
		writeError(w, http.StatusBadRequest, "page must be a positive integer")
		return
	}

	limit, err := parsePositiveInt(query.Get("limit"), 10)
	if err != nil {
		writeError(w, http.StatusBadRequest, "limit must be a positive integer")
		return
	}

	categoryID, err := parseOptionalUint(query.Get("category_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "category_id must be a positive integer")
		return
	}

	minPrice, err := parseOptionalInt64(query.Get("min_price_minor"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "min_price_minor must be an integer")
		return
	}

	maxPrice, err := parseOptionalInt64(query.Get("max_price_minor"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "max_price_minor must be an integer")
		return
	}

	result, err := h.service.List(r.Context(), services.ProductListParams{
		Search:     query.Get("search"),
		CategoryID: categoryID,
		MinPrice:   minPrice,
		MaxPrice:   maxPrice,
		Page:       page,
		Limit:      limit,
	})

	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "products retrieved successfully",
		Data:    result,
	})
}

func (h *ProductHandler) updateProduct(
	w http.ResponseWriter,
	r *http.Request,
	id uint,
) {
	var input dto.UpdateProductRequest

	if err := decodeJSON(w, r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON request body")
		return
	}

	product, err := h.service.Update(r.Context(), id, input)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "product updated successfully",
		Data:    product,
	})
}

func (h *ProductHandler) deleteProduct(
	w http.ResponseWriter,
	r *http.Request,
	id uint,
) {
	if err := h.service.Delete(r.Context(), id); err != nil {
		h.handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "product deleted successfully",
	})
}

func (h *ProductHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())

	case errors.Is(err, services.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, err.Error())

	default:
		// Log the real error on the server, but don't expose database details.
		fmt.Printf("Product API error: %v\n", err)
		writeError(w, http.StatusInternalServerError, "an internal server error occurred")
	}
}

func parsePositiveInt(raw string, defaultValue int) (int, error) {
	if raw == "" {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value < 1 {
		return 0, errors.New("value must be positive")
	}

	return value, nil
}

func parseOptionalUint(raw string) (uint, error) {
	if raw == "" {
		return 0, nil
	}

	value, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || value == 0 {
		return 0, errors.New("value must be positive")
	}

	return uint(value), nil
}

func parseOptionalInt64(raw string) (*int64, error) {
	if raw == "" {
		return nil, nil
	}

	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return nil, errors.New("value must be an integer")
	}

	return &value, nil
}
