package main

import (
	"fmt"

	"io/ioutil"
	"net/http"

	"github.com/bitly/go-simplejson"
	"github.com/mgutz/ansi"
)

// cache escape codes and build strings
var lime = ansi.ColorCode("green+buh")
var yellow = ansi.ColorCode("yellow+h")
var reset = ansi.ColorCode("reset")
var red = ansi.ColorCode("red+h")
var white = ansi.ColorCode("white")
var cyan = ansi.ColorCode("cyan+h")

func DisplayResponse(r *http.Response) {
	colorCode := white

	switch r.StatusCode {
	case http.StatusOK:
		colorCode = cyan
	default:
		colorCode = red
	}

	contents, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Println(colorCode, string(contents), reset)
}

func DisplayString(code, message string) {
	fmt.Println(code, message, reset)
}

func DisplayTweets(resp *http.Response) {
	fmt.Println(cyan, "tweets", reset)

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
		fmt.Print(lime, tw["memberId"], reset)
		fmt.Println(yellow, tw["text"], reset)
	}
}
