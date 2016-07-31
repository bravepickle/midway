// DB-related business logic
package stubman

import (
	"database/sql"
	"errors"
	"fmt"
	//	"log"
	"io/ioutil"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const sqlSchemaInit = `
create table stub (id integer not null primary key AUTOINCREMENT, name text, request_method text, request_uri text, request text, response text, created datetime);
`

// main db connection for app
var DefaultDb *Db

type Db struct {
	DbName     string
	Connection *sql.DB
}

// MakeDefault sets current DB as default
func (d *Db) MakeDefault() {
	DefaultDb = d // just specify as default DB
}

func (d *Db) Close() {
	d.Connection.Close()
}

func (d *Db) Init() error {
	if d.DbName == `` {
		return errors.New(`DB name is no set`)
	}

	conn, err := sql.Open("sqlite3", d.DbName)
	if err != nil {
		return err
	}

	d.Connection = conn

	return nil
}

// Reset resets database to its prestine format
func (d *Db) Reset() error {
	if d.DbName == `` {
		return errors.New(`DB name is no set`)
	}

	if d.Connection != nil {
		d.Close() // close current connection
	}

	// remove current DB
	os.Remove(d.DbName)

	// reconnect
	err := d.Init()
	//	defer db.Close()
	if err != nil {
		return errors.New(fmt.Sprintf("%q: %s\n", err, sqlSchemaInit))
	}

	_, err = d.Connection.Exec(sqlSchemaInit)
	if err != nil {
		return errors.New(fmt.Sprintf("%q: %s\n", err, sqlSchemaInit))
	}

	return nil
}

// ImportFromFile imports SQL from file
func (d *Db) ImportFromFile(path string) error {
	if d.Connection == nil {
		return errors.New(`DB connection is not set`)
	}

	if path == `` {
		return errors.New(`Import file is not specified`)
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	sql, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	_, err = d.Connection.Exec(string(sql))
	if err != nil {
		return err
	}

	return nil
}

func NewDb(dbname string, setAsDefault bool) *Db {
	db := &Db{DbName: dbname}
	db.Init()

	if setAsDefault {
		db.MakeDefault()
	}

	return db
}
