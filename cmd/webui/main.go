package main

import (
	hdrwatch "github.com/vegasq/HDRWatch"
	"encoding/json"
	"io"
	"net/http"
	"time"
	"log"
)

func getHDRReleases(w http.ResponseWriter, r *http.Request) {
	sc := hdrwatch.SearchController{}
	sc.Tick = time.Now().Local()

	releases := hdrwatch.ToTable(sc.SearchInCategory("", 52))
	releasesJSON, err := json.Marshal(releases)
	if err != nil {
		log.Fatal("Cannot encode to JSON ", err)
	}

	io.WriteString(w, string(releasesJSON))
}

type IndexHandler struct {}
func (ih IndexHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	_, err := rw.Write([]byte(html))
	if err != nil {
		panic(err)
	}
}

func main() {
	http.Handle("/", IndexHandler{})
	http.HandleFunc("/v1/getHDRReleases", getHDRReleases)

	http.ListenAndServe(":80", nil)
}
