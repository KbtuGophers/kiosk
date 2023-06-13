package handler

import (
	"github.com/KbtuGophers/kiosk/account/internal/handler/http"
	"github.com/KbtuGophers/kiosk/account/internal/service/account"
	"github.com/KbtuGophers/kiosk/account/internal/service/otp"
	"github.com/KbtuGophers/kiosk/account/pkg/server/router"
	"github.com/go-chi/chi/v5"
)

type Configuration func(h *Handler) error

type Dependencies struct {
	AccountService *account.Service
	OtpService     *otp.Service
}

type Handler struct {
	//Service *service.Service
	dependencies Dependencies
	HTTP         *chi.Mux
}

func New(d Dependencies, configs ...Configuration) (h *Handler, err error) {
	// Create the handler
	h = &Handler{
		dependencies: d,
	}

	// Apply all Configurations passed in
	for _, cfg := range configs {
		// Pass the service into the configuration function
		if err = cfg(h); err != nil {
			return
		}
	}

	return
}

func WithHTTPHandler() Configuration {
	return func(h *Handler) (err error) {
		h.HTTP = router.New()
		accountHandler := http.NewAccountHandler(h.dependencies.AccountService)
		otpHandler := http.NewOtpHandler(h.dependencies.OtpService)

		h.HTTP.Route("/api/v1", func(r chi.Router) {
			r.Mount("/user", accountHandler.Routes())
			r.Mount("/otp", otpHandler.Routes())
		})
		return
	}
}
