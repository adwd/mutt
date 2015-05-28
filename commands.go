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
	"code.google.com/p/gopass"
	"github.com/bitly/go-simplejson"
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
	commandUnfollow,
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

var commandUnfollow = cli.Command{
	Name:  "unfollow",
	Usage: "Unfollows a user.",
	Description: `mutter unfollow "someone"
`,
	Action: doUnfollow,
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

	DisplayString(cyan, "URL: "+Conf.URL)

	// name, password
	fmt.Print("username: ")
	name, _ := reader.ReadString('\n')
	fmt.Print("password: ")
	pass, _ := gopass.GetPass("")

	js := simplejson.New()
	js.Set("name", strings.Trim(name, "\n"))
	js.Set("password", strings.Trim(pass, "\n"))
	jsbin, _ := js.MarshalJSON()

	// send request
	resp, err := http.Post(Conf.URL+"api/authenticate", "application/json", bytes.NewReader(jsbin))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		sessionID := strings.Split(resp.Header["Set-Cookie"][0], "; ")[0]
		Conf.SessionID = sessionID
		err := SaveConfig(&Conf)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	DisplayResponse(resp)
}

func doLogout(c *cli.Context) {
	simpleReqRes("POST", Conf.URL+"api/logout/")
}

func doRegister(c *cli.Context) {
	reader := bufio.NewReader(os.Stdin)

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

	// name, password
	fmt.Print("username: ")
	name, _ := reader.ReadString('\n')
	fmt.Print("password: ")
	pass, _ := gopass.GetPass("")
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
	url := Conf.URL + "api/create"
	resp, err := client.Post(url, "application/json", bytes.NewReader(jsbin))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		sessionID := strings.Split(resp.Header["Set-Cookie"][0], "; ")[0]
		Conf.SessionID = sessionID
		err := SaveConfig(&Conf)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	DisplayResponse(resp)
}

func doTweet(c *cli.Context) {
	text := c.Args().First()
	js := simplejson.New()
	js.Set("text", text)
	jsbin, _ := js.MarshalJSON()

	req, err := http.NewRequest("POST", Conf.URL+"api/tweet"+url.Values{}.Encode(), bytes.NewReader(jsbin))
	req.Header.Set("Cookie", Conf.SessionID)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		DisplayResponse(resp)
		return
	}

	DisplayTweets(resp)
}

/**
GET /
*/
func doShow(c *cli.Context) {
	req, err := http.NewRequest("GET", Conf.URL+"api/tweets"+url.Values{}.Encode(), nil)
	req.Header.Set("Cookie", Conf.SessionID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		DisplayResponse(resp)
		return
	}

	DisplayTweets(resp)
}

func doRecommends(c *cli.Context) {
	req, err := http.NewRequest("GET", Conf.URL+"api/recommends"+url.Values{}.Encode(), nil)
	req.Header.Set("Cookie", Conf.SessionID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	bodyjs, _ := simplejson.NewJson(body)

	for _, v := range bodyjs.MustArray() {
		tw := v.(map[string]interface{})
		fmt.Println(tw["memberId"], tw["mailAddress"])
	}
}

func doFollow(c *cli.Context) {
	following := c.Args().First()
	simpleReqRes("POST", Conf.URL+"api/follow/"+following)
}

func doUnfollow(c *cli.Context) {
	following := c.Args().First()
	simpleReqRes("POST", Conf.URL+"api/unfollow/"+following)
}

// GET /api/recents
func doRecents(c *cli.Context) {
	resp, err := http.Get(Conf.URL + "api/recents" + url.Values{}.Encode())
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	DisplayTweets(resp)
}

func doClear(c *cli.Context) {
	ClearConfig()
}

// urlに中身のないリクエストを、指定したmethodで送り、
// レスポンスのボディのtextを表示する
func simpleReqRes(method, url string) {
	req, err := http.NewRequest(method, url, nil)
	req.Header.Set("Cookie", Conf.SessionID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	DisplayResponse(resp)
}
