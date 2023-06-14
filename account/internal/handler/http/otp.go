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
	r.Delete("/", o.DeleteExpiredOtps)

	return r
}

func (o *OtpHandler) GetOtp(w http.ResponseWriter, r *http.Request) {
	req := secret.Request{}
	httpResponse := status.Response{}
	var err error
	req.PhoneNumber = r.URL.Query().Get("phone")
	req.DebugMode, err = strconv.ParseBool(r.URL.Query().Get("debug"))

	if err != nil {
		httpResponse = status.BadRequest(err, req)
		httpResponse.Render(w, r)
		render.JSON(w, r, httpResponse)
		return
	}

	//if err := render.Bind(r, &req); err != nil {
	//	render.JSON(w, r, status.BadRequest(err, req))
	//	return
	//}

	res, err := o.otpService.Create(r.Context(), req)
	if err != nil {
		httpResponse = status.InternalServerError(err)
		httpResponse.Render(w, r)
		render.JSON(w, r, httpResponse)
		return
	}

	httpResponse = status.OK(res)
	httpResponse.Render(w, r)
	render.JSON(w, r, httpResponse)

}

func (o *OtpHandler) CheckOtp(w http.ResponseWriter, r *http.Request) {
	req := secret.Request{}
	httpResponse := status.Response{}
	if err := render.Bind(r, &req); err != nil {
		httpResponse = status.BadRequest(err, req)
		httpResponse.Render(w, r)
		render.JSON(w, r, httpResponse)
		return
	}

	res, err := o.otpService.Check(r.Context(), req)
	if err != nil {

		httpResponse = status.BadRequest(err, req)
		httpResponse.Render(w, r)
		render.JSON(w, r, httpResponse)
		return
	}

	accountInfo, _ := o.otpService.GetAccountByPhone(res.PhoneNumber)
	//fmt.Println(res.PhoneNumber)
	//if err != nil {
	//	httpResponse = status.InternalServerError(err)
	//	httpResponse.Render(w, r)
	//	render.JSON(w, r, httpResponse)
	//	return
	//}
	//
	//fmt.Println(accountInfo)
	//

	//if accountInfo.ID == "" {
	//	httpResponse = status.InternalServerError(errors.New("account id is invalid"))
	//	httpResponse.Render(w, r)
	//	render.JSON(w, r, httpResponse)
	//	return
	//}
	//
	//if err = o.otpService.InsertActivities(accountInfo.ID); err != nil {
	//	httpResponse = status.InternalServerError(err)
	//	httpResponse.Render(w, r)
	//	render.JSON(w, r, httpResponse)
	//	return
	//}

	httpResponse = status.OK(accountInfo)
	httpResponse.Render(w, r)
	render.JSON(w, r, httpResponse)

}

func (o *OtpHandler) DeleteExpiredOtps(w http.ResponseWriter, r *http.Request) {
	httpResponse := status.Response{}

	if err := o.otpService.DeleteExpiredTokens(); err != nil {
		httpResponse = status.InternalServerError(err)
		httpResponse.Render(w, r)
		render.JSON(w, r, httpResponse)
		return
	}

	httpResponse = status.OK("expired otps deleted")
	httpResponse.Render(w, r)
	render.JSON(w, r, httpResponse)
}
