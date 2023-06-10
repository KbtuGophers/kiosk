package category

type Entity struct {
	ID   string  `db:"id"`
	Name *string `db:"name"`
}
