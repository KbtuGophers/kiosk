package handler

import (
	"github.com/KbtuGophers/kiosk/account/docs"
	"github.com/KbtuGophers/kiosk/account/internal/config"
	"github.com/KbtuGophers/kiosk/account/internal/handler/http"
	"github.com/KbtuGophers/kiosk/account/internal/service/account"
	"github.com/KbtuGophers/kiosk/account/internal/service/otp"
	"github.com/KbtuGophers/kiosk/account/pkg/server/router"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"net/url"
)

type Configuration func(h *Handler) error

type Dependencies struct {
	Configs        config.Config
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

		docs.SwaggerInfo.BasePath = "/api/v1"
		docs.SwaggerInfo.Host = h.dependencies.Configs.HTTP.Host
		docs.SwaggerInfo.Schemes = []string{h.dependencies.Configs.HTTP.Schema}

		swaggerURL := url.URL{
			Scheme: h.dependencies.Configs.HTTP.Schema,
			Host:   h.dependencies.Configs.HTTP.Host,
			Path:   "swagger/doc.json",
		}

		h.HTTP.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(swaggerURL.String()),
		))

		accountHandler := http.NewAccountHandler(h.dependencies.AccountService)
		otpHandler := http.NewOtpHandler(h.dependencies.OtpService)

		h.HTTP.Route("/api/v1", func(r chi.Router) {
			r.Mount("/user", accountHandler.Routes())
			r.Mount("/otp", otpHandler.Routes())
		})
		return
	}
}
