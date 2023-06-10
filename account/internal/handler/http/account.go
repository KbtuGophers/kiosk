package http

import (
	"github.com/KbtuGophers/kiosk/account/internal/domain/user"
	"github.com/KbtuGophers/kiosk/account/internal/service/account"
	"github.com/KbtuGophers/kiosk/account/pkg/server/status"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
)

type AccountHandler struct {
	accountService *account.Service
}

func NewAccountHandler(service *account.Service) *AccountHandler {
	return &AccountHandler{accountService: service}
}

func (a *AccountHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", a.add)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", a.get)
		r.Put("/", a.update)
		r.Delete("/", a.delete)
	})

	return r
}

func (a *AccountHandler) add(w http.ResponseWriter, r *http.Request) {
	req := user.Request{}
	if err := render.Bind(r, &req); err != nil {
		render.JSON(w, r, status.BadRequest(err, req))
		return
	}

	res, err := a.accountService.AddAccount(r.Context(), req)
	if err != nil {
		render.JSON(w, r, status.InternalServerError(err))
		return
	}

	render.JSON(w, r, status.OK(res))
}

func (a *AccountHandler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	res, err := a.accountService.GetAuthor(r.Context(), id)
	if err != nil {
		render.JSON(w, r, status.InternalServerError(err))
		return
	}
	render.JSON(w, r, status.OK(res))
}
func (a *AccountHandler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	req := user.Request{}
	if err := render.Bind(r, &req); err != nil {
		render.JSON(w, r, status.BadRequest(err, req))
		return
	}

	err := a.accountService.UpdateAccount(r.Context(), id, req)
	if err != nil {
		render.JSON(w, r, status.InternalServerError(err))
		return
	}

}
func (a *AccountHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := a.accountService.DeleteAuthor(r.Context(), id); err != nil {
		render.JSON(w, r, status.InternalServerError(err))
		return
	}
}
