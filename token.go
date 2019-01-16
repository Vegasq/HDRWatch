package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type tokenJSONStruct struct {
	Token string
}

type tokenController struct {
	TokenValue   string
	TokenExpires time.Time
}

func (tc tokenController) getTokenURL() string {
	url := "https://torrentapi.org/pubapi_v2.php"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := req.URL.Query()
	q.Add("get_token", "get_token")
	q.Add("format", "json")
	q.Add("app_id", "hdr_watch")
	req.URL.RawQuery = q.Encode()

	return req.URL.String()
}

func (tc *tokenController) getToken() string {
	url := tc.getTokenURL()
	resp, err := http.Get(url)

	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatal("Incorrect return code from ", url, " Code: ", resp.StatusCode)
	}

	target := new(tokenJSONStruct)
	json.NewDecoder(resp.Body).Decode(target)

	return target.Token
}

func (tc *tokenController) doAuth() {
	tc.TokenValue = tc.getToken()
	tc.TokenExpires = time.Now().Local().Add(time.Minute * time.Duration(14))
}

func (tc *tokenController) auth() {
	if tc.TokenValue != "" {
		if tc.TokenExpires.After(time.Now().Local()) {
			log.Print("Token expired. Getting new one.")
			tc.doAuth()
		}
	} else {
		log.Print("Request new token.")
		tc.doAuth()
	}
}

// GetToken return torrentapi token as a string value
func (tc *tokenController) GetToken() string {
	tc.auth()
	return tc.TokenValue
}
