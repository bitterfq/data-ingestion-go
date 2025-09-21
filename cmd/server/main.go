package main

import (
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	resp, err := http.Get("http://github.com/bitterfq")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Response Status:", resp.Status)
	//fmt.Println("Response Headers:", resp.Header)
	fmt.Println("Response Body:", resp.Body)

}
