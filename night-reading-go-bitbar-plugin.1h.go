package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/johnmccabe/go-bitbar"
)

var (
	// label: 预分享
	preShareLabel = "https://github.com/developer-learning/night-reading-go/labels/%E9%A2%84%E5%88%86%E4%BA%AB"

	// label: 已排期
	planedLabel = "https://github.com/developer-learning/night-reading-go/labels/%E5%B7%B2%E6%8E%92%E6%9C%9F"

	// label: 已分享
	sharedLabel = "https://github.com/developer-learning/night-reading-go/issues?q=label%3A%E5%B7%B2%E5%88%86%E4%BA%AB+is%3Aclosed"
)

// readingShare
type readingShare struct {
	title string
	link  string
}

func main() {
	app := bitbar.New()
	submenu := app.NewSubMenu()

	preShares := fetchLabel(preShareLabel)
	planeds := fetchLabel(planedLabel)
	shareds := fetchLabel(sharedLabel)

	total := fmt.Sprintf("预分享：%d， 已排期：%d，已分享：%d", len(preShares), len(planeds), len(shareds))
	nightReadingGo := " Go 夜读 - " + total
	app.StatusLine(nightReadingGo)

	submenu.Line("Refresh...").Color("black").Refresh()

	planed := fmt.Sprintf("已排期（%d)", len(planeds))
	submenu.Line(planed).Color("#c115a7")
	for _, rs := range planeds {
		subsubmenu := submenu.NewSubMenu()
		subsubmenu.Line(rs.title).Href(rs.link).Color("black")
	}

	preShare := fmt.Sprintf("预分享（%d）", len(preShares))
	submenu.Line(preShare).Color("#f2f096")

	for _, rs := range preShares {
		subsubmenu := submenu.NewSubMenu()
		subsubmenu.Line(rs.title).Href(rs.link).Color("black")
	}

	splitLine := "----------------"
	submenu.Line(splitLine).Color("black")

	shared := fmt.Sprintf("已分享（%d）", len(shareds))
	submenu.Line(shared).Color("#f7b9f1")

	for _, rs := range shareds {
		subsubmenu := submenu.NewSubMenu()
		subsubmenu.Line(rs.title).Href(rs.link).Color("black")
	}

	app.Render()
}

func fetchLabel(urlStr string) []readingShare {
	var rs []readingShare
	dom := fetchURLHTML(urlStr)
	dom.Find(".link-gray-dark.v-align-middle.no-underline.h4.js-navigation-open").Each(func(i int, selection *goquery.Selection) {
		href, ok := selection.Attr("href")
		if !ok {
			fmt.Println("get issues link error")
			return
		}
		href = "https://github.com" + href
		title := selection.Text()
		rs = append(rs, readingShare{title: title, link: href})
	})
	return rs
}

func fetchURLHTML(urlStr string) *goquery.Document {
	resp, err := http.Get(urlStr)
	if err != nil {
		fmt.Println("http get error", err)
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read error", err)
		return nil
	}
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		log.Fatalln(err)
	}
	return dom
}
