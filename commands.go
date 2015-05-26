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
	"github.com/k0kubun/pp"
	"io"
	"strings"
)

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

const (
	ReqURL = "http://localhost:9000/"
)

type Tweet struct {
	TweetID          int    `json:"tweetId"`
	Text             string `json:"text"`
	MemberID         string `json:"memberId"`
	TimestampCreated int    `json:"timestampCreated"`
	TimestampUpdated int    `json:"timestampUpdated"`
}

/**
POST /authenticate
*/
func doLogin(c *cli.Context) {
	url := ReqURL + "api/authenticate"
	pp.Println("URL:>", url)

	jsonStr := []byte(`{"name": "asd", "password": "asd"}`)

	client := &http.Client{}
	jar, _ := cookiejar.New(nil)
	client.Jar = jar
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}

	pp.Println("response Status:", resp.Status)
	pp.Println("response Headers:", resp.Header)

	if strings.Contains(resp.Status, "400") {
		body, _ := ioutil.ReadAll(resp.Body)
		pp.Println("response Body:", string(body))
	} else {
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
	sessionID, _ := sessionID()
	req, err := http.NewRequest("GET", ReqURL+"api/tweets"+url.Values{}.Encode(), nil)
	req.Header.Set("Cookie", sessionID)

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
	following := c.Args().First()
	sessionID, _ := sessionID()
	req, err := http.NewRequest("POST", ReqURL+"api/follow/"+following+url.Values{}.Encode(), nil)
	req.Header.Set("Cookie", sessionID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	defer resp.Body.Close()
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

func sessionID() (sID string, err error) {
	file := "mutterTemp/sessionID"
	fl, err := os.Open(file)

	if err == nil {
		buf := bytes.NewBuffer(nil)
		io.Copy(buf, fl)
		sID = string(buf.Bytes())
	}

	fl.Close()
	return sID, err
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

	for _, v := range js.MustArray() {
		tw := v.(map[string]interface{})
		fmt.Println(tw["memberId"], tw["text"])
	}
}
