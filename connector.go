package sqldb

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // mysql
	_ "github.com/lib/pq"              // postgres
)

const (
	DRIVER_POSTGRES = "postgres"
	DRIVER_MYSQL    = "mysql"

	TICK_INTERVAL = 50 * time.Millisecond
)

// Conf structure to open the connection to the database.
type Conf struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

// Connector is the connector used to communicate with PostgreSQL database server.
type Connector struct {
	*sql.DB
}

// TryConnection test one ping every 50ms with timeout of t second, it ends if the ping is a success or timeout.
func (con *Connector) TryConnection(t int) error {
	ticker := time.NewTicker(TICK_INTERVAL)
	defer ticker.Stop()

	timeout := time.After(time.Duration(t) * time.Second)

	for {
		select {
		case <-ticker.C:
			err := con.Ping()
			if err == nil {
				return nil
			}

		case <-timeout:
			return fmt.Errorf("can't ping SqlDB: timeout after %d s", t)
		}
	}
}

// Commit the given transaction.
func (con *Connector) Commit(tx *sql.Tx) error {
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("fail to commit transaction: %w", err)
	}
	return nil
}

// Exec a query with a prepared statement.
func (con *Connector) Exec(tx *sql.Tx, query string, args ...interface{}) (sql.Result, error) {
	var res sql.Result
	stm, err := tx.Prepare(query)
	if err != nil {
		return res, fmt.Errorf("can't prepare query: %w", err)
	}

	res, err = stm.Exec(args...)
	if err != nil {
		return res, fmt.Errorf("error when executing prepared statement: %w", err)
	}

	return res, nil
}

// ExecQueryRow queries the database in a transaction and returns the only row (i.e. for an insert with a RETURNING).
func (con *Connector) ExecQueryRow(tx *sql.Tx, query string, args ...interface{}) (*sql.Row, error) {
	var row *sql.Row
	stm, err := tx.Prepare(query)
	if err != nil {
		return row, fmt.Errorf("can't prepare query: %w", err)
	}

	row = stm.QueryRow(args...)

	return row, nil
}

// FactoryConnector instanciates a new *Connector with the given params.
func FactoryConnector(c Conf) (*Connector, error) {
	var err error
	con := new(Connector)
	con.DB, err = sql.Open(c.Driver, c.DSN)
	if err != nil {
		return con, fmt.Errorf("can't open connection to SqlDB(driver: %s): %w", c.Driver, err)
	}

	return con, nil
}
