package main

import "fmt"
import "net/http"
import "log"
import "os"
import "time"
// import "io/ioutil"
import "encoding/json"


type Token struct {
    Token   string
}

type TokenController struct {
    TokenValue string
    TokenExpires time.Time
}

func (tc TokenController) getTokenUrl () string {
	url := "http://torrentapi.org/pubapi_v2.php"

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

func (tc *TokenController) getToken() string {
    url := tc.getTokenUrl()
    resp, err := http.Get(url)
    defer resp.Body.Close()

    if err != nil {
        log.Print(err)
        os.Exit(1)
    }

    target := new(Token)
    if resp.StatusCode == http.StatusOK {
        json.NewDecoder(resp.Body).Decode(target)
    }

    return target.Token
}

func (tc *TokenController) DoAuth(){
    tc.TokenValue = tc.getToken()
    tc.TokenExpires = time.Now().Local().Add(time.Minute * time.Duration(14))    
}

func (tc *TokenController) Auth() {
    if tc.TokenValue == "" {
        if tc.TokenExpires >= time.Now().Local() {
            log.Print("Token expired. Getting new one.")
            tc.DoAuth()
        }
    } else {
        log.Print("Request new token.")
        tc.DoAuth()
    }
}

func main() {
    tc := TokenController{}
    tc.Auth()

	fmt.Println(tc.TokenValue)
	fmt.Println(tc.TokenExpires)
}
