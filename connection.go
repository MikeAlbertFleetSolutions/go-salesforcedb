package salesforcedb

import (
	"context"
	"database/sql/driver"
	"fmt"

	"github.com/simpleforce/simpleforce"
)

// Conn is a salesforce connection
type Conn struct {
	client *simpleforce.Client
}

var (
	_ driver.Conn = &Conn{}

	// errors
	errNoPrepared     = fmt.Errorf("no support for prepared statements")
	errNoTransactions = fmt.Errorf("no support for transactions")
)

// Begin not supported but satisfies the interface requirements
func (c *Conn) Begin() (driver.Tx, error) {
	return nil, errNoTransactions
}

// Prepare not supported but satisfies the interface requirements
func (c *Conn) Prepare(query string) (driver.Stmt, error) {
	return nil, errNoPrepared
}

// Close closes the salesforce connection
func (c *Conn) Close() error {
	c.client = nil
	return nil
}

// Ping NOP but satisfies the interface requirements
func (c *Conn) Ping(ctx context.Context) error {
	return nil
}

// Query carries out the simpleforce query
func (c *Conn) Query(query string, args []driver.Value) (driver.Rows, error) {
	if len(args) > 0 {
		return nil, errNoPrepared
	}

	// run query using api
	rs, err := c.client.Query(query)
	if err != nil {
		return nil, err
	}

	rows := &Rows{
		connection: c,
		results:    rs,
	}

	return rows, nil
}
