package model

import (
	"math"
)

type Spawn struct {
	EncounterID        int64   `json:"eid"`
	SpawnID            string  `json:"sid"`
	NameID             int     `json:"pid"`
	Move1ID            int     `json:"move1_id"`
	Move2ID            int     `json:"move2_id"`
	IVAttack           int     `json:"iv_attack"`
	IVDefense          int     `json:"iv_defense"`
	IVStamina          int     `json:"iv_stamina"`
	TimeToHide         int64   `json:"tth"`
	DespawnUnixSeconds int64   `json:"despawn_unix_seconds"`
	Latitude           float64 `json:"lat"`
	Longitude          float64 `json:"lon"`
	IsShiny            bool    `json:"is_shiny"`
}

func (spawn Spawn) HasIV() bool {
	// No way to set default values to a -1 flag, so check for the presence of moves.
	return spawn.Move1ID > 0
}

func (s Spawn) IVPercent() int {
	if !s.HasIV() {
		return -1
	}
	percentFloat := float64((s.IVAttack+s.IVDefense+s.IVStamina)*100) / 45
	return int(math.Floor(percentFloat + .5))
}
