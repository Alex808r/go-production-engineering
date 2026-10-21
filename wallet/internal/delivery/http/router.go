package http

import (
	"net/http"
)

// NewRouter registers all HTTP routes and wraps them with logging and recovery middlewares.
func NewRouter(h *WalletHandler) http.Handler {
	mux := http.NewServeMux()

	// POST /api/v1/wallet (and /api/v1/wallet/)
	mux.HandleFunc("POST /api/v1/wallet", h.HandleOperation)
	mux.HandleFunc("POST /api/v1/wallet/", h.HandleOperation)

	// GET /api/v1/wallets/{id} (and /api/v1/wallet/{id})
	mux.HandleFunc("GET /api/v1/wallets/{id}", h.HandleGetBalance)
	mux.HandleFunc("GET /api/v1/wallets/{id}/", h.HandleGetBalance)
	mux.HandleFunc("GET /api/v1/wallet/{id}", h.HandleGetBalance)
	mux.HandleFunc("GET /api/v1/wallet/{id}/", h.HandleGetBalance)

	// POST /api/v1/wallets (creation helper)
	mux.HandleFunc("POST /api/v1/wallets", h.HandleCreateWallet)

	// Health check for Docker / Kubernetes
	mux.HandleFunc("GET /health", h.HandleHealth)
	mux.HandleFunc("GET /api/v1/health", h.HandleHealth)

	// Wrap with middlewares
	handler := RecoveryMiddleware(mux)
	handler = LoggingMiddleware(handler)

	return handler
}
