package trap

import "net/http"

/*HandleInternalServerError checks for non-nil errors, and, if found, sends error data
via the ResponseWriter*/
func (s *Service) HandleInternalServerError(e error, res http.ResponseWriter) bool {
	// aught to be should be extended to idiomatically handle more errors
	if e != nil {
		res.WriteHeader(500)
		res.Write([]byte(e.Error()))
		return true
	}
	return false
}

/*HandleUnprocessableEntity checks for non-nil errors, and, if found, sends error data
via the ResponseWriter*/
func (s *Service) HandleUnprocessableEntity(e error, res http.ResponseWriter) bool {
	// aught to be should be extended to idiomatically handle more errors
	if e != nil {
		res.WriteHeader(422)
		res.Write([]byte(e.Error()))
		return true
	}
	return false
}

// func (s *Service) HandleErrors(e error, res http.ResponseWriter, code int) bool {
// 	if e != nil {
// 		res.WriteHeader(code)
// 		res.Write([]byte(e.Error()))
// 		return true
// 	}
// 	return false
// }
