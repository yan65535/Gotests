package main

import (
	"net/http"

	"github.com/aeilang/httpz"
	"github.com/aeilang/httpz/middleware"
)

func main() {
	// Create a new mux
	mux := httpz.NewServeMux()

	// add logger middleware, it 's copy from chi/middleware
	mux.Use(middleware.Logger)

	// register a GET /hello route
	// GET /hello
	mux.Get("/hello", func(w http.ResponseWriter, r *http.Request) error {

		// rw is a helper responsewriter to send response
		rw := httpz.NewHelperRW(w)
		return rw.String(http.StatusOK, "hello httpz")

		// or you can write it by yourself.
		// hw.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		// hw.WriteHeader(http.StatusOK)
		// hw.Write([]byte("hello httpz"))
		// return nil
	})

	// just like net/http's ServerMux
	http.ListenAndServe(":8080", mux)
}
