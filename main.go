package main

import (
	"encoding/json"
	"flag"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var flagHost = flag.String("host", ":1234", "the host to listen to")
var flagRawPath = flag.String("raw_path", "/var/bigbluebutton/recording/raw/", "the base path of raw folder")
var flagPublisedPath = flag.String("published_path", "/var/bigbluebutton/published/presentation/", "the base path of published folder")
var flagPathName = flag.String("path_name", "record", "the base path in URL")

type output map[string]interface{}

func main() {
	flag.Parse()

	pathPrefix := "/" + *flagPathName + "/"
	r := http.NewServeMux()

	r.HandleFunc(pathPrefix, func(w http.ResponseWriter, r *http.Request) {
		var size int64
		id := strings.TrimPrefix(r.URL.Path, pathPrefix)

		if id == "" {
			renderJSON(w, r, http.StatusInternalServerError, output{
				"id":    id,
				"size":  0,
				"error": "id is missing",
			})
			return
		}

		if s, err := folderSize(*flagRawPath + id); err == nil {
			size += s
		} else {
			renderJSON(w, r, http.StatusInternalServerError, output{
				"id":    id,
				"size":  0,
				"error": err.Error(),
			})
			return
		}

		if s, err := folderSize(*flagPublisedPath + id); err == nil {
			size += s
		} else {
			renderJSON(w, r, http.StatusInternalServerError, output{
				"id":    id,
				"size":  0,
				"error": err.Error(),
			})
			return
		}

		renderJSON(w, r, http.StatusOK, output{
			"id":    id,
			"size":  round(float64(size)/1024, 2), // kb
			"error": "",
		})
	})

	http.ListenAndServe(*flagHost, r)
}

func round(f float64, n int) float64 {
	p10 := math.Pow10(n)
	return math.Trunc(f*p10+0.5) / p10
}

func renderJSON(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func folderSize(fpath string) (size int64, err error) {
	err = filepath.Walk(fpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		size += info.Size()
		return nil
	})

	return
}
