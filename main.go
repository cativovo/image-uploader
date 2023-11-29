package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

var (
	homePage       *template.Template = template.Must(template.ParseFiles("web/templates/base.html"))
	imageComponent *template.Template = template.Must(template.ParseFiles("web/templates/components/image.html"))
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	storageDirPath := "images"

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		homePage.Execute(w, nil)
	})

	r.Post("/upload", func(w http.ResponseWriter, r *http.Request) {
		if stat, err := os.Stat(storageDirPath); err != nil || !stat.IsDir() {
			log.Println(err)
			os.Mkdir(storageDirPath, 0777)
		}

		f, fHeaders, err := r.FormFile("image")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer f.Close()

		ext := path.Ext(fHeaders.Filename)
		filename := uuid.NewString() + ext

		dst, err := os.Create(path.Join(storageDirPath, filename))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, f); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		protocol := "http"

		if r.TLS != nil {
			protocol += "s"
		}

		filePath := fmt.Sprintf("%s://%s/images/%s", protocol, r.Host, filename)

		w.WriteHeader(http.StatusCreated)
		imageComponent.Execute(w, map[string]any{
			"ImageSrc": filePath,
			"Filename": filename,
		})
	})

	r.Get("/images/*", func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(http.Dir(storageDirPath)))

		fs.ServeHTTP(w, r)
	})

	r.Get("/public/*", func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(http.Dir("public")))

		fs.ServeHTTP(w, r)
	})

	http.ListenAndServe(":4000", r)
}
