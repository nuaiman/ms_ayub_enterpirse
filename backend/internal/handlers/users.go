package handlers

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"backend/internal/utils"
	"context"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Username == "" || req.Password == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "missing required fields")
		return
	}

	user, err := h.app.Models.User.GetByUsername(r.Context(), req.Username)
	if err != nil || user == nil {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if user.IsActive == false {
		utils.ErrorJson(w, http.StatusForbidden, "account is deactivated")
		return
	}

	if !utils.ComparePassword(user.Password, req.Password) {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "could not generate session")
		return
	}

	if err := h.app.Models.User.UpdateRefreshToken(r.Context(), user.ID, &refreshToken); err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "could not save session")
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, h.app.Config.JWTKey)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "could not generate access token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 30,
	})

	user.Password = ""

	utils.SuccessJson(w, http.StatusOK, "login successful", map[string]any{
		"access_token": accessToken,
	})
}

func (h *Handler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil || cookie.Value == "" {
		utils.ErrorJson(w, http.StatusUnauthorized, "missing session")
		return
	}

	user, err := h.app.Models.User.GetByRefreshToken(r.Context(), cookie.Value)
	if err != nil || user == nil {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid session")
		return
	}

	newRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "could not generate session")
		return
	}

	if err := h.app.Models.User.UpdateRefreshToken(r.Context(), user.ID, &newRefreshToken); err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "could not update session")
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, h.app.Config.JWTKey)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "could not generate access token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 30,
	})

	utils.SuccessJson(w, http.StatusOK, "session refreshed", map[string]any{
		"access_token": accessToken,
	})
}

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	if err := h.app.Models.User.UpdateRefreshToken(r.Context(), userID, nil); err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to logout")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	utils.SuccessJson(w, http.StatusOK, "logged out successfully", nil)
}

func (h *Handler) GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	user, err := h.app.Models.User.GetByID(r.Context(), userID)
	if err != nil || user == nil {
		utils.ErrorJson(w, http.StatusUnauthorized, "user not found")
		return
	}

	user.Password = ""

	utils.SuccessJson(w, http.StatusOK, "user details fetched", user)
}

func (h *Handler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetParamID(w, r)
	if !ok {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid parameter")
		return
	}

	user, err := h.app.Models.User.GetByID(r.Context(), userID)
	if err != nil || user == nil {
		utils.ErrorJson(w, http.StatusNotFound, "user not found")
		return
	}

	user.Password = ""

	utils.SuccessJson(w, http.StatusOK, "user details fetched", user)
}

func (h *Handler) ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.CurrentPassword == "" || req.NewPassword == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "missing required fields")
		return
	}

	if err := utils.ValidatePassword(req.NewPassword); err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.CurrentPassword == req.NewPassword {
		utils.ErrorJson(w, http.StatusBadRequest, "new password must be different from current password")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	user, err := h.app.Models.User.GetByID(r.Context(), userID)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}

	if user == nil {
		utils.ErrorJson(w, http.StatusNotFound, "user not found")
		return
	}

	if !utils.ComparePassword(user.Password, req.CurrentPassword) {
		utils.ErrorJson(w, http.StatusUnauthorized, "current password is incorrect")
		return
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	err = h.app.Models.User.UpdatePassword(r.Context(), userID, hashedPassword)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update password")
		return
	}

	// Audit log
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Changed own password",
		EntityType:  "users",
		EntityID:    userID,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "password changed successfully", nil)
}

func (h *Handler) UpdateUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Name     string  `json:"name"`
		Email    *string `json:"email"`
		Phone    *string `json:"phone"`
		Address  *string `json:"address"`
		IDType   *string `json:"id_type"`
		IDNumber *string `json:"id_number"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "name is required")
		return
	}

	currentUserID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	targetUserID, ok := utils.GetParamID(w, r)
	if !ok {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid user id")
		return
	}

	targetUser, err := h.app.Models.User.GetByID(r.Context(), targetUserID)
	if err != nil || targetUser == nil {
		utils.ErrorJson(w, http.StatusNotFound, "user not found")
		return
	}
	if targetUser.Role == "admin" && targetUserID != currentUserID {
		utils.ErrorJson(w, http.StatusForbidden, "cannot modify admin users")
		return
	}

	if targetUserID != currentUserID {
		currentUser, err := h.app.Models.User.GetByID(r.Context(), currentUserID)
		if err != nil || currentUser == nil {
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify permissions")
			return
		}
		if currentUser.Role != "admin" && currentUser.Role != "manager" {
			utils.ErrorJson(w, http.StatusForbidden, "you can only update your own profile")
			return
		}
	}

	if req.Email != nil && *req.Email != "" {
		if !strings.Contains(*req.Email, "@") {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid email format")
			return
		}

		existingUser, err := h.app.Models.User.GetByEmail(r.Context(), *req.Email)
		if err == nil && existingUser != nil && existingUser.ID != targetUserID {
			utils.ErrorJson(w, http.StatusConflict, "email already in use")
			return
		}
	}

	if req.IDType != nil && *req.IDType != "" {
		validTypes := map[string]bool{
			"nid": true, "passport": true, "driving_license": true,
			"birth_certificate": true, "trade_license": true, "other": true,
		}
		if !validTypes[*req.IDType] {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid id_type")
			return
		}
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Address:  req.Address,
		IDType:   req.IDType,
		IDNumber: req.IDNumber,
	}

	err = h.app.Models.User.UpdateProfile(r.Context(), targetUserID, user)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update user")
		return
	}

	updatedUser, err := h.app.Models.User.GetByID(r.Context(), targetUserID)
	if err != nil || updatedUser == nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch updated user")
		return
	}

	updatedUser.Password = ""

	utils.SuccessJson(w, http.StatusOK, "user updated successfully", updatedUser)
}

func (h *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Name     string  `json:"name"`
		Username string  `json:"username"`
		Password string  `json:"password"`
		Role     string  `json:"role"`
		Email    *string `json:"email"`
		Phone    *string `json:"phone"`
		Address  *string `json:"address"`
		IDType   *string `json:"id_type"`
		IDNumber *string `json:"id_number"`
		ImageURL *string `json:"image_url"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" || req.Username == "" || req.Password == "" || req.Role == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "missing required fields")
		return
	}

	switch req.Role {
	case "manager", "staff", "accounts":
	default:
		utils.ErrorJson(w, http.StatusBadRequest, "invalid role")
		return
	}

	if req.IDType != nil && *req.IDType != "" {
		validTypes := map[string]bool{
			"nid": true, "passport": true, "driving_license": true,
			"birth_certificate": true, "trade_license": true, "other": true,
		}
		if !validTypes[*req.IDType] {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid id_type")
			return
		}
	}

	if req.Email != nil && *req.Email != "" {
		if !strings.Contains(*req.Email, "@") {
			utils.ErrorJson(w, http.StatusBadRequest, "invalid email format")
			return
		}

		existingUser, err := h.app.Models.User.GetByEmail(r.Context(), *req.Email)
		if err == nil && existingUser != nil {
			utils.ErrorJson(w, http.StatusConflict, "email already in use")
			return
		}
	}

	if err := utils.ValidatePassword(req.Password); err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	existing, err := h.app.Models.User.GetByUsername(r.Context(), req.Username)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to check username")
		return
	}

	if existing != nil {
		utils.ErrorJson(w, http.StatusConflict, "username already exists")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	user := models.User{
		Name:     req.Name,
		Username: req.Username,
		Password: hashedPassword,
		Role:     req.Role,
		Email:    req.Email,
		Phone:    req.Phone,
		Address:  req.Address,
		IDType:   req.IDType,
		IDNumber: req.IDNumber,
		ImageURL: req.ImageURL,
	}

	id, err := h.app.Models.User.Insert(r.Context(), &user)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	user.ID = id
	user.Password = ""

	// Audit log
	currentUserID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      currentUserID,
		Action:      "create",
		Description: "Created user: " + req.Username,
		EntityType:  "users",
		EntityID:    id,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusCreated, "user created successfully", user)
}

func (h *Handler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.app.Models.User.GetAll(r.Context())
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch users")
		return
	}

	for i := range users {
		users[i].Password = ""
		users[i].RefreshToken = nil
	}

	utils.SuccessJson(w, http.StatusOK, "users fetched successfully", users)
}

func (h *Handler) ResetAllPasswordsHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		NewPassword string `json:"new_password"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.NewPassword == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "missing new_password")
		return
	}

	if err := utils.ValidatePassword(req.NewPassword); err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	hashed, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	err = h.app.Models.User.ResetAllPasswordsExceptRoot(r.Context(), hashed)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to reset passwords")
		return
	}

	// Audit log
	userID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      userID,
		Action:      "update",
		Description: "Reset all user passwords",
		EntityType:  "users",
		EntityID:    0,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "all user passwords reset successfully", nil)
}

func (h *Handler) ChangeUserPasswordHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		NewPassword string `json:"new_password"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.NewPassword == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "missing new_password")
		return
	}

	if err := utils.ValidatePassword(req.NewPassword); err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	userID, ok := utils.GetParamID(w, r)
	if !ok {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.app.Models.User.GetByID(r.Context(), userID)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}

	if user == nil {
		utils.ErrorJson(w, http.StatusNotFound, "user not found")
		return
	}

	if user.Role == "admin" {
		utils.ErrorJson(w, http.StatusForbidden, "cannot modify admin users")
		return
	}

	hashed, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	err = h.app.Models.User.AdminUpdatePasswordAndInvalidateSessions(r.Context(), userID, hashed)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update password")
		return
	}

	// Audit log
	currentUserID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      currentUserID,
		Action:      "update",
		Description: "Admin changed password for user #" + strconv.FormatInt(userID, 10),
		EntityType:  "users",
		EntityID:    userID,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "user password updated successfully", nil)
}

func (h *Handler) ChangeUserRoleHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Role string `json:"role"`
	}

	var req request

	if err := utils.ReadJson(w, r, &req); err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Role == "" {
		utils.ErrorJson(w, http.StatusBadRequest, "missing role")
		return
	}

	switch req.Role {
	case "admin", "manager", "accounts", "staff":
	default:
		utils.ErrorJson(w, http.StatusBadRequest, "invalid role")
		return
	}

	userID, ok := utils.GetParamID(w, r)
	if !ok {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.app.Models.User.GetByID(r.Context(), userID)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}

	if user == nil {
		utils.ErrorJson(w, http.StatusNotFound, "user not found")
		return
	}

	if user.Role == "admin" {
		utils.ErrorJson(w, http.StatusForbidden, "cannot modify admin users")
		return
	}

	err = h.app.Models.User.UpdateRole(r.Context(), userID, req.Role)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update role")
		return
	}

	// Audit log
	currentUserID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      currentUserID,
		Action:      "update",
		Description: "Changed role for user #" + strconv.FormatInt(userID, 10) + " to " + req.Role,
		EntityType:  "users",
		EntityID:    userID,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "user role updated successfully", nil)
}

func (h *Handler) ToggleUserActiveHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetParamID(w, r)
	if !ok {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid user id")
		return
	}

	if userID == 1 {
		utils.ErrorJson(w, http.StatusForbidden, "cannot modify root user")
		return
	}

	user, err := h.app.Models.User.GetByID(r.Context(), userID)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}

	if user == nil {
		utils.ErrorJson(w, http.StatusNotFound, "user not found")
		return
	}

	if user.Role == "admin" {
		utils.ErrorJson(w, http.StatusForbidden, "cannot modify admin users")
		return
	}

	newStatus := !user.IsActive

	err = h.app.Models.User.SetActiveStatus(r.Context(), userID, newStatus)
	if err != nil {
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update user status")
		return
	}

	user.IsActive = newStatus

	// Audit log
	currentUserID, _ := r.Context().Value(middlewares.UserIDKey).(int64)
	ip := r.RemoteAddr
	ua := r.UserAgent()
	action := "deactivated"
	if newStatus {
		action = "activated"
	}
	_, _ = h.app.Models.Log.Insert(context.Background(), &models.Log{
		UserID:      currentUserID,
		Action:      "update",
		Description: "User " + action + ": #" + strconv.FormatInt(userID, 10),
		EntityType:  "users",
		EntityID:    userID,
		IPAddress:   &ip,
		UserAgent:   &ua,
	})

	utils.SuccessJson(w, http.StatusOK, "user status updated successfully", map[string]any{
		"user": user,
	})
}
