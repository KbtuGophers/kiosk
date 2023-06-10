package product

type Entity struct {
	ID          string  `db:"id"`
	Name        *string `db:"name"`
	Description *string `db:"genre"`
	Cost        *int    `db:"isbn"`
	//Images      postgres.Array `db:"authors"`
	//Characteristics ???
	Category *string `db:"category"`
}
