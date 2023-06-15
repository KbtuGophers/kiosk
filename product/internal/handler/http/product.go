package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"product/internal/domain/product"
	"product/internal/service"
	"product/pkg/server/status"
	"product/pkg/store"
)

type ProductHandler struct {
	Service *service.Service
}

func NewProductHandler(s *service.Service) *ProductHandler {
	return &ProductHandler{Service: s}
}

func (h *ProductHandler) Routes() chi.Router {
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

// List of products from the database
//
//	@Summary	List of products from the database
//	@Tags		products
//	@Accept		json
//	@Produce	json
//	@Success	200			{array}		product.Response
//	@Failure	500			{object}	status.Response
//	@Router		/products 	[get]
func (h *ProductHandler) list(w http.ResponseWriter, r *http.Request) {
	res, err := h.Service.ListProduct(r.Context(), r)
	if err != nil {
		render.JSON(w, r, status.InternalServerError(err))
		return
	}

	render.JSON(w, r, status.OK(res))
}

// Add a new product to the database
//
//	@Summary	Add a new product to the database
//	@Tags		products
//	@Accept		json
//	@Produce	json
//	@Param		request	body		product.Request	true	"body param"
//	@Success	200		{object}	product.Response
//	@Failure	400		{object}	status.Response
//	@Failure	500		{object}	status.Response
//	@Router		/products [post]
func (h *ProductHandler) add(w http.ResponseWriter, r *http.Request) {
	req := product.Request{}
	if err := render.Bind(r, &req); err != nil {
		render.JSON(w, r, status.BadRequest(err, req))
		return
	}

	res, err := h.Service.AddProduct(r.Context(), req)
	if err != nil {
		render.JSON(w, r, status.InternalServerError(err))
		return
	}

	render.JSON(w, r, status.OK(res))
}

// Read the product from the database
//
//	@Summary	Read the product from the database
//	@Tags		products
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"path param"
//	@Success	200	{object}	product.Response
//	@Failure	404	{object}	status.Response
//	@Failure	500	{object}	status.Response
//	@Router		/products/{id} [get]
func (h *ProductHandler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	res, err := h.Service.GetProduct(r.Context(), id)
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

// Update the product in the database
//
//	@Summary	Update the product in the database
//	@Tags		products
//	@Accept		json
//	@Produce	json
//	@Param		id		path	int				true	"path param"
//	@Param		request	body	product.Request	true	"body param"
//	@Success	200
//	@Failure	400	{object}	status.Response
//	@Failure	404	{object}	status.Response
//	@Failure	500	{object}	status.Response
//	@Router		/products/{id} [put]
func (h *ProductHandler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	req := product.Request{}
	if err := render.Bind(r, &req); err != nil {
		render.JSON(w, r, status.BadRequest(err, req))
		return
	}

	res, err := h.Service.UpdateProduct(r.Context(), id, req)
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

// Delete the product from the database
//
//	@Summary	Delete the product from the database
//	@Tags		products
//	@Accept		json
//	@Produce	json
//	@Param		id	path	int	true	"path param"
//	@Success	200
//	@Failure	404	{object}	status.Response
//	@Failure	500	{object}	status.Response
//	@Router		/products/{id} [delete]
func (h *ProductHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.Service.DeleteProduct(r.Context(), id)
	if err != nil && err != store.ErrorNotFound {
		render.JSON(w, r, status.InternalServerError(err))
		return
	}

	if err == store.ErrorNotFound {
		render.JSON(w, r, status.NotFound(err))
		return
	}
}
