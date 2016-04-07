package aar

import "time"

type Mission struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	World string `json:"world"`
}

type Event struct {
	ID         int32       `json:"id"`
	Data       *EventData  `json:"data,omitempty"`
	Timestamp  time.Time   `json:"timestamp"`
	Player     *Player     `json:"player,omitempty"`
	Projectile *Projectile `json:"projectile,omitempty"`
	Unit       *Unit       `json:"unit,omitempty"`
	Vehicle    *Vehicle    `json:"vehicle,omitempty"`
}

type EventData struct {
	Player     *Player     `json:"player,omitempty"`
	Projectile *Projectile `json:"projectile,omitempty"`
	Unit       *Unit       `json:"unit,omitempty"`
	Vehicle    *Vehicle    `json:"vehicle,omitempty"`
}

type Player struct {
	Name string `json:"name"`
	UID  string `json:"uid"`
}

type Position struct {
	Dir float64 `json:"dir"`
	X   float64 `json:"x"`
	Y   float64 `json:"y"`
	Z   float64 `json:"z"`
}

type Projectile struct {
	ID         string   `json:"id"`
	Position   Position `json:"position"`
	Simulation string   `json:"simulation"`
}

type Unit struct {
	ID        string   `json:"id"`
	LifeState string   `json:"life_state"`
	Name      string   `json:"name"`
	Position  Position `json:"position"`
	Side      string   `json:"side"`
}

type Vehicle struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Position   Position `json:"position"`
	Side       string   `json:"side"`
	Simulation string   `json:"simulation"`
}
