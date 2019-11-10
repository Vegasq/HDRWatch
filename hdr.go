package hdrwatch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/vegasq/GoPrintTable"
)


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
	Token string
	Category int
}

func (sc *SearchController) waitForNextAPIWindow() {
	for i := 0; i < 3; i++ {
		currentTime := time.Now().Local()
		tickTimePlus := sc.Tick.Add(time.Second * time.Duration(3))

		if currentTime.After(tickTimePlus) {
			sc.Tick = currentTime
			return
		}

		time.Sleep(1 * time.Second)
	}
}

func (sc *SearchController) createSearchURL(searchString string) string {
	url := "https://torrentapi.org/pubapi_v2.php"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := req.URL.Query()

	q.Add("search_string", searchString)
	q.Add("app_id", "hdr_watch")
	q.Add("category", strconv.Itoa(sc.Category))
	q.Add("token", sc.Token)
	q.Add("min_seeders", "1")
	q.Add("min_leechers", "1")
	q.Add("limit", "100")
	q.Add("format", "json_extended")
	q.Add("ranked", "0")
	req.URL.RawQuery = q.Encode()

	return req.URL.String()
}

func httpCall(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatal("Incorrect return code from ", url, " Code: ", resp.StatusCode)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return data
}

func ToTable(target torrentAPIResponse) [][]string {
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
}

func (sc *SearchController) Search(searchString string) torrentAPIResponse {
	sc.Token = GetToken()
	url := sc.createSearchURL(searchString)
	sc.waitForNextAPIWindow()
	fmt.Println(url)

	response := httpCall(url)
	fmt.Println(response)

	target := torrentAPIResponse{}
	json.NewDecoder(bytes.NewReader(response)).Decode(&target)

	return target
}

func (sc *SearchController) SearchInCategory(searchString string, category int) torrentAPIResponse {
	sc.Category = category
	r := sc.Search(searchString)
	fmt.Println(r)
	return r
}

func Search(searchFor string) torrentAPIResponse {
	sc := SearchController{}
	sc.Tick = time.Now().Local()
	return sc.Search(searchFor)
}