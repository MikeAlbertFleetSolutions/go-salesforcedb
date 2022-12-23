# go-salesforcedb
Golang Salesforce database driver conforming to the Go database/sql interface.

## Usage
This is an implementation of Go's database/sql/driver interface. In order to use it, you need to import the package and use the database/sql API.

Only SELECT operations are supported.

```go
import (
	"database/sql"
	"log"

	_ "github.com/MikeAlbertFleetSolutions/go-salesforcedb"
)

func main() {
	conn, err := sql.Open("salesforcedb", "testuser/password1/XkjHhusah@https://company.my.salesforce.com")
	if err != nil {
		log.Fatalf("%+v", err)
	}

	rows, err := conn.Query("select Name from Account")
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer rows.Close()

	...

}
```