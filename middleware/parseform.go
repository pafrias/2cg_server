package middleware

import "net/http"

/*ParseForm automatically parses forms for certain requests
Likely an unneccesary abstraction, may delete in final version*/
func ParseForm(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
			if req.Method == "POST" || req.Method == "PATCH" {
				req.ParseForm()
			}
			next.ServeHTTP(res, req)
		})
}
