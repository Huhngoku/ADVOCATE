package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/ErikKassubek/deadlockDetectorGo/src/dedego"
	"github.com/PuerkitoBio/goquery"
)

type DelayInfo struct {
	OptsDelay   time.Duration
	RobotsDelay time.Duration
	LastDelay   time.Duration
}

type FetchInfo struct {
	Ctx           *URLContext
	Duration      time.Duration
	StatusCode    int
	IsHeadRequest bool
}

type Extender interface {
	Start(interface{}) interface{}
	End(error)
	Error(*CrawlError)
	Log(LogFlags, LogFlags, string)

	ComputeDelay(string, *DelayInfo, *FetchInfo) time.Duration

	Fetch(*URLContext, string, bool) (*http.Response, error)
	RequestGet(*URLContext, *http.Response) bool
	RequestRobots(*URLContext, string) ([]byte, bool)
	FetchedRobots(*URLContext, *http.Response)
	Filter(*URLContext, bool) bool
	Enqueued(*URLContext)
	Visit(*URLContext, *http.Response, *goquery.Document) (interface{}, bool)
	Visited(*URLContext, interface{})
	Disallowed(*URLContext)
}

var HttpClient = &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {

	if isRobotsURL(req.URL) {
		if len(via) >= 10 {
			return errors.New("stopped after 10 redirects")
		}
		if len(via) > 0 {
			req.Header.Set("User-Agent", via[0].Header.Get("User-Agent"))
		}
		return nil
	}

	return ErrEnqueueRedirect
}}

type DefaultExtender struct {
	EnqueueChan dedego.Chan[interface{}]
}

func (de *DefaultExtender) Start(seeds interface{}) interface{} {
	return seeds
}

func (de *DefaultExtender) End(err error) {}

func (de *DefaultExtender) Error(err *CrawlError) {}

func (de *DefaultExtender) Log(logFlags LogFlags, msgLevel LogFlags, msg string) {
	if logFlags&msgLevel == msgLevel {
		log.Println(msg)
	}
}

func (de *DefaultExtender) ComputeDelay(host string, di *DelayInfo, lastFetch *FetchInfo) time.Duration {
	if di.RobotsDelay > 0 {
		return di.RobotsDelay
	}
	return di.OptsDelay
}

func (de *DefaultExtender) Fetch(ctx *URLContext, userAgent string, headRequest bool) (*http.Response, error) {
	var reqType string

	if headRequest {
		reqType = "HEAD"
	} else {
		reqType = "GET"
	}
	req, e := http.NewRequest(reqType, ctx.url.String(), nil)
	if e != nil {
		return nil, e
	}
	req.Header.Set("User-Agent", userAgent)
	return HttpClient.Do(req)
}

func (de *DefaultExtender) RequestGet(ctx *URLContext, headRes *http.Response) bool {
	return headRes.StatusCode >= 200 && headRes.StatusCode < 300
}

func (de *DefaultExtender) RequestRobots(ctx *URLContext, robotAgent string) (data []byte, doRequest bool) {
	return nil, true
}

func (de *DefaultExtender) FetchedRobots(ctx *URLContext, res *http.Response) {}

func (de *DefaultExtender) Filter(ctx *URLContext, isVisited bool) bool {
	return !isVisited
}

func (de *DefaultExtender) Enqueued(ctx *URLContext) {}

func (de *DefaultExtender) Visit(ctx *URLContext, res *http.Response, doc *goquery.Document) (harvested interface{}, findLinks bool) {
	return nil, true
}

func (de *DefaultExtender) Visited(ctx *URLContext, harvested interface{}) {}

func (de *DefaultExtender) Disallowed(ctx *URLContext) {}
