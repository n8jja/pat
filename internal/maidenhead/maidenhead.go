package maidenhead

import (
	"encoding/json"
	"github.com/la5nta/pat/internal/gpsd"
	"github.com/pd0mz/go-maidenhead"
)

// WatchGPSd watches the gpsd daemon for position reports
func WatchGPSd(addr string) (*gpsd.Conn, error) {
	conn, err := gpsd.Dial(addr)
	if err != nil {
		return nil, err
	}
	conn.Watch(true)
	return conn, nil
}

// GetGridSquare gets the grid square from the maidenhead package
func GetGridSquare(conn *gpsd.Conn, newLat json.Number, newLon json.Number) (string, error) {
	lat, err := newLat.Float64()
	if err != nil {
		return "", err
	}
	lon, err := newLon.Float64()
	if err != nil {
		return "", err
	}
	point := maidenhead.NewPoint(lat, lon)
	return point.GridSquare()
}

// CheckGridSquare checks the grid square in the config file to see if it matches the grid square from the gpsd daemon
func CheckGridSquare(conn *gpsd.Conn, currGridSquare string) (string, error) {
	// get tpv from gpsd.Next
	tpv, err := conn.Next()
	if err != nil {
		return "", err
	}
	// assure tpv is a TPV object
	tpvObj, ok := tpv.(*gpsd.TPV)
	if !ok {
		return "", err
	}
	// get lat and lon from tpv
	lat, err := tpvObj.Lat.Float64()
	if err != nil {
		return "", err
	}
	lon, err := tpvObj.Lon.Float64()
	if err != nil {
		return "", err
	}

	point := maidenhead.NewPoint(lat, lon)

	newGridSquare, err := point.GridSquare()
	if err != nil {
		return "", err
	}
	if newGridSquare != currGridSquare {
		return newGridSquare, nil
	}

	return currGridSquare, nil
}

// TODO: Add ability to update position via web browser information.
