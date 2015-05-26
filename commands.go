package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"

	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"bytes"
	"github.com/bitly/go-simplejson"
	"io"
	"strings"
)

const (
	ReqURL = "http://localhost:9000/"
)

type Tweet struct {
	TweetId          int    `json:"tweetId"`
	Text             string `json:"text"`
	MemberId         string `json:"memberId"`
	TimestampCreated int    `json:"timestampCreated"`
	TimestampUpdated int    `json:"timestampUpdated"`
}

var Commands = []cli.Command{
	commandLogin,
	commandTweet,
	commandShow,
	commandRecommends,
	commandFollow,
	commandRecents,
}

var commandLogin = cli.Command{
	Name:  "login",
	Usage: "",
	Description: `
`,
	Action: doLogin,
}

var commandTweet = cli.Command{
	Name:  "tweet",
	Usage: "",
	Description: `
`,
	Action: doTweet,
}

var commandShow = cli.Command{
	Name:  "show",
	Usage: "",
	Description: `
`,
	Action: doShow,
}

var commandRecommends = cli.Command{
	Name:  "recommends",
	Usage: "",
	Description: `
`,
	Action: doRecommends,
}

var commandFollow = cli.Command{
	Name:  "follow",
	Usage: "",
	Description: `
`,
	Action: doFollow,
}

var commandRecents = cli.Command{
	Name:  "recents",
	Usage: "",
	Description: `
`,
	Action: doRecents,
}

func debug(v ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		log.Println(v...)
	}
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

/**
POST /authenticate
*/
func doLogin(c *cli.Context) {
	url := ReqURL + "api/authenticate"
	fmt.Println("URL:>", url)

	jsonStr := []byte(`{"name": "asd", "password": "asd"}`)

	client := &http.Client{}
	jar, _ := cookiejar.New(nil)
	client.Jar = jar
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	if strings.Contains(resp.Status, "400") {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
	} else {
		fmt.Println("Set-Cookie:", resp.Header["Set-Cookie"])

		sessionID := strings.Split(resp.Header["Set-Cookie"][0], "; ")[0]
		os.Mkdir("mutterTemp", 0777)
		file := "mutterTemp/sessionID"
		fout, err := os.Create(file)
		if err != nil {
			fmt.Println(file, err)
			return
		}

		defer fout.Close()
		fout.WriteString(sessionID)
	}

	defer resp.Body.Close()
}

func doTweet(c *cli.Context) {
}

/**
GET /
*/
func doShow(c *cli.Context) {
	file := "mutterTemp/sessionID"
	fl, err := os.Open(file)
	if err != nil {
		fmt.Println(file, err)
		return
	}

	defer fl.Close()

	buf := bytes.NewBuffer(nil)
	io.Copy(buf, fl)

	req, err := http.NewRequest("GET", ReqURL+"api/tweets"+url.Values{}.Encode(), nil)
	req.Header.Set("Cookie", string(buf.Bytes()))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	displayTweets(resp)
	defer resp.Body.Close()
}

func doRecommends(c *cli.Context) {
}

func doFollow(c *cli.Context) {
}

// GET /api/recents
func doRecents(c *cli.Context) {
	resp, err := http.Get(ReqURL + "api/recents" + url.Values{}.Encode())
	if err != nil {
		fmt.Println(err)
	}

	displayTweets(resp)
	defer resp.Body.Close()
}

func displayTweets(resp *http.Response) {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	js, err := simplejson.NewJson(b)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i, v := range js.MustArray() {
		tw := v.(map[string]interface{})
		fmt.Println("%d, %s, %s", i, tw["memberId"], tw["text"])
	}
}
