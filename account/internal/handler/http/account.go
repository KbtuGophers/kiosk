package http

import (
	"database/sql"
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
	r.Get("/", a.list)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", a.get)
		r.Put("/", a.update)
		r.Delete("/", a.delete)
	})
	r.Route("/type/{name}", func(r chi.Router) {
		r.Post("/", a.InsertType)
	})

	return r
}

func (a *AccountHandler) list(w http.ResponseWriter, r *http.Request) {
	res, err := a.accountService.GetAllAccounts(r.Context())
	httpRespoonse := status.Response{}
	if err != nil {
		httpRespoonse = status.InternalServerError(err)
		render.JSON(w, r, httpRespoonse)
		httpRespoonse.Render(w, r)
		return
	}

	httpRespoonse = status.OK(res)
	httpRespoonse.Render(w, r)
	render.JSON(w, r, httpRespoonse)
}

func (a *AccountHandler) add(w http.ResponseWriter, r *http.Request) {
	req := user.Request{}
	httpResponse := status.Response{}
	if err := render.Bind(r, &req); err != nil {
		httpResponse = status.BadRequest(err, req)
		httpResponse.Render(w, r)
		render.JSON(w, r, httpResponse)
		return
	}

	res, err := a.accountService.AddAccount(r.Context(), req)
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

func (a *AccountHandler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	httpResponse := status.Response{}

	res, err := a.accountService.GetAccount(r.Context(), id)
	if err != nil && err == sql.ErrNoRows {
		//render.JSON(w, r, status.BadRequest(err, req))
		httpResponse = status.NotFoundError(err)
		httpResponse.Render(w, r)
		render.JSON(w, r, httpResponse)
		return
	} else if err != nil {
		httpResponse = status.InternalServerError(err)
		httpResponse.Render(w, r)
		render.JSON(w, r, httpResponse)
		return
	}
	httpResponse = status.OK(res)
	httpResponse.Render(w, r)
	render.JSON(w, r, httpResponse)
}

func (a *AccountHandler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	httpResponse := status.Response{}
	req := user.Request{}
	if err := render.Bind(r, &req); err != nil {
		httpResponse = status.BadRequest(err, req)
		httpResponse.Render(w, r)
		render.JSON(w, r, httpResponse)
		return
	}

	err := a.accountService.UpdateAccount(r.Context(), id, req)
	if err != nil {
		httpResponse = status.InternalServerError(err)
		httpResponse.Render(w, r)
		render.JSON(w, r, httpResponse)
		return
	}

	httpResponse = status.OK("updated")
	httpResponse.Render(w, r)
	render.JSON(w, r, httpResponse)
}
func (a *AccountHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	httpResponse := status.Response{}

	if err := a.accountService.DeleteAccount(r.Context(), id); err != nil {
		httpResponse = status.InternalServerError(err)
		httpResponse.Render(w, r)
		render.JSON(w, r, httpResponse)
		return
	}

	httpResponse = status.OK("deleted")
	httpResponse.Render(w, r)
	render.JSON(w, r, httpResponse)
}

func (a *AccountHandler) InsertType(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	httpResponse := status.Response{}
	id, err := a.accountService.CreateAccountType(user.Types{Name: name})
	if err != nil {
		httpResponse = status.InternalServerError(err)
		httpResponse.Render(w, r)
		render.JSON(w, r, httpResponse)
		return
	}
	httpResponse = status.OK(map[string]interface{}{
		"type_id": id,
	})
	httpResponse.Render(w, r)

	render.JSON(w, r, httpResponse)

	return
}
