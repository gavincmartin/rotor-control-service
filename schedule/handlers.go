package schedule

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var passSchedule Schedule

// PassesHandleFunc handles GET requests for the full schedule or POST requests
// to add new passes
func PassesHandleFunc(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodGet:
		// TODO: Implement
		fmt.Println("GET called in PassesHandleFunc")
	case http.MethodPost:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		pass := FromJSON(body)
		passSchedule.AddPass(pass)
		fmt.Println(passSchedule)
		fmt.Println()
	default:
		// TODO: Implement
		fmt.Println("Not supported")
	}
}

// PassHandleFunc handles GET, PUT, and DEL for individual passes
func PassHandleFunc(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodGet:
		// TODO: Implement
		fmt.Println("GET called in PassHandleFunc")
	case http.MethodPut:
		// TODO: Implement
		fmt.Println("PUT called in PassHandleFunc")
	case http.MethodDelete:
		// TODO: Implement
		fmt.Println("DELETE called in PassHandleFunc")
	default:
		// TODO: Implement
		fmt.Println("Not supported")
	}
}
