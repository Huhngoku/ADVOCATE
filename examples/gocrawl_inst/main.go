package main

import (
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ErikKassubek/deadlockDetectorGo/src/dedego"
	"github.com/PuerkitoBio/goquery"
)

var rxOk = regexp.MustCompile(`http://duckduckgo\.com(/a.*)?$`)

type ExampleExtender struct {
	DefaultExtender
}

func (x *ExampleExtender) Visit(ctx *URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {

	return nil, true
}

func (x *ExampleExtender) Filter(ctx *URLContext, isVisited bool) bool {
	return !isVisited && rxOk.MatchString(ctx.NormalizedURL().String())
}

func main() {
	var order string
	if len(os.Args) > 0 {
		order = os.Args[1]
	}
	order_split := strings.Split(order, ";")
	for _, ord := range order_split {
		ord_split := strings.Split(ord, ",")
		id, err1 := strconv.Atoi(ord_split[0])
		c, err2 := strconv.Atoi(ord_split[1])
		if err1 == nil && err2 == nil {
			dedegoFetchOrder[id] = c
		}
	}
	dedego.Init(20)
	defer dedego.RunAnalyzer()
	defer time.Sleep(time.Millisecond)

	opts := NewOptions(new(ExampleExtender))

	opts.RobotUserAgent = "Example"

	opts.UserAgent = "Mozilla/5.0 (compatible; Example/1.0; +http://example.com)"

	opts.CrawlDelay = 1 * time.Second
	opts.LogFlags = LogAll

	opts.MaxVisits = 2

	c := NewCrawlerWithOptions(opts)
	c.Run("https://duckduckgo.com/")

}

var dedegoFetchOrder = make(map[int]int)
