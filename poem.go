package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type AlfredItem struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Arg      string `json:"arg"`
	Icon     struct {
		Path string `json:"path"`
	} `json:"icon"`
}

func getItems(query string) {
	res, err := http.Get(fmt.Sprintf("https://www.gushiwen.cn/shiwen2017/ajaxSearchSo.aspx?valuekey=%s", html.EscapeString(query)))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	r := make([]AlfredItem, 0)
	doc.Find(".main").Each(func(i int, cs *goquery.Selection) {
		category := cs.Find(".mleft > span").Text()
		cs.Find(".mright > div > a").Each(func(i int, s *goquery.Selection) {
			r = append(r, AlfredItem{
				Type:     "file",
				Title:    strings.ReplaceAll(s.Text(), "\n", ""),
				Subtitle: category,
				Arg:      s.AttrOr("href", ""),
				Icon: struct {
					Path string `json:"path"`
				}{
					Path: "icon.png",
				},
			})

		})
	})

	output, _ := json.Marshal(struct {
		Items []AlfredItem `json:"items"`
	}{
		Items: r,
	})

	fmt.Println(string(output))
}

func main() {
	getItems(os.Args[1])
}
