package utils

import (
	"net/http"
	"strconv"
)

func GetParamID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ErrorJson(w, http.StatusBadRequest, "invalid parameter")
		return 0, false
	}

	return id, true
}
