package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"backend/internal/middlewares"
	"backend/internal/utils"
)

// =======================================
// UPLOAD IMAGE (USERS / ITEMS / EXPENSES)
// POST /images/{type}/{id}
// =======================================
func (h *Handler) UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	entityType := r.PathValue("type")

	id, ok := utils.GetParamID(w, r)
	if !ok {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid id")
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid form")
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "image required")
		return
	}
	defer file.Close()

	filename := "img_" + header.Filename

	var (
		basePath     string
		imageURL     string
		oldImagePath string
	)

	switch entityType {

	case "users":
		user, err := h.app.Models.User.GetByID(r.Context(), id)
		if err != nil || user == nil {
			utils.ErrorJson(w, http.StatusNotFound, "user not found")
			return
		}

		if user.ID != userID {
			currentUser, err := h.app.Models.User.GetByID(r.Context(), userID)
			if err != nil || currentUser == nil {
				utils.ErrorJson(w, http.StatusForbidden, "unauthorized")
				return
			}
			if currentUser.Role != "admin" && currentUser.Role != "manager" && currentUser.Role != "accounts" {
				utils.ErrorJson(w, http.StatusForbidden, "unauthorized")
				return
			}
		}

		if user.ImageURL != nil && *user.ImageURL != "" {
			oldImagePath = "." + *user.ImageURL
		}

		basePath = filepath.Join("public", "users", strconv.FormatInt(id, 10))

	case "items":
		item, err := h.app.Models.Item.GetByID(r.Context(), id)
		if err != nil || item == nil {
			utils.ErrorJson(w, http.StatusNotFound, "item not found")
			return
		}

		if item.UserID == nil || *item.UserID != userID {
			currentUser, err := h.app.Models.User.GetByID(r.Context(), userID)
			if err != nil || currentUser == nil {
				utils.ErrorJson(w, http.StatusForbidden, "unauthorized")
				return
			}
			if currentUser.Role != "admin" && currentUser.Role != "manager" && currentUser.Role != "accounts" {
				utils.ErrorJson(w, http.StatusForbidden, "unauthorized")
				return
			}
		}

		if item.ImageURL != nil && *item.ImageURL != "" {
			oldImagePath = "." + *item.ImageURL
		}

		basePath = filepath.Join("public", "items", strconv.FormatInt(id, 10))

	case "expenses":
		expense, err := h.app.Models.Expense.GetByID(r.Context(), id)
		if err != nil || expense == nil {
			utils.ErrorJson(w, http.StatusNotFound, "expense not found")
			return
		}

		if expense.UserID == nil || *expense.UserID != userID {
			currentUser, err := h.app.Models.User.GetByID(r.Context(), userID)
			if err != nil || currentUser == nil {
				utils.ErrorJson(w, http.StatusForbidden, "unauthorized")
				return
			}
			if currentUser.Role != "admin" && currentUser.Role != "manager" && currentUser.Role != "accounts" {
				utils.ErrorJson(w, http.StatusForbidden, "unauthorized")
				return
			}
		}

		if expense.ImageURL != nil && *expense.ImageURL != "" {
			oldImagePath = "." + *expense.ImageURL
		}

		basePath = filepath.Join("public", "expenses", strconv.FormatInt(id, 10))

	default:
		utils.ErrorJson(w, http.StatusBadRequest, "invalid type")
		return
	}

	err = os.MkdirAll(basePath, 0755)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create directory")
		return
	}

	dstPath := filepath.Join(basePath, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to save file")
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to write file")
		os.Remove(dstPath)
		return
	}

	imageURL = "/" + filepath.ToSlash(dstPath)

	switch entityType {
	case "users":
		err = h.app.Models.User.UpdateImage(r.Context(), id, &imageURL)
	case "items":
		err = h.app.Models.Item.UpdateImage(r.Context(), id, &imageURL)
	case "expenses":
		err = h.app.Models.Expense.UpdateImage(r.Context(), id, &imageURL)
	}

	if err != nil {
		os.Remove(dstPath)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update database")
		return
	}

	if oldImagePath != "" {
		if err := os.Remove(oldImagePath); err != nil {
			println("Warning: Failed to delete old image:", oldImagePath, err.Error())
		}
	}

	utils.SuccessJson(w, http.StatusOK, "image uploaded", map[string]any{
		"image_url": imageURL,
	})
}

// =======================================
// DELETE IMAGE (USERS / ITEMS / EXPENSES)
// DELETE /images/{type}/{id}
// =======================================
func (h *Handler) DeleteImageHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	entityType := r.PathValue("type")

	id, ok := utils.GetParamID(w, r)
	if !ok {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid id")
		return
	}

	var imageURL string

	switch entityType {

	case "users":
		user, err := h.app.Models.User.GetByID(r.Context(), id)
		if err != nil || user == nil {
			utils.ErrorJson(w, http.StatusNotFound, "user not found")
			return
		}

		if user.ID != userID {
			currentUser, err := h.app.Models.User.GetByID(r.Context(), userID)
			if err != nil || currentUser == nil {
				utils.ErrorJson(w, http.StatusForbidden, "unauthorized")
				return
			}
			if currentUser.Role != "admin" && currentUser.Role != "manager" && currentUser.Role != "accounts" {
				utils.ErrorJson(w, http.StatusForbidden, "unauthorized")
				return
			}
		}

		if user.ImageURL != nil {
			imageURL = *user.ImageURL
		}

	case "items":
		item, err := h.app.Models.Item.GetByID(r.Context(), id)
		if err != nil || item == nil {
			utils.ErrorJson(w, http.StatusNotFound, "item not found")
			return
		}

		if item.UserID == nil || *item.UserID != userID {
			currentUser, err := h.app.Models.User.GetByID(r.Context(), userID)
			if err != nil || currentUser == nil {
				utils.ErrorJson(w, http.StatusForbidden, "unauthorized")
				return
			}
			if currentUser.Role != "admin" && currentUser.Role != "manager" && currentUser.Role != "accounts" {
				utils.ErrorJson(w, http.StatusForbidden, "unauthorized")
				return
			}
		}

		if item.ImageURL != nil {
			imageURL = *item.ImageURL
		}

	case "expenses":
		expense, err := h.app.Models.Expense.GetByID(r.Context(), id)
		if err != nil || expense == nil {
			utils.ErrorJson(w, http.StatusNotFound, "expense not found")
			return
		}

		if expense.UserID == nil || *expense.UserID != userID {
			currentUser, err := h.app.Models.User.GetByID(r.Context(), userID)
			if err != nil || currentUser == nil {
				utils.ErrorJson(w, http.StatusForbidden, "unauthorized")
				return
			}
			if currentUser.Role != "admin" && currentUser.Role != "manager" && currentUser.Role != "accounts" {
				utils.ErrorJson(w, http.StatusForbidden, "unauthorized")
				return
			}
		}

		if expense.ImageURL != nil {
			imageURL = *expense.ImageURL
		}

	default:
		utils.ErrorJson(w, http.StatusBadRequest, "invalid type")
		return
	}

	if imageURL == "" {
		utils.SuccessJson(w, http.StatusOK, "no image to delete", nil)
		return
	}

	_ = os.Remove("." + imageURL)

	switch entityType {
	case "users":
		err := h.app.Models.User.UpdateImage(r.Context(), id, nil)
		if err != nil {
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to update database")
			return
		}
	case "items":
		err := h.app.Models.Item.UpdateImage(r.Context(), id, nil)
		if err != nil {
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to update database")
			return
		}
	case "expenses":
		err := h.app.Models.Expense.UpdateImage(r.Context(), id, nil)
		if err != nil {
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to update database")
			return
		}
	}

	utils.SuccessJson(w, http.StatusOK, "image deleted", nil)
}
