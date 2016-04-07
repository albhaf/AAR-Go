package aar

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func getEvents(missionId string) ([]Event, error) {
	rows, err := DB.Query("SELECT id, data, timestamp FROM events WHERE mission_id = $1 ORDER BY timestamp ASC", missionId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]Event, 0)
	for rows.Next() {
		event := Event{}
		e := rows.Scan(&event.ID, &event.Data, &event.Timestamp)
		if e != nil {
			return nil, e
		}

		res = append(res, event)
	}

	return res, nil
}

func EventsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	missionId := params["missionId"]

	events, err := getEvents(missionId)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	json.NewEncoder(w).Encode(events)
}
