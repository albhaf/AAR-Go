package aar

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func outputMissions(w http.ResponseWriter) error {
	rows, err := DB.Query(`
		SELECT
			id,
			created_at,
			EXTRACT(
				epoch FROM (
					SELECT timestamp
					FROM events
					WHERE events.mission_id = missions.id
					ORDER BY timestamp DESC
					LIMIT 1
				) - (
					SELECT timestamp
					FROM events
					WHERE events.mission_id = missions.id
					ORDER BY timestamp ASC
					LIMIT 1
				)
			)::int AS length,
			name,
			world
		FROM missions
		ORDER BY created_at DESC
	`)

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

		mission := Mission{}
		err := rows.Scan(&mission.ID, &mission.CreatedAt, &mission.Length, &mission.Name, &mission.World)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading mission row from database: %v", err)
			continue
		}

		enc.Encode(mission)
	}

	w.Write([]byte("]"))

	return nil
}

func outputMission(missionID string, w http.ResponseWriter) error {
	row := DB.QueryRow(`
		SELECT
			id,
			name,
			world
		FROM missions
		WHERE id = $1
	`, missionID)
	mission := new(Mission)
	err := row.Scan(&mission.ID, &mission.Name, &mission.World)

	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(mission)
}

func MissionsHandler(w http.ResponseWriter, r *http.Request) {
	if err := outputMissions(w); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Error reading missions: %v", err)
	}
}

func MissionHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	missionID := params["missionId"]

	if err := outputMission(missionID, w); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Error reading mission: %v", err)
	}
}
