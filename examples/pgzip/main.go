package main

import (
	"bufio"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"advocate"
)

type City struct {
	GeonameID      int
	Name           string   //  name of geographical point (utf8) varchar(200)
	AsciiName      string   //  name of geographical point in plain ascii characters, varchar(200)
	AlternateNames []string //  alternatenames, comma separated, ascii names automatically transliterated, convenience attribute from alternatename table, varchar(10000)
	Latitude       float64  //  latitude in decimal degrees (wgs84)
	Longitude      float64  //  longitude in decimal degrees (wgs84)
	FeatureClass   string   //  see httpstring // //www.geonames.org/export/codes.html, char(1)
	FeatureCode    string   //  see httpstring // //www.geonames.org/export/codes.html, varchar(10)
	CountryCode    string   //  ISO-3166 2-letter country code, 2 characters
	CC2            string   //  alternate country codes, comma separated, ISO-3166 2-letter country code, 60 characters
	/* Skip 4 Admin */
	Population int       //  bigint (8 byte int)
	Elevation  int       //  in meters, integer
	Dem        string    //  digital elevation model, srtm3 or gtopo30, average elevation of 3''x3'' (ca 90mx90m) or 30''x30'' (ca 900mx900m) area in meters, integer. srtm processed by cgiar/ciat.
	Timezone   string    //  the timezone id (see file timeZone.txt) varchar(40)
	Modified   time.Time //  date of last modification in yyyy-MM-dd format
}

func main() {
	if false {
		// init tracing
		advocate.InitTracing(0)
		defer advocate.Finish()
	} else {
		// init replay
		advocate.EnableReplay()
		defer advocate.WaitForReplayFinish()
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	file, err := os.Open("./allCountries.txt.gz")
	if err != nil {
		panic(err)
	}
	r, err := NewReader(file)
	if err != nil {
		panic(err)
	}
	scan := bufio.NewScanner(r)
	// t := time.Now()

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
		//// fmt.Printf("%#v\n", c)
		n++
	}
	// d := time.Since(t)
	// fmt.Printf("Processed %d entries in %v, %.1f entries/sec.\n", n, d, float64(n)/(float64(d)/float64(time.Second)))
}
