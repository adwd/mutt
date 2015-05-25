package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"

	"io/ioutil"
	"net/http"
	"net/url"

	"bytes"
	"github.com/bitly/go-simplejson"
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

	jsonString := `{"name": "qwe", "password": "qwe"}`

	var jsonStr = []byte(jsonString)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	//req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func doTweet(c *cli.Context) {
}

/**
GET /recents
*/
func doShow(c *cli.Context) {

	values := url.Values{}

	getSimple(values)

	fmt.Println("hello")
}

func getSimple(values url.Values) {
	resp, err := http.Get(ReqURL + "api/recents" + values.Encode())
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	execute(resp)
}

func execute(resp *http.Response) {
	b, err := ioutil.ReadAll(resp.Body)
	if err == nil {

		js, err2 := simplejson.NewJson(b)
		if err2 == nil {

			for i, v := range js.MustArray() {
				tw := v.(map[string]interface{})
				fmt.Println("%d, %s, %s\n", i, tw["memberId"], tw["text"])
			}
		}
	}
}

func doRecommends(c *cli.Context) {
}

func doFollow(c *cli.Context) {
}

func doRecents(c *cli.Context) {
}
