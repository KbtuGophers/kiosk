package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"product/internal/domain/category"
	"product/internal/service"
	"product/pkg/server/status"
	"product/pkg/store"
)

type CategoryHandler struct {
	Service *service.Service
}

func NewCategoryHandler(s *service.Service) *CategoryHandler {
	return &CategoryHandler{Service: s}
}

func (h *CategoryHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.list)
	r.Post("/", h.add)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.get)
		r.Put("/", h.update)
		r.Delete("/", h.delete)
	})

	return r
}

// List of categories from the database
//
//	@Summary	List of categories from the database
//	@Tags		categories
//	@Accept		json
//	@Produce	json
//	@Success	200				{array}		author.Response
//	@Failure	500				{object}	status.Response
//	@Router		/categories 	[get]
func (h *CategoryHandler) list(w http.ResponseWriter, r *http.Request) {
	res, err := h.Service.ListCategories(r.Context())
	if err != nil {
		render.JSON(w, r, status.InternalServerError(err))
		return
	}

	render.JSON(w, r, status.OK(res))
}

// Add a new category to the database
//
//	@Summary	Add a new author to the database
//	@Tags		categories
//	@Accept		json
//	@Produce	json
//	@Param		request	body		category.Request	true	"body param"
//	@Success	200		{object}	category.Response
//	@Failure	400		{object}	status.Response
//	@Failure	500		{object}	status.Response
//	@Router		/categories [post]
func (h *CategoryHandler) add(w http.ResponseWriter, r *http.Request) {
	req := category.Request{}
	if err := render.Bind(r, &req); err != nil {
		render.JSON(w, r, status.BadRequest(err, req))
		return
	}

	res, err := h.Service.AddCategory(r.Context(), req)
	if err != nil {
		render.JSON(w, r, status.InternalServerError(err))
		return
	}

	render.JSON(w, r, status.OK(res))
}

// Read the category from the database
//
//	@Summary	Read the category from the database
//	@Tags		categories
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"path param"
//	@Success	200	{object}	category.Response
//	@Failure	404	{object}	status.Response
//	@Failure	500	{object}	status.Response
//	@Router		/categories/{id} [get]
func (h *CategoryHandler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	res, err := h.Service.GetCategory(r.Context(), id)
	if err != nil && err != store.ErrorNotFound {
		render.JSON(w, r, status.InternalServerError(err))
		return
	}

	if err == store.ErrorNotFound {
		render.JSON(w, r, status.NotFound(err))
		return
	}

	render.JSON(w, r, status.OK(res))
}

// Update the category in the database
//
//	@Summary	Update the category in the database
//	@Tags		categories
//	@Accept		json
//	@Produce	json
//	@Param		id		path	int					true	"path param"
//	@Param		request	body	category.Request	true	"body param"
//	@Success	200
//	@Failure	400	{object}	status.Response
//	@Failure	404	{object}	status.Response
//	@Failure	500	{object}	status.Response
//	@Router		/categories/{id} [put]
func (h *CategoryHandler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	req := category.Request{}
	if err := render.Bind(r, &req); err != nil {
		render.JSON(w, r, status.BadRequest(err, req))
		return
	}

	err := h.Service.UpdateCategory(r.Context(), id, req)
	if err != nil && err != store.ErrorNotFound {
		render.JSON(w, r, status.InternalServerError(err))
		return
	}

	if err == store.ErrorNotFound {
		render.JSON(w, r, status.NotFound(err))
		return
	}
}

// Delete the category from the database
//
//	@Summary	Delete the category from the database
//	@Tags		categories
//	@Accept		json
//	@Produce	json
//	@Param		id	path	int	true	"path param"
//	@Success	200
//	@Failure	404	{object}	status.Response
//	@Failure	500	{object}	status.Response
//	@Router		/categories/{id} [delete]
func (h *CategoryHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.Service.DeleteCategory(r.Context(), id)
	if err != nil && err != store.ErrorNotFound {
		render.JSON(w, r, status.InternalServerError(err))
		return
	}

	if err == store.ErrorNotFound {
		render.JSON(w, r, status.NotFound(err))
		return
	}
}
