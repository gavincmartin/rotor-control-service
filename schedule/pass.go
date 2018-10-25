package schedule

import (
	"encoding/json"
	"fmt"
	"github.com/gavincmartin/rotor-control-service/rotor"
	"github.com/globalsign/mgo/bson"
	"time"
)

// TrackingPass type stores a spacecraft, ID, states and times for a given
// communications pass
type TrackingPass struct {
	ID         bson.ObjectId
	Spacecraft string
	Times      []time.Time
	States     []rotor.State
}

func (t TrackingPass) String() string {
	return fmt.Sprintf("S/C: %v | Start: %v | End %v", t.Spacecraft, t.Times[0], t.Times[len(t.Times)-1])
}

// FromJSON turns a JSON of time/az/el values into a TrackingPass struct
// input JSON is in the form:
// {
//   "times": [1539227885, 1539227890, 1539227895, 1539227900],
//   "azimuths": [10.0, 11.0, 12.0, 13.0],
//   "elevations": [5.0, 5.0, 5.0, 5.0]
// }
func FromJSON(data []byte) TrackingPass {
	p := PassInfo{}
	err := json.Unmarshal(data, &p)
	if err != nil {
		panic(err)
	}
	return p.ToTrackingPass()
}

func (t TrackingPass) toPassInfo() PassInfo {
	times := make([]float64, len(t.Times))
	azimuths := make([]float64, len(t.Times))
	elevations := make([]float64, len(t.Times))
	for i := range t.Times {
		// TODO: add nanoseconds here
		times[i] = float64(t.Times[i].Unix())
		azimuths[i] = t.States[i].Az
		elevations[i] = t.States[i].El
	}
	p := PassInfo{ID: t.ID, Spacecraft: t.Spacecraft, Times: times, Azimuths: azimuths, Elevations: elevations}
	return p
}

// PassInfo stores information related to a TrackingPass in a JSON-friendly
// format (rather than a Go-friendly format). It is used to change between the two.
type PassInfo struct {
	ID         bson.ObjectId
	Spacecraft string    `json:"spacecraft" bson:"spacecraft"`
	Times      []float64 `json:"times" bson:"times"`
	Azimuths   []float64 `json:"azimuths" bson:"azimuths"`
	Elevations []float64 `json:"elevations" bson:"elevations"`
}

// ToTrackingPass converts a PassInfo struct into a TrackingPass struct
func (p PassInfo) ToTrackingPass() TrackingPass {
	t := make([]time.Time, len(p.Times))
	s := make([]rotor.State, len(p.Times))
	for i := range p.Times {
		// TODO: Change time implementation for nanoseconds

		t[i] = time.Unix(int64(p.Times[i]), 0).UTC()
		s[i] = rotor.State{Az: p.Azimuths[i], El: p.Elevations[i]}
	}
	return TrackingPass{Spacecraft: p.Spacecraft, Times: t, States: s}
}
