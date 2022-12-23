package salesforcedb

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"regexp"

	"github.com/simpleforce/simpleforce"
)

// salesforcedriver is a sql.Driver
type salesforcedriver struct{}

var (
	_ driver.Driver = &salesforcedriver{}

	// regexs
	rConnString = regexp.MustCompile(`^([^/]+)/([^/@]+)/?([^@]*)@(.+)`)

	// errors
	errBadConnString = fmt.Errorf("malformed connection string")
	errUnknown       = fmt.Errorf("unknown error creating client")
)

func init() {
	sql.Register("salesforcedb", &salesforcedriver{})
}

// Open prepares the connection to the salesforce endpoint
// connString should be of the format: username/password/token@endpoint
// token is optional if trusted IP is configured
func (d salesforcedriver) Open(connString string) (driver.Conn, error) {
	//	extract properties to use for connection
	matches := rConnString.FindStringSubmatch(connString)
	if len(matches) != 5 {
		return nil, errBadConnString
	}
	username := matches[1]
	password := matches[2]
	token := matches[3]
	endpoint := matches[4]

	// create connection to salesforce api
	client := simpleforce.NewClient(endpoint, simpleforce.DefaultClientID, simpleforce.DefaultAPIVersion)
	if client == nil {
		return nil, errUnknown
	}
	err := client.LoginPassword(username, password, token)
	if err != nil {
		return nil, err
	}

	c := &Conn{
		client: client,
	}

	return c, nil
}
