package lib

import (
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// Game represents a GDQ schedule block.
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

// ParseGame parses a game my dudes
func ParseGame(z *html.Tokenizer) (*Game, error) {
	nextToken(z) // text: newline
	nextToken(z) // <td.start-time.text-right>
	nextToken(z) // text: start time
	startTime, err := time.Parse(time.RFC3339, string(z.Text()))
	if err != nil {
		return nil, err
	}
	nextToken(z) // </td>
	nextToken(z) // text: newline
	nextToken(z) // <td>
	nextToken(z) // text: game
	game := string(z.Text())
	nextToken(z) // </td>
	nextToken(z) // text: newline
	nextToken(z) // <td>
	runners := ""
	if nextToken(z) == html.TextToken { // </td> or text: runners
		// so apparently there are runs without runners, i.e. Ninja Spirit @ SGDQ2019
		runners = string(z.Text())
		nextToken(z) // </td>
	}
	nextToken(z)                          // text: newline
	nextToken(z)                          // <td>
	if nextToken(z) != html.EndTagToken { // </td> or text: space
		// some runs also don't have setup times, i.e. Daily Recap - Monday @ SGDQ2021
		// one of these days I'm going to rewrite this POS with an actual DOM library
		nextToken(z) // <i.fa.fa-clock-o.text-gdq-red>
		nextToken(z) // </i>
		nextToken(z) // text: setup length
		nextToken(z) // </td>
	}
	// nextToken(z) // text: space
	// nextToken(z) // <i.fa.fa-clock-o.text-gdq-red>
	// nextToken(z) // </i>
	// nextToken(z) // text: setup length

	//nextToken(z) // </td>
	nextToken(z) // text: newline
	nextToken(z) // </tr>
	nextToken(z) // text: newline
	nextToken(z) // <tr.second-row>
	nextToken(z) // text: newline
	nextToken(z) // <td.text-right>
	nextToken(z) // text: space
	nextToken(z) // <i.fa.fa-clock-o>
	nextToken(z) // </i>
	nextToken(z) // text: duration
	duration, err := parseGdqDuration(string(z.Text()))
	if err != nil {
		return nil, err
	}
	nextToken(z) // </td>
	nextToken(z) // text: newline
	nextToken(z) // <td>
	nextToken(z) // text: category
	category := string(z.Text())
	nextToken(z) // </td>
	nextToken(z) // text: newline
	nextToken(z) // <td>
	nextToken(z) // <i.fa.fa-microphone>
	nextToken(z) // </i>
	nextToken(z) // text: host
	host := string(z.Text())
	nextToken(z) // </td>
	nextToken(z) // text: newline
	nextToken(z) // </tr>
	nextToken(z) // text: newline
	// next token is <tr> or </tbody>

	return &Game{startTime, game, runners, duration, category, host}, nil
}

// advance to the next token, and print it if debugging
func nextToken(z *html.Tokenizer) html.TokenType {
	const DEBUGGING = false
	ret := z.Next()
	if DEBUGGING {
		log.Println(ret)
		log.Println(string(z.Raw()))
	}
	return ret
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

	// lol unicode variable names
	const μsSecond = 1000000000
	const μsMinute = μsSecond * 60
	const μsHour = μsMinute * 60

	return time.Duration(hours*μsHour + minutes*μsMinute + seconds*μsSecond), nil
}
