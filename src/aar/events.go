package aar

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
)

type resultsWrapper interface {
	Close()
	Next() bool
	Scan(...interface{}) error
}

type eventsFetcher interface {
	GetEvents(string) (resultsWrapper, error)
}

func outputEvents(fetcher eventsFetcher, missionID string, w http.ResponseWriter) error {
	rows, err := fetcher.GetEvents(missionID)

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
		e := rows.Scan(&event.ID, &event.Data, &event.Timestamp)

		if e == nil {
			// Move properties inline to event object
			event.Player = event.Data.Player
			event.Projectile = event.Data.Projectile
			event.Unit = event.Data.Unit
			event.Vehicle = event.Data.Vehicle
			event.Data = nil

			enc.Encode(event)
		}
	}

	w.Write([]byte("]"))

	return nil
}

type dbFetcher struct {
	db *pgx.Conn
}

func (dbf *dbFetcher) GetEvents(missionID string) (resultsWrapper, error) {
	return dbf.db.Query(`
			SELECT
				id,
				data,
				timestamp
			FROM events
			WHERE mission_id = $1
			ORDER BY timestamp ASC
		`, missionID)
}

func EventsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	missionID := params["missionId"]

	fetcher := dbFetcher{DB}

	err := outputEvents(&fetcher, missionID, w)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
}
