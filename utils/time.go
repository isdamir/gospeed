package utils

import (
	"strconv"
	"strings"
	"time"
)

type Time struct {
	Layout string
	Zone   string
}

func WebTime(t time.Time) string {
	ftime := t.Format(time.RFC1123)
	if strings.HasSuffix(ftime, "UTC") {
		ftime = ftime[0:len(ftime)-3] + "GMT"
	}
	return ftime
}

func NewTime() *Time {
	return &Time{
		Layout: "2006-01-02 15:04:05",
		Zone:   "-07:00",
	}
}

func (t *Time) UnixToStr(n int64, layouts ...string) string {
	layout := t.Layout
	if len(layouts) > 0 && layouts[0] != "" {
		layout = layouts[0]
	}

	var ss, ns string
	s := strconv.FormatInt(n, 10)
	l := len(s)
	if l > 10 {
		ss = s[:10]
		ns = s[10:]
		//000000000
		fillLen := 9 - len(ns)
		ns = ns + strings.Repeat("0", fillLen)
	} else {
		ss = s[:10]
	}

	si, _ := strconv.ParseInt(ss, 10, 64)
	ni, _ := strconv.ParseInt(ns, 10, 64)
	tm := time.Unix(si, ni)

	return tm.Format(layout)
}

func (t *Time) StrToUnix(s string) int64 {
	ti, _ := time.Parse(t.Layout+" "+t.Zone, s)
	return ti.Unix()
}

// DateFormat takes a time and a layout string and returns a string with the formatted date. Used by the template parser as "dateformat"
func DateFormat(t time.Time, layout string) (datestring string) {
	datestring = t.Format(layout)
	return
}

// Date takes a PHP like date func to Go's time fomate
func Date(t time.Time, format string) (datestring string) {
	patterns := []string{
		// year
		"Y", "2006", // A full numeric representation of a year, 4 digits	Examples: 1999 or 2003
		"y", "06", //A two digit representation of a year	Examples: 99 or 03

		// month
		"m", "01", // Numeric representation of a month, with leading zeros	01 through 12
		"n", "1", // Numeric representation of a month, without leading zeros	1 through 12
		"M", "Jan", // A short textual representation of a month, three letters	Jan through Dec
		"F", "January", // A full textual representation of a month, such as January or March	January through December

		// day
		"d", "02", // Day of the month, 2 digits with leading zeros	01 to 31
		"j", "2", // Day of the month without leading zeros	1 to 31

		// week
		"D", "Mon", // A textual representation of a day, three letters	Mon through Sun
		"l", "Monday", // A full textual representation of the day of the week	Sunday through Saturday

		// time
		"g", "3", // 12-hour format of an hour without leading zeros	1 through 12
		"G", "15", // 24-hour format of an hour without leading zeros	0 through 23
		"h", "03", // 12-hour format of an hour with leading zeros	01 through 12
		"H", "15", // 24-hour format of an hour with leading zeros	00 through 23

		"a", "pm", // Lowercase Ante meridiem and Post meridiem	am or pm
		"A", "PM", // Uppercase Ante meridiem and Post meridiem	AM or PM

		"i", "04", // Minutes with leading zeros	00 to 59
		"s", "05", // Seconds, with leading zeros	00 through 59
	}
	replacer := strings.NewReplacer(patterns...)
	format = replacer.Replace(format)
	datestring = t.Format(format)
	return
}
