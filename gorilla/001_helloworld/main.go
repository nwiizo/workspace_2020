package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// ハンドラー。処理を記述する。
func helloHandler1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}

func helloHandler2(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, "hello"+" "+vars["name"])
}

func main() {
	// ルーティング設定
	r := mux.NewRouter()
	r.HandleFunc("/hello", helloHandler1)
	r.HandleFunc("/hello/{name}", helloHandler2)

	// サーバ設定
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// 起動
	log.Fatal(srv.ListenAndServe())
}
