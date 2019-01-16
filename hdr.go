package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/vegasq/GoPrintTable"
)

// import "io/ioutil"

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

type torrentAPIResponse struct {
	TorrentResults []Movie `json:"torrent_results"`
}

type Movie struct {
	Title       string `json:"title"`
	category    string
	Download    string `json:"download"`
	Seeders     int    `json:"seeders"`
	Leechers    int    `json:"leechers"`
	Size        int    `json:"size"`
	PubDate     string `json:"pubdate"`
	Ranked      int    `json:"ranked"`
	InfoPage    string `json:"info_page"`
	EpisodeInfo struct {
		IMDB       string `json:"imdb"`
		TheMovieDB string `json:"themoviedb"`
	} `json:"episode_info"`
}

type SearchController struct {
	Tick time.Time
}

func (sc *SearchController) canIMakeApiCall() bool {
	currentTime := time.Now().Local()
	tickTimePlus := sc.Tick.Add(time.Second * time.Duration(3))

	if currentTime.After(tickTimePlus) {
		sc.Tick = currentTime
		return true
	}

	return false
}

func (sc *SearchController) createSearchURL(searchString string, token string) string {
	url := "https://torrentapi.org/pubapi_v2.php"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := req.URL.Query()

	q.Add("search_string", searchString)
	q.Add("app_id", "hdr_watch")
	q.Add("category", "52")
	q.Add("token", token)
	q.Add("min_seeders", "1")
	q.Add("min_leechers", "1")
	q.Add("limit", "100")
	q.Add("format", "json_extended")
	q.Add("ranked", "0")
	req.URL.RawQuery = q.Encode()

	return req.URL.String()
}

func (sc *SearchController) Search(searchString string) [][]string {
	tc := tokenController{}

	token := tc.GetToken()
	url := sc.createSearchURL(searchString, token)
	log.Print(url)

	// Wait for API
	for sc.canIMakeApiCall() == false {
	}
	resp, err := http.Get(url)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatal("Incorrect return code from ", url, " Code: ", resp.StatusCode)
	}

	target := new(torrentAPIResponse)
	json.NewDecoder(resp.Body).Decode(target)

	tbl := [][]string{}
	tbl = append(tbl, []string{"Movie", "Year", "Seeders", "Released", "Tags", "Download"})
	for _, movie := range target.TorrentResults {
		movieMeta := extractMovieMetadata(movie.Title)

		tbl = append(tbl, []string{
			movieMeta.Name,
			strconv.Itoa(movieMeta.Year),
			strconv.Itoa(movie.Seeders),
			movie.PubDate,
			movieMeta.Tags,
			movie.Download})
	}
	GoPrintTable.PrintTableWithHeader(tbl)

	return tbl

	// responseData, err := ioutil.ReadAll(resp.Body)
	// log.Print("responseData: ", responseData)
	// log.Print("resp.StatusCode: ", resp.StatusCode)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// return string(responseData)

}

// func main2() {
// 	sc := SearchController{}
// 	sc.Tick = time.Now().Local()

// 	sc.Search("")
// }
