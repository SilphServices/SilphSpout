package model

import (
	"math"
)

type MoveIV struct {
	Move1ID            int     `json:"m1"`
	Move2ID            int     `json:"m2"`
	Attack           int     `json:"atk"`
	Defense          int     `json:"def"`
	Stamina          int     `json:"sta"`
}

type Spawn struct {
	EncounterID        int64   `json:"eid"`
	SpawnID            string  `json:"sid"`
	NameID             int     `json:"pid"`
	IVs				   MoveIV  `json:"ivs"`
	TimeToHide         int64   `json:"tth"`
	DespawnUnixSeconds int64   `json:"dts"`
	Latitude           float64 `json:"lat"`
	Longitude          float64 `json:"lon"`
	IsShiny            bool    `json:"is_shiny"`
}

func (spawn Spawn) HasIV() bool {
	// No way to set default values to a -1 flag, so check for the presence of moves.
	return spawn.IVs.Move1ID > 0
}

func (s Spawn) IVPercent() int {
	if !s.HasIV() {
		return -1
	}
	ivs := s.IVs
	percentFloat := float64((ivs.Attack+ivs.Defense+ivs.Stamina)*100) / 45
	return int(math.Floor(percentFloat + .5))
}
