package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/ErikKassubek/deadlockDetectorGo/src/dedego"
)

const version = "1.0.2"

var onlyPrintVersion = flag.Bool("version", false, "print the htcat version")

const (
	_        = iota
	KB int64 = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
)

func printUsage() {
	log.Printf("usage: %v URL", os.Args[0])
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

	url_a := "https://download.geonames.org/export/dump/cities1000.zip"
	u, err := url.Parse(url_a)
	if err != nil {
		log.Fatalf("aborting: could not parse given URL: %v", err)
	}

	client := *http.DefaultClient

	switch u.Scheme {
	case "https":
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{},
		}
	case "http":
	default:

		printUsage()
		log.Fatalf("aborting: unsupported URL scheme %v", u.Scheme)
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	htc := New(&client, u, 5)

	if _, err := htc.WriteTo(os.Stdout); err != nil {
		log.Fatalf("aborting: could not write to output stream: %v",
			err)
	}
}

var dedegoFetchOrder = make(map[int]int)
