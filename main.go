package main

import (
	"flag"
	"math"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

var flagHost = flag.String("host", ":1234", "the host to listen to")
var flagRawPath = flag.String("raw_path", "/var/bigbluebutton/recording/raw/", "the base path of raw folder")
var flagPublisedPath = flag.String("published_path", "/var/bigbluebutton/published/presentation/", "the base path of published folder")

func main() {
	r := gin.Default()
	flag.Parse()

	r.GET("record/:id", func(ctx *gin.Context) {
		var size int64
		id := ctx.Param("id")

		if s, err := folderSize(*flagRawPath + id); err == nil {
			size += s
		} else {
			ctx.JSON(500, gin.H{
				"id":    id,
				"size":  0,
				"error": *flagRawPath + id + " does not exist",
			})
			return
		}

		if s, err := folderSize(*flagPublisedPath + id); err == nil {
			size += s
		} else {
			ctx.JSON(500, gin.H{
				"id":    id,
				"size":  0,
				"error": *flagPublisedPath + id + " does not exist",
			})
			return
		}

		ctx.JSON(200, gin.H{
			"id":    id,
			"size":  round(float64(size)/1024, 2), // kb
			"error": "",
		})
	})

	r.Run(*flagHost)
}

func round(f float64, n int) float64 {
	p10 := math.Pow10(n)
	return math.Trunc(f*p10+0.5) / p10
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
