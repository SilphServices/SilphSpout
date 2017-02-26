package model

import (
	"encoding/json"
	"fmt"
	"os"
)

type Geofence interface {
	Name () string
	Inside (spawn Spawn) bool
}

type GeofenceImpl struct {
	FenceName   string
	FencePoints [][]float64
}

type AcceptAllGeofence struct {
}

func (gf AcceptAllGeofence) Inside(spawn Spawn) bool {
	return true
}

func (gf AcceptAllGeofence) Name() string {
	return "Accept All"
}

func (gf GeofenceImpl) Name() string {
	return gf.FenceName
}

func (gf GeofenceImpl) Inside(spawn Spawn) bool {
	if len(gf.FencePoints) == 0 {
		return false
	} //If there is no fence just return all is good

	x := spawn.Longitude
	y := spawn.Latitude
	points := gf.FencePoints
	inside := false
	n := len(points)
	j := 0

	for i := 0; i < n; i++ {
		j++
		if j == n {
			j = 0
		}
		if (points[i][0] < y && points[j][0] >= y) || (points[j][0] < y && points[i][0] >= y) {
			if points[i][1]+(y-points[i][0])/(points[j][0]-points[i][0])*(points[j][1]-points[i][1]) < x {
				inside = !inside
			}
		}
	}

	return inside
}

func LoadFence(pathToFenceJSON string) (geofence GeofenceImpl, err error) {

	if pathToFenceJSON == "" {
		return
	}
	fd, err := os.Open(pathToFenceJSON)
	if err != nil {
		fmt.Print(err)
	}
	decoder := json.NewDecoder(fd)
	err = decoder.Decode(&geofence)

	//Check to see if the Fence is actually closed, if not close it
	n := len(geofence.FencePoints)
	if geofence.FencePoints[0][0] != geofence.FencePoints[n-1][0] || geofence.FencePoints[0][1] != geofence.FencePoints[n-1][1] {
		geofence.FencePoints = append(geofence.FencePoints, geofence.FencePoints[0])
	}

	return
}
