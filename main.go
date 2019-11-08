package main

import (
	"log"
	"net/http"
	"os"

	"fmt"

	"github.com/gorilla/handlers"
)

func Getenv(key, defval string) string {
	val := os.Getenv(key)
	if val == "" {
		return defval
	}
	return val
}

type StaticFileSystem struct {
	NotFoundFile string
	fs           http.FileSystem
}

func (fs *StaticFileSystem) Open(name string) (http.File, error) {
	file, err := fs.fs.Open(name)
	if err != nil {
		return fs.fs.Open(fs.NotFoundFile)
	}
	return file, err
}

func main() {
	if len(os.Args) == 2 && os.Args[1] == "help" {
		fmt.Println("ENV:")
		fmt.Println("serve_path")
		fmt.Println("not_found_file")
		fmt.Println("cert_file")
		fmt.Println("key_file")
		fmt.Println("base_path")
		return
	}

	http.Handle(Getenv("base_path", "/"), handlers.CompressHandler(http.FileServer(&StaticFileSystem{
		fs:           http.Dir(Getenv("serve_path", "./dist/")),
		NotFoundFile: Getenv("not_found_file", "index.html"),
	})))

	http_server := &http.Server{}
	http_server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://"+r.Host+"/", http.StatusFound)
	})
	go http_server.ListenAndServe()

	log.Fatal(http.ListenAndServeTLS(":"+Getenv("port", "443"), Getenv("cert_file", "ws.mofon.top.cert"), Getenv("key_file", "ws.mofon.top.key"), nil))
}
