package passes

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
	ID         bson.ObjectId `json:"id" bson:"_id"`
	Spacecraft string        `json:"spacecraft" bson:"spacecraft"`
	Times      []time.Time   `json:"times" bson:"times"`
	States     []rotor.State `json:"states" bson:"states"`
	StartTime  time.Time     `json:"start_time" bson:"start_time"`
}

func (t TrackingPass) String() string {
	return fmt.Sprintf("S/C: %v | Start: %v | End %v | ID: %v", t.Spacecraft, t.Times[0], t.Times[len(t.Times)-1], t.ID.Hex())
}

// ToJSON marhsals a TrackingPass struct into JSON format
func (t TrackingPass) ToJSON() []byte {
	jsonData, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return jsonData
}

// FromJSON unmarshals a TrackingPass struct from JSON input in the form:
// {
//     "spacecraft": "ARMADILLO",
//     "times": [
//         "2018-10-11T03:18:05Z",
//         "2018-10-11T03:18:10Z",
//         "2018-10-11T03:18:15Z",
//         "2018-10-11T03:18:20Z"
//     ],
//     "states": [
//         {
//             "azimuth": 10,
//             "elevation": 5
//         },
//         {
//             "azimuth": 11,
//             "elevation": 5
//         },
//         {
//             "azimuth": 12,
//             "elevation": 5
//         },
//         {
//             "azimuth": 13,
//             "elevation": 5
//         }
//     ]
// }
func FromJSON(data []byte) TrackingPass {
	t := TrackingPass{}
	err := json.Unmarshal(data, &t)
	if err != nil {
		panic(err)
	}
	t.StartTime = t.Times[0]
	return t
}
