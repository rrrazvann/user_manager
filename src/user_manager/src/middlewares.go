package main

import "net/http"

func mwAllowedMethods(allowedMethods []string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		methodAllowed := false
		for _, method := range allowedMethods {
			if r.Method == method {
				methodAllowed = true
				break
			}
		}

		if !methodAllowed {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

			return
		}

		next(w, r)
	}
}

func mwJsonResponse(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		next(w, r)
	}
}
