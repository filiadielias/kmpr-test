// Package news contains business logic from news, store to database, etc
package news

import (
	//"log"
	"fmt"
	"time"

	"github.com/filiadielias/kmpr-test/src/general"

	"github.com/filiadielias/kmpr-test/src/helper/db"

	// PostgreSQL driver
	_ "github.com/lib/pq"
)

var kmpr *general.Module

func init() {
	kmpr = &general.KMPR
}

// Error list
var (
	ErrNotFound = fmt.Errorf("Data not Found")
)

// Newses is collection of News
type Newses []News

// News represents news information
type News struct {
	ID      int       `json:"id" db:"id"`
	Author  string    `json:"author" db:"author"`
	Body    string    `json:"body" db:"body"`
	Created time.Time `json:"created" db:"created"`
}

// Insert news
func (n *News) insert() error {
	tx := kmpr.DB.MustBegin()
	tx.QueryRowx("INSERT INTO news(author,body) values($1,$2) returning id,created", n.Author, n.Body).Scan(&n.ID, &n.Created)

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (ns *Newses) getFromDB(qb *db.QueryBuilder) error {

	qb.Query = "select id,author,body,created from news"

	query, params := qb.GetQuery()

	rows, err := kmpr.DB.Queryx(query, params...)
	if err != nil {
		return err
	}

	for rows.Next() {
		var n News

		if err := rows.StructScan(&n); err != nil {
			return err
		}

		*ns = append(*ns, n)
	}

	if len(*ns) == 0 {
		return ErrNotFound
	}

	return nil
}

func (n *News) getFromDB(qb *db.QueryBuilder) error {

	// Call Newses get function
	var ns Newses
	if err := ns.getFromDB(qb); err != nil {
		return err
	}

	if len(ns) == 0 {
		return ErrNotFound
	}

	*n = ns[0]

	return nil
}
