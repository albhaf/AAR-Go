package aar

import (
	"encoding/json"
	"log"
	"net/http"
)

func getMissions() ([]Mission, error) {
	rows, err := DB.Query("SELECT id, name, world FROM missions")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]Mission, 0, 100)
	for rows.Next() {
		mission := Mission{}
		e := rows.Scan(&mission.ID, &mission.Name, &mission.World)
		if e != nil {
			return nil, e
		}

		res = append(res, mission)
	}

	return res, nil
}

func MissionsHandler(w http.ResponseWriter, r *http.Request) {
	missions, err := getMissions()

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	json.NewEncoder(w).Encode(missions)
}
