package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)


func main() {

	logger := log.New(os.Stderr, "ERROR: ", 0)

	if len(os.Args) == 2 {
		t, err := parseArg(os.Args[1])

		if err != nil {
			logger.Fatal(err)
		}

		now := time.Now()

		mjt := diffTime(now, t)
		fmt.Println(mjt.String())
	}
}

func parseArg(s string) (time.Time, error) {

	var t time.Time
	var err error

	// Try *ALL* the formats until we get a hit!
	allConstants := []string{
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
	}

	for _, f := range allConstants {
		t, err = time.Parse(f, s)
		if err == nil {
			return t, nil
		}
	}
	return t, errors.New("could not find a suitable way to parse the date")
}

type TimeValues struct {
	year int
	month int
	day int
	hour int
	minute int
	second int
}

func (c *TimeValues) String() string {
	return fmt.Sprintf( "%d year(s), %d month(s), %d day(s), %d hour(s), %d minute(s), %d second(s)", c.year, c.month, c.day, c.hour, c.minute, c.second)
}

// diffTime -- logic inspired by https://stackoverflow.com/questions/36530251/time-since-with-months-and-years
func diffTime(a, b time.Time) TimeValues {

	var r TimeValues

	// Normalize TZ
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}

	// Make sure we're in correct order
	if a.After(b) {
		a, b = b, a
	}

	// Get YMD values
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	// Get HMS values
	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	r.year = y2 - y1
	r.month = int(M2 - M1)
	r.day = d2 - d1
	r.hour = h2 - h1
	r.minute = m2 - m1
	r.second = s2 - s1

	// Normalize negative values
	if r.second < 0 {
		r.second += 60
		r.minute--
	}
	if r.minute < 0 {
		r.minute += 60
		r.hour--
	}
	if r.hour < 0 {
		r.hour += 24
		r.day--
	}
	if r.day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		r.day += 32 - t.Day()
		r.month--
	}
	if r.month < 0 {
		r.month += 12
		r.year--
	}

	return r
}