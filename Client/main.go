package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	client := &http.Client{
		CheckRedirect: http.DefaultClient.CheckRedirect,
	}

	// creating root context
	ctx := context.Background()

	// using WithCancel and cancel function
	ctx, cancel := context.WithCancel(ctx)
	time.AfterFunc(time.Second*3, cancel)

	req, err := http.NewRequest("GET", "http://127.0.0.1:8080", nil)
	if err != nil {
		log.Fatal("Error while creating request. Err: ", err)
	}

	req = req.WithContext(ctx)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Error while receiving response. Err: ", err)
	}
	defer res.Body.Close()

	io.Copy(os.Stdout, res.Body)
}
