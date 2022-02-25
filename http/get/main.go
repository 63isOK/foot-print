package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	resp, err := http.Get("http://www.baidu.111com")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode > 299 {
		log.Fatalf("code:%d, body:%s", resp.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}

	println(string(body))
}
