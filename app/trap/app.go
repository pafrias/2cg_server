package trap

import (
	"context"
	"time"

	"github.com/pafrias/2cgaming-api/db"
)

/*Service represents the collection of route handlers and database logic that pertains solely
to the Trap Compendium property*/
type Service struct {
	*db.Connection
}

// NewHandler returns a Service exposing the Trap Compendium API
func NewHandler(db *db.Connection) Service {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	if err := db.PingContext(cxt); err != nil {
		//handle it
	}
	cancel()

	return Service{db}
}
