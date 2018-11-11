package integrations

import (
	"bytes"
	"encoding/json"
	"net/http"
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

func formatSchedule(schedule []passes.TrackingPass) []byte {
	attachments := make([]attachment, len(schedule))
	for i, pass := range schedule {
		attachments[i] = passToAttachment(pass)
	}
	payload := slackPayload{Text: "Here's today's tracking schedule! :satellite_antenna: (sent from Go)", Attachments: attachments}
	return payload.ToJSON()
}

type slackPayload struct {
	Text        string       `json:"text"`
	Attachments []attachment `json:"attachments"`
}

func (p slackPayload) ToJSON() []byte {
	jsonData, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return jsonData
}

type attachment struct {
	Fields     []field `json:"fields"`
	AuthorName string  `json:"author_name"`
}

func passToAttachment(pass passes.TrackingPass) attachment {
	startField := field{Title: "start_time", Value: pass.StartTime.In(loc).Format(time.Kitchen), Short: true}
	endField := field{Title: "end_time", Value: pass.Times[len(pass.Times)-1].In(loc).Format(time.Kitchen), Short: true}
	fields := []field{startField, endField}
	return attachment{Fields: fields, AuthorName: pass.Spacecraft}
}

type field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}
