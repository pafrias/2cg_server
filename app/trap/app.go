package trap

import (
	"github.com/pafrias/2cgaming-api/db"
)

/*Service represents the collection of route handlers and database logic that pertains solely
to the Trap Compendium property*/
type Service struct {
	*db.Connection
}

// OpenService returns a Service exposing the Trap Compendium API
func OpenService(db *db.Connection) Service {
	return Service{db}
}
