package rotor

import (
	"encoding/json"
	"math"
	"sync"
	"time"
)

// Rotor type that stores the current state and rotates
type Rotor struct {
	mu sync.RWMutex
	State
}

// State type that stores an azimuth and elevation
type State struct {
	Az float64 `json:"azimuth" bson:"azimuth"`
	El float64 `json:"elevation" bson:"elevation"`
}

// StateFromJSON used for unmarshalling of the State type
func StateFromJSON(data []byte) State {
	s := State{}
	err := json.Unmarshal(data, &s)
	if err != nil {
		panic(err)
	}
	return s
}

// ToJSON used for marshalling the Rotor type (in a concurrency-safe way)
func (r *Rotor) ToJSON() []byte {
	r.mu.RLock()
	defer r.mu.RUnlock()
	jsonData, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return jsonData
}

func (r *Rotor) rotate(s State) {
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

// Rotate used for rotating the Rotor to a desired state
func (r *Rotor) Rotate(s State) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rotate(s)
}

// GetAz is a concurrency-safe retreival method for the current azimuth of the Rotor
func (r *Rotor) GetAz() float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.Az
}

// GetEl is a concurrency-safe retreival method for the current elevation of the Rotor
func (r *Rotor) GetEl() float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.El
}
