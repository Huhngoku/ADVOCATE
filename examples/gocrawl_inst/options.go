package main

import (
	"time"

	"github.com/PuerkitoBio/purell"
)

const (
	DefaultUserAgent          string                    = `Mozilla/5.0 (Windows NT 6.1; rv:15.0) gocrawl/0.4 Gecko/20120716 Firefox/15.0a2`
	DefaultRobotUserAgent     string                    = `Googlebot (gocrawl v0.4)`
	DefaultEnqueueChanBuffer  int                       = 100
	DefaultHostBufferFactor   int                       = 10
	DefaultCrawlDelay         time.Duration             = 5 * time.Second
	DefaultIdleTTL            time.Duration             = 10 * time.Second
	DefaultNormalizationFlags purell.NormalizationFlags = purell.FlagsAllGreedy
)

type Options struct {
	UserAgent string

	RobotUserAgent string

	MaxVisits int

	EnqueueChanBuffer int

	HostBufferFactor int

	CrawlDelay time.Duration

	WorkerIdleTTL time.Duration

	SameHostOnly bool

	HeadBeforeGet bool

	URLNormalizationFlags purell.NormalizationFlags

	LogFlags LogFlags

	Extender Extender
}

func NewOptions(ext Extender) *Options {

	return &Options{
		DefaultUserAgent,
		DefaultRobotUserAgent,
		0,
		DefaultEnqueueChanBuffer,
		DefaultHostBufferFactor,
		DefaultCrawlDelay,
		DefaultIdleTTL,
		true,
		false,
		DefaultNormalizationFlags,
		LogError,
		ext,
	}
}
