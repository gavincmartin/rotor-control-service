package schedule

import (
	"fmt"
)

// Schedule stores a list of TrackingPass objects
type Schedule []TrackingPass

// Len returns the length of the schedule
func (s Schedule) Len() int {
	return len(s)
}

// Less returns true if one TrackingPass's start time is before the other's
func (s Schedule) Less(i, j int) bool {
	return s[i].Times[0].Before(s[j].Times[0])
}

// Swap switches adjacent elements of a slice
func (s Schedule) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// TODO: fix
func (s Schedule) String() string {
	return fmt.Sprintf("")
}
