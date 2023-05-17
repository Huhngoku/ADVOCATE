package main

import (
	"errors"
)

var (
	ErrEnqueueRedirect = errors.New("redirection not followed")

	ErrMaxVisits = errors.New("the maximum number of visits is reached")

	ErrInterrupted = errors.New("interrupted")
)

type CrawlErrorKind uint8

const (
	CekFetch CrawlErrorKind = iota
	CekParseRobots
	CekHttpStatusCode
	CekReadBody
	CekParseBody
	CekParseURL
	CekProcessLinks
	CekParseRedirectURL
)

var (
	lookupCek = [...]string{
		CekFetch:            "Fetch",
		CekParseRobots:      "ParseRobots",
		CekHttpStatusCode:   "HttpStatusCode",
		CekReadBody:         "ReadBody",
		CekParseBody:        "ParseBody",
		CekParseURL:         "ParseURL",
		CekProcessLinks:     "ProcessLinks",
		CekParseRedirectURL: "ParseRedirectURL",
	}
)

func (cek CrawlErrorKind) String() string {
	return lookupCek[cek]
}

type CrawlError struct {
	Ctx *URLContext

	Err error

	Kind CrawlErrorKind

	msg string
}

func (ce CrawlError) Error() string {
	if ce.Err != nil {
		return ce.Err.Error()
	}
	return ce.msg
}

func newCrawlError(ctx *URLContext, e error, kind CrawlErrorKind) *CrawlError {
	return &CrawlError{ctx, e, kind, ""}
}

func newCrawlErrorMessage(ctx *URLContext, msg string, kind CrawlErrorKind) *CrawlError {
	return &CrawlError{ctx, nil, kind, msg}
}
