package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func Get404(path string) (appErr *ApplicationError) {
	msg := "404 Not Found Yo"
	err := errors.New("User attempted to access " + path)
	appErr = NewApplicationError(msg, err, ErrCodeNotFoundFile)
	return appErr
}

// Handler for /js path
func JSHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		var appErr *ApplicationError

		role, appErr := GetHighestRoleFromRequest(r)
		if appErr != nil {
			fmt.Println(appErr)
			HttpErrorLogger(w, appErr.Msg, appErr.Code)
			return
		}

		if strings.HasPrefix(path, "/js/superadmin/") && !CompareRole(role, RoleSuperAdmin) {
			appErr = Get404(path)
		}

		if strings.HasPrefix(path, "/js/admin/") && !CompareRole(role, RoleAdmin) {
			appErr = Get404(path)
		}

		if strings.HasPrefix(path, "/js/captain/") && !CompareRole(role, RoleCaptain) {
			appErr = Get404(path)
		}

		if !strings.HasSuffix(path, ".js") && !strings.HasSuffix(path, ".map") {
			appErr = Get404(path)
		}

		if appErr != nil {
			fmt.Println(appErr)
			HttpErrorLogger(w, appErr.Msg, appErr.Code)
			return
		}

		fmt.Println(path[1:])

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		http.ServeFile(w, r, path[1:])
	}
}
