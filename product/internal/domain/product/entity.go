package product

type Entity struct {
	ID              string  `db:"id"`
	CategoryID      *string `db:"category_id"`
	Barcode         *string `db:"barcode"`
	Name            *string `db:"name"`
	Measure         *string `db:"measure"`
	Cost            *int    `db:"cost"`
	ProducerCountry *string `db:"producer_country"`
	BrandName       *string `db:"brand_name"`
	Description     *string `db:"description"`
	Image           *string `db:"image"`
	IsWeighted      *bool   `db:"is_weighted"`
}
