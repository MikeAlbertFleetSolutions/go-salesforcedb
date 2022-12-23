package salesforcedb

import (
	"database/sql/driver"
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/simpleforce/simpleforce"
)

// Rows structure to track results
type Rows struct {
	connection *Conn
	columns    []string
	results    *simpleforce.QueryResult
	currentRow int
}

var (
	// errors
	errClosed = fmt.Errorf("rows closed")
)

// stringInSlice returns true if string s is in slice p
func stringInSlice(s string, p []string) bool {
	for _, pp := range p {
		if strings.EqualFold(s, pp) {
			return true
		}
	}
	return false
}

// Columns returns the columns in the result set
func (rows *Rows) Columns() []string {
	if rows.results == nil {
		return nil
	}

	// get column names once so order is always deterministic since the underlying type with simpleforce is a map
	if rows.columns == nil {
		// columns are the map keys
		r := rows.results.Records[0]
		columns := make([]string, 0, len(r))
		for k := range r {
			// filter out keys internal to simpleforce
			if !stringInSlice(k, []string{"__client__", "attributes"}) {
				columns = append(columns, k)
			}
		}

		rows.columns = columns
	}

	return rows.columns
}

// Next navigates to next row in resultset
func (rows *Rows) Next(dest []driver.Value) error {
	if rows.results == nil {
		return errClosed
	}

	if rows.currentRow >= len(rows.results.Records) {
		if rows.currentRow >= rows.results.TotalSize {
			// at end
			return io.EOF
		}
		if rows.results.NextRecordsURL == "" {
			// no next batch to get
			return io.EOF
		}

		// get next batch of records
		rs, err := rows.connection.client.Query(rows.results.NextRecordsURL)
		if err != nil {
			return err
		}

		(*rows).results = rs
		rows.currentRow = 0
	}

	// dest values driven off column names
	if rows.columns == nil {
		rows.columns = rows.Columns()
	}

	// put values from current entry
	values := rows.results.Records[rows.currentRow]
	for i, n := range rows.columns {
		dest[i] = values[n]
	}
	rows.currentRow++

	return nil
}

// Close closes the rows
func (rows *Rows) Close() error {
	if rows.results == nil {
		return errClosed
	}
	rows.connection = nil
	rows.columns = nil
	rows.results = nil
	rows.currentRow = math.MaxInt

	return nil
}
