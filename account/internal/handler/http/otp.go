package http

import (
	"github.com/KbtuGophers/kiosk/account/internal/domain/secret"
	"github.com/KbtuGophers/kiosk/account/internal/service/otp"
	"github.com/KbtuGophers/kiosk/account/pkg/server/status"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

type OtpHandler struct {
	otpService *otp.Service
}

func NewOtpHandler(service *otp.Service) *OtpHandler {
	return &OtpHandler{otpService: service}
}

func (o *OtpHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", o.CheckOtp)
	r.Get("/", o.GetOtp)

	return r
}

func (o *OtpHandler) GetOtp(w http.ResponseWriter, r *http.Request) {
	req := secret.Request{}
	var err error
	req.PhoneNumber = r.URL.Query().Get("phone")
	req.DebugMode, err = strconv.ParseBool(r.URL.Query().Get("debug"))

	if err != nil {
		render.JSON(w, r, status.BadRequest(err, req))
		return
	}

	if err := render.Bind(r, &req); err != nil {
		render.JSON(w, r, status.BadRequest(err, req))
		return
	}

	res, err := o.otpService.Create(r.Context(), req)
	if err != nil {
		render.JSON(w, r, status.InternalServerError(err))
		return
	}

	render.JSON(w, r, status.OK(res))

}

func (o *OtpHandler) CheckOtp(w http.ResponseWriter, r *http.Request) {
	req := secret.Request{}

	if err := render.Bind(r, &req); err != nil {
		render.JSON(w, r, status.BadRequest(err, req))
		return
	}

	res, err := o.otpService.GetById(r.Context(), req)
	if err != nil {
		render.JSON(w, r, status.InternalServerError(err))
		return
	}

	render.JSON(w, r, status.OK(res))

}
