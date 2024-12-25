package main

import (
	"fmt"
	"net/http"
)

func main() {

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "127.0.0.1", nil)
	res, _ := client.Do(req)
	if res.StatusCode != 200 {
		fmt.Println("error")
	}
}
