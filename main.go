package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {

	http.HandleFunc("/", hello)
	server := &http.Server{
		Addr: ":8888",
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("server startup failed, err:%v\n", err)
	}
}

func hello(w http.ResponseWriter, _ *http.Request) {

	w.Write([]byte(fmt.Sprint("Now time is:", time.Now())))
}
