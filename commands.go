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

	"bufio"
	"bytes"
	"github.com/bitly/go-simplejson"
	"github.com/k0kubun/pp"
	"strings"
)

var Commands = []cli.Command{
	commandLogin,
	commandLogout,
	commandRegister,
	commandTweet,
	commandShow,
	commandRecommends,
	commandFollow,
	commandRecents,
	commandClear,
}

var commandLogin = cli.Command{
	Name:  "login",
	Usage: "Logins to mutter.",
	Description: `mutter creates config file(.mutterrc) to record URL of mutter and sessionID.
`,
	Action: doLogin,
}

var commandLogout = cli.Command{
	Name:  "logout",
	Usage: "Logouts from mutter.",
	Description: `
`,
	Action: doLogout,
}

var commandRegister = cli.Command{
	Name:  "register",
	Usage: "Creates a user account.",
	Description: `
`,
	Action: doRegister,
}

var commandTweet = cli.Command{
	Name:  "tweet",
	Usage: "Posts a tweet.",
	Description: `mutter tweet "hello mutter"
`,
	Action: doTweet,
}

var commandShow = cli.Command{
	Name:  "show",
	Usage: "Shows tweets by your following-users and you.",
	Description: `
`,
	Action: doShow,
}

var commandRecommends = cli.Command{
	Name:  "recommends",
	Usage: "Shows a list of users you don't follow.",
	Description: `
`,
	Action: doRecommends,
}

var commandFollow = cli.Command{
	Name:  "follow",
	Usage: "Follows a user.",
	Description: `mutter follow "someone"
`,
	Action: doFollow,
}

var commandRecents = cli.Command{
	Name:  "recents",
	Usage: "Shows recent tweets by all users.",
	Description: `
`,
	Action: doRecents,
}

var commandClear = cli.Command{
	Name:  "clear",
	Usage: "Cleans mutter config file. (.mutterrc)",
	Description: `Cleans mutter config file. (.mutterrc)
`,
	Action: doClear,
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
	reader := bufio.NewReader(os.Stdin)

	// URL
	if Conf.SessionID == "" {
		fmt.Print("URL: ")
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, "\n")
		if strings.HasSuffix(text, "/") {
			Conf.URL = text
		} else {
			Conf.URL = text + "/"
		}
	}

	url := Conf.URL + "api/authenticate"
	pp.Println("URL:>", url)

	// name, password
	fmt.Print("username: ")
	name, _ := reader.ReadString('\n')
	fmt.Print("password: ")
	pass, _ := reader.ReadString('\n')

	js := simplejson.New()
	js.Set("name", strings.Trim(name, "\n"))
	js.Set("password", strings.Trim(pass, "\n"))
	jsbin, _ := js.MarshalJSON()

	// send request
	client := &http.Client{}
	jar, _ := cookiejar.New(nil)
	client.Jar = jar
	resp, err := client.Post(url, "application/json", bytes.NewReader(jsbin))
	if err != nil {
		panic(err)
	}

	// show response
	pp.Println("response Status:", resp.Status)
	pp.Println("response Headers:", resp.Header)

	if strings.Contains(resp.Status, "400") {
		body, _ := ioutil.ReadAll(resp.Body)
		pp.Println("response Body:", body)
	} else {
		sessionID := strings.Split(resp.Header["Set-Cookie"][0], "; ")[0]
		Conf.SessionID = sessionID
		err := SaveConfig(&Conf)
		if err != nil {
			fmt.Println(err)
			return
		}
		body, _ := ioutil.ReadAll(resp.Body)
		jsBody, _ := simplejson.NewJson(body)
		pp.Println("response Body:", jsBody)

		fmt.Println("login ")
	}

	defer resp.Body.Close()
}

func doLogout(c *cli.Context) {
	conf, _ := LoadConfig()
	req, err := http.NewRequest("POST", conf.URL+"api/logout"+url.Values{}.Encode(), nil)
	req.Header.Set("Cookie", conf.SessionID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	json, _ := simplejson.NewJson(body)
	pp.Println("response Body:", json)
	defer resp.Body.Close()
}

func doRegister(c *cli.Context) {
	reader := bufio.NewReader(os.Stdin)

	// URL
	conf, _ := LoadConfig()
	if conf.SessionID == "" {
		fmt.Print("URL: ")
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, "\n")
		if strings.HasSuffix(text, "/") {
			conf.URL = text
		} else {
			conf.URL = text + "/"
		}
	}

	// name, password
	fmt.Print("username: ")
	name, _ := reader.ReadString('\n')
	fmt.Print("password: ")
	pass, _ := reader.ReadString('\n')
	fmt.Print("mail address: ")
	mail, _ := reader.ReadString('\n')

	js := simplejson.New()
	js.Set("name", strings.Trim(name, "\n"))
	js.Set("password", strings.Trim(pass, "\n"))
	js.Set("mail", strings.Trim(mail, "\n"))
	jsbin, _ := js.MarshalJSON()

	// send request
	client := &http.Client{}
	jar, _ := cookiejar.New(nil)
	client.Jar = jar
	url := conf.URL + "api/create"
	resp, err := client.Post(url, "application/json", bytes.NewReader(jsbin))
	if err != nil {
		panic(err)
	}

	// show response
	pp.Println("response Status:", resp.Status)
	pp.Println("response Headers:", resp.Header)

	if strings.Contains(resp.Status, "400") {
		body, _ := ioutil.ReadAll(resp.Body)
		js, _ := simplejson.NewJson(body)
		pp.Println("response Body:", js)
	} else {
		sessionID := strings.Split(resp.Header["Set-Cookie"][0], "; ")[0]
		conf.SessionID = sessionID
		err := SaveConfig(&conf)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	defer resp.Body.Close()
}

func doTweet(c *cli.Context) {
	text := c.Args().First()
	js := simplejson.New()
	js.Set("text", text)
	jsbin, _ := js.MarshalJSON()

	conf, _ := LoadConfig()
	req, err := http.NewRequest("POST", conf.URL+"api/tweet"+url.Values{}.Encode(), bytes.NewReader(jsbin))
	req.Header.Set("Cookie", conf.SessionID)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	if strings.Contains(resp.Status, "400") {
		body, _ := ioutil.ReadAll(resp.Body)
		js, _ := simplejson.NewJson(body)
		pp.Println("response Body:", js)
		return
	}

	displayTweets(resp)
	defer resp.Body.Close()
}

/**
GET /
*/
func doShow(c *cli.Context) {
	conf, _ := LoadConfig()
	req, err := http.NewRequest("GET", conf.URL+"api/tweets"+url.Values{}.Encode(), nil)
	req.Header.Set("Cookie", conf.SessionID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	displayTweets(resp)
	defer resp.Body.Close()
}

func doRecommends(c *cli.Context) {
	conf, _ := LoadConfig()
	req, err := http.NewRequest("GET", conf.URL+"api/recommends"+url.Values{}.Encode(), nil)
	req.Header.Set("Cookie", conf.SessionID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	bodyjs, _ := simplejson.NewJson(body)
	//pp.Println("response Body:", bodyjs)
	for _, v := range bodyjs.MustArray() {
		tw := v.(map[string]interface{})
		fmt.Println(tw["memberId"], tw["mailAddress"])
	}
	defer resp.Body.Close()
}

func doFollow(c *cli.Context) {
	following := c.Args().First()
	conf, _ := LoadConfig()
	req, err := http.NewRequest("POST", conf.URL+"api/follow/"+following+url.Values{}.Encode(), nil)
	req.Header.Set("Cookie", conf.SessionID)

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
	conf, _ := LoadConfig()
	resp, err := http.Get(conf.URL + "api/recents" + url.Values{}.Encode())
	if err != nil {
		fmt.Println(err)
	}

	displayTweets(resp)
	defer resp.Body.Close()
}

func doClear(c *cli.Context) {
	ClearConfig()
}

func sessionID() (sID string, err error) {
	conf, err := LoadConfig()
	return conf.SessionID, err
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
