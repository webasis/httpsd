package main

import (
	"log"
	"net/http"
	"os"

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
	http.Handle("/", handlers.CompressHandler(http.FileServer(&StaticFileSystem{
		fs:           http.Dir(Getenv("serve_path", "./dist/")),
		NotFoundFile: Getenv("not_found_file", "index.html"),
	})))

	log.Fatal(http.ListenAndServeTLS(":443", Getenv("cert_file", "ws.mofon.top.cert"), Getenv("key_file", "ws.mofon.top.key"), nil))
}
