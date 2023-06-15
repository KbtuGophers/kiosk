package product

import (
	"errors"
	"net/http"
)

type Request struct {
	ID              string `json:"id"`
	CategoryID      string `json:"category_id"`
	Barcode         string `json:"barcode"`
	Name            string `json:"name"`
	Measure         string `json:"measure"`
	Cost            int    `json:"cost"`
	ProducerCountry string `json:"producer_country"`
	BrandName       string `json:"brand_name"`
	Description     string `json:"description"`
	Image           string `json:"image"`
	IsWeighted      bool   `json:"is_weighted"`
}

func (s *Request) Bind(r *http.Request) error {
	if s.Name == "" {
		return errors.New("name: cannot be blank")
	}

	if s.CategoryID == "" {
		return errors.New("category_id: cannot be blank")
	}
	return nil
}

type Response struct {
	ID              string `json:"id"`
	CategoryID      string `json:"category_id"`
	Barcode         string `json:"barcode"`
	Name            string `json:"name"`
	Measure         string `json:"measure"`
	Cost            int    `json:"cost"`
	ProducerCountry string `json:"producer_country"`
	BrandName       string `json:"brand_name"`
	Description     string `json:"description"`
	Image           string `json:"image"`
	IsWeighted      bool   `json:"is_weighted"`
}

func ParseFromEntity(data Entity) (res Response) {
	res = Response{
		ID:              data.ID,
		CategoryID:      *data.CategoryID,
		Barcode:         *data.Barcode,
		Name:            *data.Name,
		Measure:         *data.Measure,
		Cost:            *data.Cost,
		ProducerCountry: *data.ProducerCountry,
		BrandName:       *data.BrandName,
		Description:     *data.Description,
		Image:           *data.Image,
		IsWeighted:      *data.IsWeighted,
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
