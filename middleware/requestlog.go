package middleware

import (
	"fmt"
	"net/http"
)

/*LogRequests returns an http.Handler middleware that prints method and route.
It should be extended later to be sensitive do a DEV environment variable*/
func LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		fmt.Println(req.Method, req.URL)
		next.ServeHTTP(res, req)
	})
}
