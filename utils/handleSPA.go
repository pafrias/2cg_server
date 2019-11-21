package utils

import (
	"net/http"
	"strings"
)

// single page app helper function, redirecting back past the base url
func HandleRedirect(path string) http.Handler {
	return http.RedirectHandler(path, http.StatusSeeOther)
}

func RequiresRedirect(path string, prefix string) bool {
	if path != prefix && strings.HasPrefix(path, prefix) {
		path = strings.TrimPrefix(path, prefix)
		if strings.Contains(path, ".") {
			return false
		}
		return true
	}
	return false
}
