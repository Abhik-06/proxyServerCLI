package main

import (
	"fmt"
	"net/http"
)

func handleSearch() {}

func main() {
	http.HandleFunc(" /search", handleSearch)
	fmt.Println("Proxy server : Request acknowledged")

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
