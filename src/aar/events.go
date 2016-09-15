package aar

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func outputEvents(missionID string, w http.ResponseWriter) error {
	rows, err := DB.Query(`
		SELECT
			id,
			data,
			timestamp
		FROM events
		WHERE mission_id = $1
		ORDER BY timestamp ASC
	`, missionID)

	if err != nil {
		return err
	}
	defer rows.Close()

	enc := json.NewEncoder(w)
	w.Write([]byte("["))

	var first = true

	for rows.Next() {
		if first {
			first = false
		} else {
			w.Write([]byte(","))
		}

		event := Event{}
		err := rows.Scan(&event.ID, &event.Data, &event.Timestamp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading event row from database: %v", err)
			continue
		}

		// Move properties inline to event object
		event.Player = event.Data.Player
		event.Projectile = event.Data.Projectile
		event.Unit = event.Data.Unit
		event.Vehicle = event.Data.Vehicle
		event.Data = nil

		enc.Encode(event)
	}

	w.Write([]byte("]"))

	return nil
}

func EventsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	missionID := params["missionId"]

	if err := outputEvents(missionID, w); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Error reading events: %v", err)
	}
}
