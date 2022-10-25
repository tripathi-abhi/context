package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {

	http.HandleFunc("/", handleSayHey)

	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		fmt.Println("Error starting server. Error: ", err)
	}
}

func handleSayHey(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second * 5)
	fmt.Fprintf(w, "Hey handler said hey!")
}
