package product

import (
	"errors"
	"net/http"
)

type Request struct {
	ID          string `json:"id"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Cost        int    `json:"cost" validate:"required"`
	Category    string `json:"category" validate:"required"`
}

func (s *Request) Bind(r *http.Request) error {
	if s.Name == "" {
		return errors.New("name: cannot be blank")
	}

	if s.Description == "" {
		return errors.New("discription: cannot be blank")
	}

	if s.Cost == 0 {
		return errors.New("cost: cannot be blank")
	}

	if s.Category == "" {
		return errors.New("category: cannot be blank")
	}

	return nil
}

type Response struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Cost        int    `json:"cost"`
	Category    string `json:"category"`
}

func ParseFromEntity(data Entity) (res Response) {
	res = Response{
		ID:          data.ID,
		Name:        *data.Name,
		Description: *data.Description,
		Cost:        *data.Cost,
		Category:    *data.Category,
	}
	return
}

func ParseFromEntities(data []Entity) (res []Response) {
	res = make([]Response, 0)
	for _, object := range data {
		res = append(res, ParseFromEntity(object))
	}
	return
}
