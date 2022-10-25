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
	ctx := r.Context()
	fmt.Println("Request received!")

	select {
	case <-ctx.Done():
		ctxErr := ctx.Err()
		fmt.Println(ctxErr)
		http.Error(w, ctxErr.Error(), http.StatusInternalServerError)
	case <-time.After(time.Second * 5):
		fmt.Println("Response sent!")
		fmt.Fprintf(w, "SayHey handler said hey!")
	}
}
