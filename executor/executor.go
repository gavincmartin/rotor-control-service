package executor

import (
	"github.com/gavincmartin/rotor-control-service/passes"
	"github.com/gavincmartin/rotor-control-service/rotor"
	"math"
	"time"
)

// Executor stores the relevant rotor controller object, the database in which
// TrackingPass objects are stored, a channel that receives updates when a
// POST, PUT, or DELETE request is made to the service, the next TrackingPass in the
// future, and a boolean value for whether the Executor is engaged (so that
// only one pass will be tracked at once)
type Executor struct {
	Rotctl        *rotor.Rotor
	DB            passes.DAO
	Updates       <-chan struct{}
	AbortCommands <-chan struct{}
	NextPass      passes.TrackingPass
	Engaged       bool
}

// Run loops the executor indefinitely, updating its NextPass attribute if
// a POST, PUT, or DELETE request is made at the service level. If a TrackingPass
// is about to start (in < 1 min) and the Executor is not currently engaged,
// it will start a goroutine that performs rotor rotation for the duration of
// the TrackingPass
func (e *Executor) Run() {
	for {
		select {
		// there was an update
		case <-e.Updates:
			e.NextPass, _ = e.DB.GetNextPass()
		// no update
		default:
			// if the next TrackingPass starts within 1 minute (and is in the future)
			if time.Until(e.NextPass.StartTime) <= 1*time.Minute && time.Now().Before(e.NextPass.StartTime) && !e.Engaged {
				e.engage()
				go func() {
					e.TrackPass(e.NextPass)
					e.NextPass, _ = e.DB.GetNextPass()
				}()
			} else {
				time.Sleep(3 * time.Second)
			}
		}
	}
}

// TrackPass carries out the automated execution of a given TrackingPass,
// rotating the rotor to ensure that it is always within 1 degree of the target
// State at a given time. Linear interpolation is used between times to estimate
// the appropriate Az/El. Additionally, there exists an abort channel that will
// stop the tracking and disengage the Executor.
func (e *Executor) TrackPass(pass passes.TrackingPass) error {
	endTime := pass.Times[len(pass.Times)-1]

	// Perform the initial rotation
	e.Rotctl.Rotate(pass.States[0])

	// Sleep until the pass starts
	for time.Now().Before(pass.StartTime) {
		time.Sleep(1 * time.Second)
	}

	// Loop until the pass is over, interpolating between state values
	idxNextTime := 0
	for now := time.Now(); now.Before(endTime) || now.Equal(endTime); now = time.Now() {
		select {
		case <-e.AbortCommands:
			e.disengage()
			return nil
		default:
			for pass.Times[idxNextTime].Before(now) {
				idxNextTime++
			}
			targetState := interpolateState(pass.States[idxNextTime], pass.States[idxNextTime-1], pass.Times[idxNextTime], pass.Times[idxNextTime-1], now)
			if math.Abs(targetState.Az-e.Rotctl.Az) > 1.0 || math.Abs(targetState.El-e.Rotctl.El) > 1.0 {
				e.Rotctl.Rotate(targetState)
			} else {
				time.Sleep(1 * time.Second)
			}
		}
	}
	e.disengage()
	return nil
}

func interpolateState(s1, s2 rotor.State, t1, t2, targetTime time.Time) rotor.State {
	ratio := targetTime.Sub(t1).Seconds() / t2.Sub(t1).Seconds()
	az := (s2.Az-s1.Az)*ratio + s1.Az
	el := (s2.El-s1.El)*ratio + s1.El
	return rotor.State{Az: az, El: el}
}

func (e *Executor) engage() {
	if !e.Engaged {
		e.Engaged = true
	}
}

func (e *Executor) disengage() {
	if e.Engaged {
		e.Engaged = false
	}
}
