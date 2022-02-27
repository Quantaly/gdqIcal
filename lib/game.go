package lib

import (
	"fmt"
	"hash/crc32"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// Game represents a GDQ schedule block
type Game struct {
	StartTime time.Time
	Game      string
	Runners   string
	Duration  time.Duration
	Category  string
	Host      string
}

// WriteIcalEvent writes an iCalendar event for a game into w.
func (g *Game) WriteIcalEvent(w io.Writer) {
	startTime := g.StartTime.Format("20060102T150405Z07:00")
	endTime := (g.StartTime.Add(g.Duration)).Format("20060102T150405Z07:00")
	fmt.Fprintf(w, "BEGIN:VEVENT\r\n")
	fmt.Fprintf(w, "UID:%d@gamesdonequick.com\r\n", g.uid())
	fmt.Fprintf(w, "DTSTAMP:%s\r\n", startTime)
	fmt.Fprintf(w, "DTSTART:%s\r\n", startTime)
	fmt.Fprintf(w, "DTEND:%s\r\n", endTime)
	fmt.Fprintf(w, "SUMMARY:%s\r\n", g.Game)
	fmt.Fprintf(w, "DESCRIPTION:%s by %s\\nHosted by %s\r\n", g.Category, g.Runners, strings.TrimSpace(g.Host))
	fmt.Fprintf(w, "END:VEVENT\r\n")
}

func (g *Game) uid() uint32 {
	return crc32.ChecksumIEEE([]byte(g.Game + g.Runners + g.Category))
}

func ParseGame(firstRow *html.Node) (*Game, error) {
	secondRow := nextElement(firstRow)

	td := nextElement(firstRow.FirstChild)
	startTime, err := time.Parse(time.RFC3339, td.FirstChild.Data)
	if err != nil {
		return nil, err
	}

	td = nextElement(td)
	game := td.FirstChild.Data

	// they have had runs without runners in the table before
	// see Ninja Spirit @ SGDQ2019 (https://gamesdonequick.com/schedule/26)
	td = nextElement(td)
	var runners string
	if td.FirstChild != nil {
		runners = td.FirstChild.Data
	}

	// they have had runs with a setup length but no duration
	// see PreShow @ Frost Fatales 2022 (https://gamesdonequick.com/schedule/38)
	// the setup length seems to have been the intended duration in that case, so:
	// prefer the duration, but parse the setup length just in case
	// and then default to like 30 min? so it still shows up on google calendar
	td = nextElement(secondRow.FirstChild)
	var duration time.Duration
	if td.LastChild != nil {
		duration, err = parseGdqDuration(td.LastChild.Data)
	}
	if td.LastChild == nil || err != nil {
		setupTd := prevElement(firstRow.LastChild)
		duration, err = parseGdqDuration(setupTd.LastChild.Data)
		if err != nil {
			duration = time.Duration(30) * time.Minute
		}
	}

	td = nextElement(td)
	category := td.FirstChild.Data

	td = nextElement(td)
	host := td.LastChild.Data

	return &Game{startTime, game, runners, duration, category, host}, nil
}

var gdqDurationRegex = regexp.MustCompile(`(\d+):(\d+):(\d+)`)

func parseGdqDuration(s string) (time.Duration, error) {
	matches := gdqDurationRegex.FindStringSubmatch(s)
	hours, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, err
	}
	minutes, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, err
	}
	seconds, err := strconv.Atoi(matches[3])
	if err != nil {
		return 0, err
	}

	return time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second, nil
}
