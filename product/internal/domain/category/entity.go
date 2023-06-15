package category

type Entity struct {
	ID       string   `db:"id"`
	ParentId *string  `db:"parent_id"`
	Name     *string  `db:"name"`
	Child    []Entity `db:"child"`
}
