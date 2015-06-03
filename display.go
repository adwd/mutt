package main

import (
	"fmt"
	"os"
	"strings"

	"io/ioutil"
	"net/http"

	"github.com/bitly/go-simplejson"
	"github.com/mgutz/ansi"
	"github.com/moznion/go-text-visual-width"
	"golang.org/x/crypto/ssh/terminal"
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

	// pretty print
	w, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println(err)
		return
	}

	if w < 20 {
		fmt.Println("terminal too narrow to show tweets")
		return
	}

	// ツイートをひとつずつ表示する。
	// その際、名前をツイートを左右に分けて表示する。
	for _, v := range js.MustArray() {
		tw := v.(map[string]interface{})
		name, _ := tw["memberId"].(string)
		text, _ := tw["text"].(string)
		texts := strings.Split(text, "\n")

		for _, v := range texts {
			prettyPrintTweet(name, v, w)
			name = ""
		}
	}
}

func prettyPrintTweet(name, text string, width int) {
	rows := visualwidth.Width(text)/(width-17) + 1
	var textstr = text
	var textprint string
	var namestr = name
	for i := 0; i < rows; i++ {
		// make and print name string
		blanks := 16 - visualwidth.Width(namestr)
		fmt.Print(lime, namestr, reset)
		for j := blanks; j > 0; j-- {
			fmt.Print(reset, " ", reset)
		}
		namestr = ""

		// trim and save next for tweeted text
		textprint, textstr = visualwidth.Separate(textstr, width-18)
		fmt.Println(yellow, textprint, reset)
	}
}
