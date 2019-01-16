package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

func getHDRReleases(w http.ResponseWriter, r *http.Request) {
	sc := SearchController{}
	sc.Tick = time.Now().Local()

	releases := sc.Search("")
	releasesJSON, err := json.Marshal(releases)
	if err != nil {
		log.Fatal("Cannot encode to JSON ", err)
	}

	io.WriteString(w, string(releasesJSON))
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./webui")))
	http.HandleFunc("/v1/getHDRReleases", getHDRReleases)

	http.ListenAndServe(":80", nil)
}
