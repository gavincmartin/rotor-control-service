package schedule

import (
	"encoding/json"
	"time"
	"tutorials/rotor-controller/rotor"
)

// TrackingPass type stores States and Times
// should probably add an ID & SC here too
type TrackingPass struct {
	Times  []time.Time
	States []rotor.State
}

// FromJSON turns a JSON of time/az/el values into a TrackingPass struct
// input JSON is in the form:
// {
//   "times": [1539227885, 1539227890, 1539227895, 1539227900],
//   "azimuths": [10.0, 11.0, 12.0, 13.0],
//   "elevations": [5.0, 5.0, 5.0, 5.0]
// }
func FromJSON(data []byte) TrackingPass {
	var m map[string][]interface{}
	err := json.Unmarshal(data, &m)
	if err != nil {
		panic(err)
	}
	times := m["times"]
	azimuths := m["azimuths"]
	elevations := m["elevations"]

	t := make([]time.Time, len(times))
	s := make([]rotor.State, len(times))
	for i := 0; i < len(times); i++ {
		// TODO:Change time implementation for nanoseconds

		// type assertions
		tI64 := int64(times[i].(float64))
		azF64 := azimuths[i].(float64)
		elF64 := elevations[i].(float64)

		t[i] = time.Unix(tI64, 0).UTC()
		s[i] = rotor.State{Az: azF64, El: elF64}
	}
	return TrackingPass{Times: t, States: s}
}
