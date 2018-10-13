package schedule

import (
	"log"
	"sort"
	"strconv"
	"strings"

	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type PassesDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "passes"
)

func (p *PassesDAO) Connect() {
	session, err := mgo.Dial(p.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(p.Database)
}

func (p *PassesDAO) FindAll() ([]TrackingPass, error) {
	var passes []TrackingPass
	err := db.C(COLLECTION).Find(bson.M{}).All(&passes)
	return passes, err
}

func (p *PassesDAO) Insert(pass TrackingPass) error {
	err := db.C(COLLECTION).Insert(&pass)
	return err
}

////////////////////////

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

// TODO: make prettier print later
func (s Schedule) String() string {
	var table strings.Builder
	for i, pass := range s {
		table.WriteString(strconv.Itoa(i) + ": " + pass.String())
		if i == len(s)-1 {
			continue
		}
		table.WriteRune('\n')
	}
	return table.String()
}

// AddPass adds a TrackingPass struct to a Schedule struct by performing a
// binary search and inserting into the Schedule by start time
func (s *Schedule) AddPass(t TrackingPass) {
	// *s = append(*s, t)
	// sort.Sort(*s)
	idx := sort.Search(len(*s), func(i int) bool { return (*s)[i].Times[0].After(t.Times[0]) })
	*s = append(*s, TrackingPass{})
	copy((*s)[idx+1:], (*s)[idx:])
	(*s)[idx] = t
}
