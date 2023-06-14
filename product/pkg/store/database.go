package store

import (
	_ "database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"net/url"
	"strings"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	//_ "github.com/sijms/go-ora/v2"
)

// postgres://username:password@localhost:5432/dbname?sslmode=disable

type Database struct {
	schema         string
	driverName     string
	dataSourceName string

	Client *sqlx.DB
}

func (s *Database) Print() {
	fmt.Println("Schema: ", s.schema)
	fmt.Println("DriverName: ", s.driverName)
	fmt.Println("Data source: ", s.dataSourceName)
	fmt.Println("client: ", s.Client)

}

// NewDatabase established connection to a database instance using provided URI and auth credentials.
func NewDatabase(schema, dataSourceName string) (database *Database, err error) {
	database = &Database{
		schema:         schema,
		dataSourceName: dataSourceName,
	}

	if err = database.connection(); err != nil {
		return
	}
	err = database.createSchema()

	database.Print()

	return
}

func (s *Database) Migrate() (err error) {
	migrations, err := migrate.New("file://migrations", s.dataSourceName)
	if err != nil {
		return
	}

	if err = migrations.Up(); err != nil && err != migrate.ErrNoChange {
		return
	}

	return nil
}

func (s *Database) connection() error {
	err := s.parseDSN()
	if err != nil {
		return err
	}

	s.Client, err = sqlx.Connect(s.driverName, s.dataSourceName)
	if err != nil {
		return err
	}
	s.Client.SetMaxOpenConns(20)

	return nil
}

func (s *Database) parseDSN() (err error) {
	if !strings.Contains(s.dataSourceName, "://") {
		err = errors.New("sql: undefined data source name " + s.dataSourceName)
		return
	}
	s.driverName = strings.ToLower(strings.Split(s.dataSourceName, "://")[0])

	source, err := url.Parse(s.dataSourceName)
	if err != nil {
		return
	}
	sourceQuery := source.Query()

	if s.schema != "" {
		switch s.driverName {
		case "postgres":
			sourceQuery.Set("search_path", s.schema)
		}
	}

	source.RawQuery = sourceQuery.Encode()
	s.dataSourceName = source.String()

	return
}

func (s *Database) createSchema() (err error) {
	if s.schema == "" {
		return
	}

	query := ""
	switch s.driverName {
	case "postgres":
		query = fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", s.schema)
	}

	_, err = s.Client.Exec(query)

	return
}
