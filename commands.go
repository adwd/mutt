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
	"encoding/json"
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
}

var commandLogin = cli.Command{
	Name:  "login",
	Usage: "",
	Description: `
`,
	Action: doLogin,
}

var commandLogout = cli.Command{
	Name:  "logout",
	Usage: "",
	Description: `
`,
	Action: doLogout,
}

var commandRegister = cli.Command{
	Name:  "register",
	Usage: "",
	Description: `
`,
	Action: doRegister,
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

type Tweet struct {
	TweetID          int    `json:"tweetId"`
	Text             string `json:"text"`
	MemberID         string `json:"memberId"`
	TimestampCreated int    `json:"timestampCreated"`
	TimestampUpdated int    `json:"timestampUpdated"`
}

type Config struct {
	URL       string
	SessionID string
}

func saveConfig(conf *Config) (err error) {
	os.Mkdir("mutterTemp", 0777)
	file := "mutterTemp/sessionID"

	b, err := json.Marshal(conf)
	if err == nil {
		err = ioutil.WriteFile(file, b, 0644)
	}

	return err
}

func loadConfig() (conf Config, err error) {
	filename := "mutterTemp/sessionID"
	file, err := ioutil.ReadFile(filename)

	if err == nil {
		err = json.Unmarshal(file, &conf)
	}

	return conf, err
}

/**
POST /authenticate
*/
func doLogin(c *cli.Context) {
	reader := bufio.NewReader(os.Stdin)

	// URL
	conf, _ := loadConfig()
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

	url := conf.URL + "api/authenticate"
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
		conf.SessionID = sessionID
		err := saveConfig(&conf)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	defer resp.Body.Close()
}

func doLogout(c *cli.Context) {
	conf, _ := loadConfig()
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
	conf, _ := loadConfig()
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
		pp.Println("response Body:", body)
	} else {
		sessionID := strings.Split(resp.Header["Set-Cookie"][0], "; ")[0]
		conf.SessionID = sessionID
		err := saveConfig(&conf)
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

	conf, _ := loadConfig()
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
	conf, _ := loadConfig()
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
	conf, _ := loadConfig()
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
	conf, _ := loadConfig()
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
	conf, _ := loadConfig()
	resp, err := http.Get(conf.URL + "api/recents" + url.Values{}.Encode())
	if err != nil {
		fmt.Println(err)
	}

	displayTweets(resp)
	defer resp.Body.Close()
}

func sessionID() (sID string, err error) {
	conf, err := loadConfig()
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
