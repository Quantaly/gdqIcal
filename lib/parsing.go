package lib

import (
	"fmt"
	"golang.org/x/net/html"
	"hash/crc32"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
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
	z.Next() // text: newline
	z.Next() // <td.start-time.text-right>
	z.Next() // text: start time
	startTime, err := time.Parse(time.RFC3339, string(z.Text()))
	if err != nil {
		return nil, err
	}
	z.Next() // </td>
	z.Next() // text: newline
	z.Next() // <td>
	z.Next() // text: game
	game := string(z.Text())
	z.Next() // </td>
	z.Next() // text: newline
	z.Next() // <td>
	runners := ""
	if z.Next() == html.TextToken {
		// so apparently there are runs without runners, i.e. Ninja Spirit @ SGDQ2019
		runners = string(z.Text())
		z.Next() // </td>
	}
	z.Next() // text: newline
	z.Next() // <td>
	z.Next() // text: space
	z.Next() // <i.fa.fa-clock-o.text-gdq-red>
	z.Next() // </i>
	z.Next() // text: setup length
	z.Next() // </td>
	z.Next() // text: newline
	z.Next() // </tr>
	z.Next() // text: newline
	z.Next() // <tr.second-row>
	z.Next() // text: newline
	z.Next() // <td.text-right>
	z.Next() // text: space
	z.Next() // <i.fa.fa-clock-o>
	z.Next() // </i>
	z.Next() // text: duration
	duration, err := parseGdqDuration(string(z.Text()))
	if err != nil {
		return nil, err
	}
	z.Next() // </td>
	z.Next() // text: newline
	z.Next() // <td>
	z.Next() // text: category
	category := string(z.Text())
	z.Next() // </td>
	z.Next() // text: newline
	z.Next() // <td>
	z.Next() // <i>
	z.Next() // </i>
	z.Next() // text: host
	host := string(z.Text())
	z.Next() // </td>
	z.Next() // text: newline
	z.Next() // </tr>
	z.Next() // text: newline
	// next token is <tr> or </tbody>

	return &Game{startTime, game, runners, duration, category, host}, nil
}

func debugNextToken(z *html.Tokenizer) {
	fmt.Println(z.Next())
	fmt.Println(string(z.Raw()))
}

var gdqDurationRegex = regexp.MustCompile("(\\d+):(\\d+):(\\d+)")

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
