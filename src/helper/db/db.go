// Package db contains database helper function such as connect, query builder, etc
package db

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Conn is db connection builder, stores db connect information
type Conn struct {
	hostname string
	port     int
	user     string
	password string
	database string
}

// InitDB : Initialize DB Connection Builder
func InitDB(hostname string, port int) *Conn {
	var d Conn
	d.hostname = hostname
	d.port = port
	return &d
}

// Credential : Add credential
func (d *Conn) Credential(user, password string) *Conn {
	d.user = user
	d.password = password

	return d
}

// Database : Add Database name
func (d *Conn) Database(dbname string) *Conn {
	d.database = dbname

	return d
}

// Connect using db connection builder information
// Usage : InitDB( ... , ...)
//			.Credential( ... , ... )
//			.Database( ... )
//			.Connect()
func (d *Conn) Connect() (db *sqlx.DB, err error) {
	switch {
	case len(d.hostname) == 0:
		err = fmt.Errorf("invalid hostname")
	case d.port <= 0:
		// Default port
		d.port = 5432
	case len(d.user) == 0:
		// Default user
		d.user = "postgres"
	case len(d.password) == 0:
		err = fmt.Errorf("invalid password")
	case len(d.database) == 0:
		err = fmt.Errorf("invalid database")
	}
	if err != nil {
		return db, err
	}

	connString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		d.hostname, d.port, d.user, d.password, d.database)

	db, err = sqlx.Connect("postgres", connString)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// QueryBuilder is simple query builder to sorting, filtering and pagination
type QueryBuilder struct {
	Query   string
	filters []struct {
		Column   string
		Value    interface{}
		Operator string
	}
	Page  int
	Limit int
	Sort  map[string]string
}

// AddFilter used for adding filter condition.
func (qb *QueryBuilder) AddFilter(column string, value interface{}, operator string) error {
	//validate empty column
	if len(column) == 0 {
		return fmt.Errorf("Column cannot be empty")
	}

	//validate operator
	if !isValidOperator(operator) {
		operator = "="
	}

	data := struct {
		Column   string
		Value    interface{}
		Operator string
	}{
		column,
		value,
		operator,
	}

	qb.filters = append(qb.filters, data)

	return nil
}

// GetQuery : construct query
func (qb *QueryBuilder) GetQuery() (string, []interface{}) {
	var query strings.Builder

	query.WriteString(qb.Query)
	query.WriteString(" WHERE 1 = 1 ")

	var filterValues []interface{}

	if len(qb.filters) > 0 {
		count := 1
		for _, flt := range qb.filters {
			query.WriteString(" AND ")
			switch {
			case flt.Operator == "is not null" || flt.Operator == "is null":
				query.WriteString(fmt.Sprintf(" %s %s ", flt.Column, flt.Operator))
				continue //doesn't add param value
			case flt.Operator == "in" || flt.Operator == "not in":
				query.WriteString(fmt.Sprintf(" %s %s ($%d) ", flt.Column, flt.Operator, count))
			default:
				query.WriteString(fmt.Sprintf(" %s %s $%d ", flt.Column, flt.Operator, count))
			}

			filterValues = append(filterValues, flt.Value)
			count++
		}
	}

	if len(qb.Sort) > 0 {
		query.WriteString(" ORDER BY ")
		for key, value := range qb.Sort {
			query.WriteString(fmt.Sprintf(" %s %s ", key, value))
		}
	}

	if qb.Limit > 0 {
		query.WriteString(fmt.Sprintf(" LIMIT %d ", qb.Limit))

		if qb.Page > 0 {
			offset := (qb.Page - 1) * qb.Limit //Page 1 means item 0..Limit

			query.WriteString(fmt.Sprintf(" OFFSET %d ", offset))
		}
	}

	return query.String(), filterValues
}

// Check filter operator
func isValidOperator(operator string) bool {

	oprs := []string{"=", "!=", "<", ">", "<=", ">=", "in", "not in", "like", "not like", "is null", "is not null"}
	for _, o := range oprs {
		if operator == o {
			return true
		}
	}
	return false
}
