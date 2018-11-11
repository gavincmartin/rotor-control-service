package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gavincmartin/rotor-control-service/config"
	"github.com/gavincmartin/rotor-control-service/executor"
	"github.com/gavincmartin/rotor-control-service/integrations"
	"github.com/gavincmartin/rotor-control-service/passes"
	"github.com/gavincmartin/rotor-control-service/rotor"
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
)

var (
	cfg           = config.Config{}
	db            = passes.DAO{}
	rotctl        = rotor.Rotor{State: rotor.State{Az: 0.0, El: 0.0}}
	updates       = make(chan struct{})
	abortCommands = make(chan struct{})
	passTracker   executor.Executor
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/rotor", GetRotorStateEndpoint).Methods("GET")
	r.HandleFunc("/api/rotor", SetRotorStateEndpoint).Methods("POST")
	r.HandleFunc("/api/passes", GetPassesEndpoint).Methods("GET")
	r.HandleFunc("/api/passes", AddPassEndpoint).Methods("POST")
	r.HandleFunc("/api/passes/{id}", GetPassByIDEndpoint).Methods("GET")
	r.HandleFunc("/api/passes/{id}", UpdatePassEndpoint).Methods("PUT")
	r.HandleFunc("/api/passes/{id}", DeletePassEndpoint).Methods("DELETE")
	r.HandleFunc("/api/test", TestEndpoint).Methods("GET")
	http.ListenAndServe(port(), r)
}

// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	cfg.Read()

	db.Server = cfg.Server
	db.Database = cfg.Database
	db.Connect()

	nextPass, err := db.GetNextPass()
	if err != nil {
		panic(err)
	}

	passTracker = executor.Executor{Rotctl: &rotctl, DB: db, Updates: updates, AbortCommands: abortCommands, NextPass: nextPass}
	go passTracker.Run()
}

func port() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	return ":" + port
}

// TestEndpoint allows me to test methods' behavior
func TestEndpoint(w http.ResponseWriter, r *http.Request) {
	schedule, err := db.FindAll()
	if err != nil {
		panic(err)
	}
	integrations.SendSlackSchedule(schedule)
	w.WriteHeader(http.StatusOK)

}

// GetRotorStateEndpoint delivers the Rotor's State upon a GET request
func GetRotorStateEndpoint(w http.ResponseWriter, r *http.Request) {
	safeRespondWithJSON(w, http.StatusOK, &rotctl)
}

// SetRotorStateEndpoint alters the Rotor's State upon a POST request
func SetRotorStateEndpoint(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	state := rotor.StateFromJSON(body)
	rotctl.Rotate(state)
}

// GetPassesEndpoint delivers either all TrackingPasses from MongoDB or
// TrackingPasses with a specific ID or for a specific spacecraft if a query
// parameter is added to the URL (triggered by GET request)
func GetPassesEndpoint(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var err error
	var results interface{}
	if len(q) == 0 {
		results, err = db.FindAll()
	} else {
		// Build the query
		query := make(bson.M)
		if val, ok := q["spacecraft"]; ok {
			query["spacecraft"] = val[0]
		}

		start, startDefined := q["from"]
		end, endDefined := q["to"]
		if startDefined && endDefined {
			s, _ := time.Parse(time.RFC3339, start[0])
			e, _ := time.Parse(time.RFC3339, end[0])
			query["start_time"] = bson.M{"$gte": s, "$lte": e}
		} else if startDefined {
			s, _ := time.Parse(time.RFC3339, start[0])
			query["start_time"] = bson.M{"$gte": s}
		} else if endDefined {
			e, _ := time.Parse(time.RFC3339, end[0])
			query["start_time"] = bson.M{"$lte": e}
		}
		results, err = db.FindByQuery(query)
	}
	if err != nil {
		panic(err)
	}
	respondWithJSON(w, http.StatusOK, results)
}

// AddPassEndpoint adds a TrackingPass to MongoDB upon a POST request
func AddPassEndpoint(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	pass := passes.FromJSON(body)
	pass.ID = bson.NewObjectId()
	// TODO: implement conflict check

	err = db.Insert(pass)
	if err != nil {
		panic(err)
	}

	w.Header().Add("Location", "/passes/"+pass.ID.Hex())
	respondWithJSON(w, http.StatusCreated, pass)
	go sendUpdate()
}

// GetPassByIDEndpoint retrieves a specific TrackingPass from MongoDB by ID
// upon a GET request
func GetPassByIDEndpoint(w http.ResponseWriter, r *http.Request) {
	// panics if ID isn't Mongo-compliant
	params := mux.Vars(r)
	pass, err := db.FindByID(params["id"])
	if err != nil {
		http.NotFound(w, r)
		return
	}
	respondWithJSON(w, http.StatusOK, pass)
}

// UpdatePassEndpoint updates a specific TrackingPass in MongoDB upon a PUT request
func UpdatePassEndpoint(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	pass := passes.FromJSON(body)
	if pass.ID == bson.ObjectId("") {
		pass.ID = bson.ObjectIdHex(mux.Vars(r)["id"])
	}
	err = db.Update(pass)
	if err != nil {
		panic(err)
	}
	respondWithJSON(w, http.StatusOK, pass)
	go sendUpdate()
}

// DeletePassEndpoint deletes a specific TrackingPass in MongoDB upon a DEL request
func DeletePassEndpoint(w http.ResponseWriter, r *http.Request) {
	// panics if ID isn't Mongo-compliant
	params := mux.Vars(r)
	pass, err := db.FindByID(params["id"])
	if err != nil {
		http.NotFound(w, r)
		return
	}
	err = db.Delete(pass)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusNoContent)
	go sendUpdate()
}

func respondWithJSON(w http.ResponseWriter, code int, i interface{}) {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(b)
}

func safeRespondWithJSON(w http.ResponseWriter, code int, i JSONMarshallable) {
	b := i.ToJSON()
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(b)
}

// JSONMarshallable interface is for types with the method ToJSON (used
// for concurrency-safe JSON marshalling)
type JSONMarshallable interface {
	ToJSON() []byte
}

func sendUpdate() {
	updates <- struct{}{}
}

func abortPass() {
	if passTracker.Engaged {
		abortCommands <- struct{}{}
	}
}
