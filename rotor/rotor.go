package rotor

import (
	"encoding/json"
	"math"
	"time"
)

// Rotor type that stores the current state and rotates
type Rotor struct {
	State
}

// State type that stores an azimuth and elevation
type State struct {
	Az float64 `json:"azimuth" bson:"azimuth"`
	El float64 `json:"elevation" bson:"elevation"`
}

// FromJSON used for unmarshalling of the State type
func StateFromJSON(data []byte) State {
	s := State{}
	err := json.Unmarshal(data, &s)
	if err != nil {
		panic(err)
	}
	return s
}

// Rotate used for rotating the Rotor to a desired state
func (r *Rotor) Rotate(s State) {
	// TODO: need to add actual rotation stuff here

	deltaAz := math.Copysign(0.1, s.Az-r.Az)
	for a := r.Az; a <= s.Az; a += deltaAz {
		time.Sleep(10 * time.Millisecond)
		r.Az = a
	}

	deltaEl := math.Copysign(0.1, s.El-r.El)
	for e := r.El; e <= s.El; e += deltaEl {
		time.Sleep(10 * time.Millisecond)
		r.El = e
	}
}
