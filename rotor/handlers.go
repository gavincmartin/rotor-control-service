package rotor

import (
	"io/ioutil"
	"net/http"
)

var rotor = Rotor{State{Az: 0.0, El: 0.0}}

// HandleFunc used as a handler function for the rotor
func HandleFunc(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodGet:
		stateJSON := rotor.State.ToJSON()
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.Write(stateJSON)
	case http.MethodPost:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		state := FromJSON(body)
		rotor.Rotate(state)
	}
}
