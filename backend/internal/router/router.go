package router

import (
	"backend/internal/app"
	"backend/internal/handlers"
	"backend/internal/middlewares"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func RegisterRouter(app *app.Application, handler *handlers.Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowOriginFunc: func(r *http.Request, origin string) bool {
			return true
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	limiter := middlewares.NewRateLimiter(app.Config.RateLimit, app.Config.RateInterval)
	r.Use(limiter.LimitRate)

	fs := http.FileServer(http.Dir("./public"))
	r.Handle("/public/*", http.StripPrefix("/public/", fs))

	// Helper functions
	protected := func(h http.HandlerFunc) http.HandlerFunc {
		return middlewares.Authenticate(app.Config.JWTKey, h)
	}

	adminOnly := func(h http.HandlerFunc) http.HandlerFunc {
		return protected(
			middlewares.Authorize(app, "admin")(h),
		)
	}

	higherManagementOnly := func(h http.HandlerFunc) http.HandlerFunc {
		return protected(
			middlewares.Authorize(app, "admin", "manager")(h),
		)
	}

	managementOnly := func(h http.HandlerFunc) http.HandlerFunc {
		return protected(
			middlewares.Authorize(app, "admin", "manager", "accounts")(h),
		)
	}

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", handler.HealthCheckHandler)

		r.Route("/users", func(r chi.Router) {
			// Public routes
			r.Post("/login", handler.LoginHandler)
			r.Post("/refresh-session", handler.RefreshHandler)

			// Protected routes
			r.Delete("/logout", protected(handler.LogoutHandler))
			r.Get("/current-user", protected(handler.GetCurrentUserHandler))
			r.Patch("/change-password", protected(handler.ChangePasswordHandler))
			// r.Patch("/profile", protected(handler.UpdateProfileHandler))
			r.Patch("/{id}/profile", higherManagementOnly(handler.UpdateUserProfileHandler))

			// Management routes
			r.Get("/", higherManagementOnly(handler.GetAllUsersHandler))
			r.Post("/", higherManagementOnly(handler.CreateUserHandler))

			// Admin only routes
			r.Patch("/reset-all-passwords", adminOnly(handler.ResetAllPasswordsHandler))

			// Parameterized routes (must come after specific routes)
			r.Get("/{id}", protected(handler.GetUserHandler))
			r.Patch("/{id}/change-password", higherManagementOnly(handler.ChangeUserPasswordHandler))
			r.Patch("/{id}/change-role", higherManagementOnly(handler.ChangeUserRoleHandler))
			r.Patch("/{id}/toggle-active", higherManagementOnly(handler.ToggleUserActiveHandler))
		})

		r.Route("/customers", func(r chi.Router) {
			// Management routes
			r.Get("/", higherManagementOnly(handler.GetAllCustomersHandler))
			r.Post("/", higherManagementOnly(handler.CreateCustomerHandler))

			// Parameterized routes
			r.Get("/{id}", protected(handler.GetCustomerHandler))
			r.Patch("/{id}", higherManagementOnly(handler.UpdateCustomerHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteCustomerHandler))
		})

		r.Route("/contracts", func(r chi.Router) {
			// Management routes
			r.Get("/", higherManagementOnly(handler.GetAllContractsHandler))
			r.Post("/", higherManagementOnly(handler.CreateContractHandler))

			// Parameterized routes (must come after specific routes)
			r.Get("/{id}", higherManagementOnly(handler.GetContractHandler))
			r.Patch("/{id}", higherManagementOnly(handler.UpdateContractHandler))
			r.Patch("/{id}/status", higherManagementOnly(handler.ChangeContractStatusHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteContractHandler))
		})

		r.Route("/items", func(r chi.Router) {
			r.Get("/", protected(handler.GetAllItemsHandler))
			r.Post("/", protected(handler.CreateItemHandler))

			r.Get("/{id}", protected(handler.GetItemHandler))
			r.Patch("/{id}", higherManagementOnly(handler.UpdateItemHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteItemHandler))
		})

		r.Route("/shipments", func(r chi.Router) {
			r.Get("/", higherManagementOnly(handler.GetAllShipmentsHandler))
			r.Post("/", higherManagementOnly(handler.CreateShipmentHandler))

			r.Get("/{id}", higherManagementOnly(handler.GetShipmentHandler))
			r.Patch("/{id}", higherManagementOnly(handler.UpdateShipmentHandler))
			r.Patch("/{id}/status", higherManagementOnly(handler.ChangeShipmentStatusHandler))
			r.Delete("/{id}", higherManagementOnly(handler.DeleteShipmentHandler))
		})

		r.Route("/expenses", func(r chi.Router) {
			r.Get("/", managementOnly(handler.GetAllExpensesHandler))
			r.Post("/", managementOnly(handler.CreateExpenseHandler))

			r.Get("/{id}", managementOnly(handler.GetExpenseHandler))
			r.Patch("/{id}", managementOnly(handler.UpdateExpenseHandler))
			r.Delete("/{id}", managementOnly(handler.DeleteExpenseHandler))
		})

		r.Route("/logs", func(r chi.Router) {
			r.Get("/", higherManagementOnly(handler.GetAllLogsHandler))
		})

		r.Route("/images", func(r chi.Router) {
			r.Post("/{type}/{id}", protected(handler.UploadImageHandler))
			r.Delete("/{type}/{id}", protected(handler.DeleteImageHandler))
		})

		r.Get("/backup", adminOnly(handler.BackupHandler))
	})

	return r
}
