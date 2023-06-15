package service

import (
	"context"
	"github.com/google/uuid"
	"net/http"
	"product/internal/domain/product"
)

func (s *Service) ListProduct(ctx context.Context, r *http.Request) (res []product.Response, err error) {
	data, err := s.productRepository.Select(ctx, r)
	if err != nil {
		return
	}
	res = product.ParseFromEntities(data)

	return
}

func (s *Service) AddProduct(ctx context.Context, req product.Request) (res product.Response, err error) {
	data := product.Entity{
		ID:              uuid.New().String(),
		CategoryID:      &req.CategoryID,
		Barcode:         &req.Barcode,
		Name:            &req.Name,
		Measure:         &req.Measure,
		Cost:            &req.Cost,
		ProducerCountry: &req.ProducerCountry,
		BrandName:       &req.BrandName,
		Description:     &req.Description,
		Image:           &req.Image,
		IsWeighted:      &req.IsWeighted,
	}

	data.ID, err = s.productRepository.Create(ctx, data)
	if err != nil {
		return
	}
	res = product.ParseFromEntity(data)

	return
}

func (s *Service) GetProduct(ctx context.Context, id string) (res product.Response, err error) {
	data, err := s.productRepository.Get(ctx, id)
	if err != nil {
		return
	}
	res = product.ParseFromEntity(data)

	return
}

func (s *Service) UpdateProduct(ctx context.Context, id string, req product.Request) (res product.Response, err error) {
	data := product.Entity{
		ID:              id,
		CategoryID:      &req.CategoryID,
		Barcode:         &req.Barcode,
		Name:            &req.Name,
		Measure:         &req.Measure,
		Cost:            &req.Cost,
		ProducerCountry: &req.ProducerCountry,
		BrandName:       &req.BrandName,
		Description:     &req.Description,
		Image:           &req.Image,
		IsWeighted:      &req.IsWeighted,
	}

	err = s.productRepository.Update(ctx, id, data)

	if err != nil {
		return
	}
	res = product.ParseFromEntity(data)

	return
}

func (s *Service) DeleteProduct(ctx context.Context, id string) (err error) {
	return s.productRepository.Delete(ctx, id)
}
