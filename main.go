package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var flagHost = flag.String("host", ":1234", "the host to listen to")
var flagRawPath = flag.String("raw_path", "/var/bigbluebutton/recording/raw/", "the base path of raw folder")
var flagPublisedPath = flag.String("published_path", "/var/bigbluebutton/published/presentation/", "the base path of published folder")

type output map[string]interface{}

func main() {
	flag.Parse()

	r := http.NewServeMux()
	r.HandleFunc("/meeting/", func(w http.ResponseWriter, r *http.Request) {
		var size int64
		id := strings.TrimPrefix(r.URL.Path, "/meeting/")

		if id == "" {
			renderJSON(w, r, 500, output{
				"id":    id,
				"size":  0,
				"error": "id is missing",
			})
			return
		}

		if s, err := folderSize(*flagRawPath + id); err == nil {
			size += s
		} else {
			renderJSON(w, r, 500, output{
				"id":    id,
				"size":  0,
				"error": *flagRawPath + id + " does not exist",
			})
			return
		}

		if s, err := folderSize(*flagPublisedPath + id); err == nil {
			size += s
		} else {
			renderJSON(w, r, 500, output{
				"id":    id,
				"size":  0,
				"error": *flagPublisedPath + id + " does not exist",
			})
			return
		}

		renderJSON(w, r, 200, output{
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
	d, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(d))
	}
}

func folderSize(fpath string) (size int64, err error) {
	err = filepath.Walk(fpath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		size += info.Size()
		return nil
	})

	return
}
