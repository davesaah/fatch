package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/davesaah/fatch/internal/database"
	"github.com/davesaah/fatch/types"
	"github.com/jackc/pgx/v5/pgtype"
)

// GetCategories godoc
// @Summary Get all categories
// @Description Get all categories
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {object} types.SuccessResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /categories [get]
func (h *Handler) GetCategories(w http.ResponseWriter, r *http.Request) *types.ErrorDetails {
	if r.URL.Query().Get("category_id") != "" {
		return h.GetCategoryByID(w, r)
	}

	userID := r.Context().Value("userID").(pgtype.UUID)
	categories, errResponse, err := h.Service.GetCategories(r.Context(), userID)
	if err != nil {
		types.ReturnJSON(w, errResponse)
		return &types.ErrorDetails{
			Message: "Unable to get all categories",
			Level:   "DEBUG",
			Trace:   err,
		}
	}

	return types.ReturnJSON(w, types.OKResponse("", categories))
}

// GetCategoryByID godoc
// @Summary Get category by ID
// @Description Get category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param category_id query int true "Category ID"
// @Success 200 {object} types.SuccessResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /categories [get]
func (h *Handler) GetCategoryByID(w http.ResponseWriter, r *http.Request) *types.ErrorDetails {
	userID := r.Context().Value("userID").(pgtype.UUID)
	idStr := r.URL.Query().Get("category_id")
	categoryID, err := strconv.Atoi(idStr)
	if err != nil {
		types.ReturnJSON(w, types.BadRequestErrorResponse("Invalid category id"))
		return &types.ErrorDetails{
			Message: "Unable to convert id to integer",
			Level:   "ERROR",
			Trace:   err,
		}
	}

	params := database.GetCategoryByIDParams{
		UserID:     userID,
		CategoryID: categoryID,
	}

	category, errResponse, err := h.Service.GetCategoryByID(r.Context(), params)
	if err != nil {
		types.ReturnJSON(w, errResponse)
		return &types.ErrorDetails{
			Message: "Unable to get category by id",
			Level:   "DEBUG",
			Trace:   err,
		}
	}

	return types.ReturnJSON(w, types.OKResponse("", category))
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new category
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {object} types.SuccessResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /categories [post]
func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) *types.ErrorDetails {
	var params database.CreateCategoryParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		types.ReturnJSON(w, types.BadRequestErrorResponse("Invalid request body"))
		return &types.ErrorDetails{
			Message: "Unable to decode request body",
			Level:   "ERROR",
			Trace:   err,
		}
	}

	// no empty values
	if params.CategoryName == "" || params.Type == "" {
		types.ReturnJSON(w, types.BadRequestErrorResponse("No empty fields allowed"))
		return &types.ErrorDetails{
			Message: "No empty fields allowed",
			Level:   "ERROR",
		}
	}

	// type must be "expense" or "income"
	if params.Type != "expense" && params.Type != "income" {
		types.ReturnJSON(w, types.BadRequestErrorResponse("Invalid category type"))
		return &types.ErrorDetails{
			Message: "Invalid category type",
			Level:   "ERROR",
		}
	}

	params.UserID = r.Context().Value("userID").(pgtype.UUID)

	if errResponse, err := h.Service.AddCategory(r.Context(), params); err != nil {
		types.ReturnJSON(w, errResponse)
		return &types.ErrorDetails{
			Message: "Unable to add category",
			Level:   "DEBUG",
			Trace:   err,
		}
	}

	return types.ReturnJSON(w, types.OKResponse("Category added successfully", nil))
}

// UpdateCategory godoc
// @Summary Update a category
// @Description Update a category
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {object} types.SuccessResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /categories [put]
func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) *types.ErrorDetails {
	return nil
}

// DeleteCategory godoc
// @Summary Delete a category
// @Description Delete a category
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {object} types.SuccessResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /categories [delete]
func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) *types.ErrorDetails {
	return nil
}
