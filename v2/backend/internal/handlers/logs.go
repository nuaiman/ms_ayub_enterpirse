package handlers

import (
	"backend/internal/models"
	"backend/internal/utils"
	"net/http"
)

func (h *Handler) GetAllLogsHandler(w http.ResponseWriter, r *http.Request) {
	logs, err := h.app.Models.Log.GetAll(r.Context())
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch logs")
		return
	}

	if logs == nil {
		logs = []models.Log{}
	}

	utils.SuccessJson(w, http.StatusOK, "logs fetched successfully", logs)
}
