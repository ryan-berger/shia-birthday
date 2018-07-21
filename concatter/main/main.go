package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"github.com/ryan-berger/shia-birthday/concatter"
	"fmt"
	"io"
	"os"
)

func main() {
	pool := concatter.NewWorkerPool()
	r := mux.NewRouter()
	r.HandleFunc("/shia-surprise/", func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err != nil {
			// Handle error
		}
		fmt.Println(request.Form["response_url"])
		pool.MakeRequest(request.Form["text"][0], request.Form["response_url"][0])
	})

	r.HandleFunc("/gif/{name}", func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		fmt.Println(vars)
		file, e := os.Open(fmt.Sprintf("gifs/%s.gif", vars["name"]))
		if e != nil {
			panic(e)
		}
		request.Header.Set("Content-Type", "image/gif")
		io.Copy(writer, file)
	})

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8111", r))
}
