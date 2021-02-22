package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	port = 4273
)

func main() {
	http.DefaultServeMux.HandleFunc(
		"/users",
		mwJsonResponse(
			mwAllowedMethods(
				[]string{"GET", "POST"},
				getAndAddUsers,
			),
		),
	)

	http.DefaultServeMux.HandleFunc(
		"/users/",
		mwJsonResponse(
			mwAllowedMethods(
				[]string{"DELETE", "PUT"},
				deleteAndUpdateUser,
			),
		),
	)

	http.DefaultServeMux.HandleFunc(
		"/health",
		mwJsonResponse(
			mwAllowedMethods(
				[]string{"GET"},
				healthCheck,
			),
		),
	)

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	log.Println(fmt.Sprintf("Private server started on port: %s", addr))

	err := http.ListenAndServe(addr, http.DefaultServeMux)
	log.Fatal(err)
}
