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

	base_path := Getenv("base_path", "/")
	serve_path := Getenv("serve_path", "./dist/")
	not_found_file := Getenv("not_found_file", "index.html")
	port := Getenv("port", "443")
	cert_file := Getenv("cert_file", "ws.mofon.top.cert")
	key_file := Getenv("key_file", "ws.mofon.top.key")

	fmt.Println("base_path:", base_path)
	fmt.Println("serve_path:", serve_path)
	fmt.Println("not_found_file:", not_found_file)
	fmt.Println("port:", port)
	fmt.Println("cert_file:", cert_file)
	fmt.Println("key_file:", key_file)

	http.Handle(base_path, handlers.CompressHandler(http.FileServer(&StaticFileSystem{
		fs:           http.Dir(serve_path),
		NotFoundFile: not_found_file,
	})))

	http_server := &http.Server{}
	http_server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://"+r.Host+"/", http.StatusFound)
	})
	go http_server.ListenAndServe()

	log.Fatal(http.ListenAndServeTLS(":"+port, cert_file, key_file, nil))
}
