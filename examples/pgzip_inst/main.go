// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/ErikKassubek/deadlockDetectorGo/src/dedego"
)

type City struct {
	GeonameID      int
	Name           string
	AsciiName      string
	AlternateNames []string
	Latitude       float64
	Longitude      float64
	FeatureClass   string
	FeatureCode    string
	CountryCode    string
	CC2            string

	Population int
	Elevation  int
	Dem        string
	Timezone   string
	Modified   time.Time
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
	dedego.Init(300)
	defer dedego.RunAnalyzer()
	defer time.Sleep(time.Millisecond)
	// defer dedego.PrintTrace()
	runtime.GOMAXPROCS(runtime.NumCPU())
	file, err := os.Open("allCountries.txt.gz")
	if err != nil {
		panic(err)
	}
	r, err := NewReader(file)
	if err != nil {
		panic(err)
	}
	scan := bufio.NewScanner(r)
	t := time.Now()

	n := 0
	for scan.Scan() {
		line := scan.Text()
		s := strings.Split(line, "\t")
		if len(s) < 19 {
			continue
		}
		c := City{}
		c.GeonameID, _ = strconv.Atoi(s[0])
		c.Name = s[1]
		c.AsciiName = s[2]
		c.AlternateNames = strings.Split(s[3], ",")
		c.Latitude, _ = strconv.ParseFloat(s[4], 64)
		c.Longitude, _ = strconv.ParseFloat(s[5], 64)
		c.FeatureClass = s[6]
		c.FeatureCode = s[7]
		c.CountryCode = s[8]
		c.CC2 = s[9]
		c.Population, _ = strconv.Atoi(s[14])
		c.Elevation, _ = strconv.Atoi(s[15])
		c.Dem = s[16]
		c.Timezone = s[17]
		c.Modified, _ = time.Parse("2006-01-02", s[18])

		n++
	}
	d := time.Since(t)
	fmt.Printf("Processed %d entries in %v, %.1f entries/sec.", n, d, float64(n)/(float64(d)/float64(time.Second)))
}

var dedegoFetchOrder = make(map[int]int)
