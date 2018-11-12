package integrations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gavincmartin/rotor-control-service/passes"
	"github.com/spf13/viper"
)

var loc, _ = time.LoadLocation("America/Chicago")

// SendSlackSchedule POSTs a slice of TrackingPass structs to a specified
// slack URL in a schedule format
func SendSlackSchedule(schedule []passes.TrackingPass) {
	slackPOSTUrl := viper.GetString("SlackPOSTUrl")
	if len(slackPOSTUrl) == 0 {
		return
	}
	payload := formatSchedule(schedule)
	resp, err := http.Post(slackPOSTUrl, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

// SendSlackPass POSTs a TrackingPass struct to a specified slack URL
func SendSlackPass(pass passes.TrackingPass) {
	slackPOSTUrl := viper.GetString("SlackPOSTUrl")
	if len(slackPOSTUrl) == 0 {
		return
	}
	payload := formatSinglePass(pass)
	resp, err := http.Post(slackPOSTUrl, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func formatSchedule(schedule []passes.TrackingPass) []byte {
	attachments := make([]attachment, len(schedule))
	for i, pass := range schedule {
		attachments[i] = passToAttachment(pass)
	}
	payload := slackPayload{Text: "Here's today's tracking schedule! :satellite_antenna:", Attachments: attachments}
	return payload.ToJSON()
}

func formatSinglePass(pass passes.TrackingPass) []byte {
	attachments := make([]attachment, 1)
	attachments[0] = passToAttachment(pass)
	payload := slackPayload{Text: "A pass is about to start! :satellite:", Attachments: attachments}
	return payload.ToJSON()
}

type slackPayload struct {
	Text        string       `json:"text"`
	Attachments []attachment `json:"attachments"`
}

func (p slackPayload) ToJSON() []byte {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(p)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

type attachment struct {
	Fields     []field `json:"fields"`
	AuthorName string  `json:"author_name"`
}

func passToAttachment(pass passes.TrackingPass) attachment {
	endTime := pass.Times[len(pass.Times)-1]
	calendarLink := generateCalendarLink(pass.Spacecraft, pass.StartTime, endTime)
	startField := field{Title: "start_time", Value: timeToSlackFormat(pass.StartTime, calendarLink), Short: true}
	endField := field{Title: "end_time", Value: timeToSlackFormat(endTime, calendarLink), Short: true}
	fields := []field{startField, endField}
	return attachment{Fields: fields, AuthorName: pass.Spacecraft}
}

func timeToSlackFormat(t time.Time, calendarLink string) string {
	//ex: "<!date^1542004812^{date_short_pretty} at {time}^https://calendar.google.com/calendar/r/eventedit?text=My+Custom+Event&dates=20180512T230000Z/20180513T030000Z|8:39 AM (CT)>"
	return fmt.Sprintf("<!date^%v^{date_short_pretty} at {time}^%v|%v (CT)>", t.Unix(), calendarLink, t.In(loc).Format(time.Kitchen))
}

func generateCalendarLink(spacecraft string, start, end time.Time) string {
	title := strings.Replace(spacecraft+" Tracking Pass", " ", "+", -1)
	timeFormat := "20060102T150405Z"
	link := fmt.Sprintf("https://calendar.google.com/calendar/r/eventedit?text=%v&dates=%v/%v", title, start.Format(timeFormat), end.Format(timeFormat))
	return link
}

type field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}
