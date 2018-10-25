package executor

import (
	"time"
	"tutorials/rotor-controller/passes"
	"tutorials/rotor-controller/rotor"
)

type Executor struct {
	rotctl   rotor.Rotor
	db       passes.DAO
	updates  <-chan struct{}
	nextPass passes.TrackingPass
}

func (e *Executor) Run() {
	for {
		select {
		// there was an update
		case <-e.updates:
			e.nextPass, _ = e.db.GetNextPass()
		// no update
		default:
			// if the next TrackingPass starts within 1 minute
			if time.Until(e.nextPass.StartTime) <= 1*time.Minute {
				// go e.TrackPass(e.nextPass)
			} else {
				time.Sleep(3 * time.Second)
			}
		}
	}
}

// TODO: add channel to communicate when done? or to abort/broadcast? (8.9)
func (e *Executor) TrackPass(pass passes.TrackingPass, done <-chan struct{}) error {
	select {
	case <-done:
		return nil
	}
}
