package main

import (
	"fmt"
	"net/http"
)

func main() {
	server := &http.Server{
		Addr:              ":8080",
		Handler:           nil,
		TLSConfig:         nil,
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error")
		return
	}
}
