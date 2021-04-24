// Package sampler can be used to sample 2D trajectories based on distance as well as 3D trajectories based on time
package sampler

import (
	"errors"
	"sort"
	"time"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkt"
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/resample"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkbhex"
	geomwkt "github.com/twpayne/go-geom/encoding/wkt"
)

// Metric either time or distance used for sampling
type Metric string

const (
	// Time sampling uses hours as SamplRate unit
	Time Metric = "time"
	// Distance sampling uses meters as SamplRate unit
	Distance Metric = "distance"
)

// Instance of sampler
type Instance struct {
	// Trajectory in (E)WKT format, e.g. (LINESTRING Z (x, y, timestamp))
	WKTTrajectory string
	// Trajectory in (E)WBT format, e.g. (01020000800801058A1A4CC312424025B5548F949B5454240000)
	WKBHexTrajectory string
	Metric           Metric
	// SampleRate unit is hours for time and meters for distance
	SampleRate int
}

// Parsers

// parse2DTrajectory parses a trajectory not containing Z coordinates
// we want orb.Geometry for the resampling method which doesn't support 3D geometry,
// so we use parse3DTrajectory and flatten the coords
func (s *Instance) parse2DTrajectory() (orb.LineString, error) {
	geoLine, err := s.parse3DTrajectory()
	if err != nil {
		return nil, err
	}
	lineString := orb.LineString{}

	coords := geoLine.Coords()
	for _, c := range coords {
		lineString = append(lineString, orb.Point{c[0], c[1]})
	}

	return lineString, nil
}

// parse3DTrajectory parses a 3D trajectory where the Z coordinate contains time
func (s *Instance) parse3DTrajectory() (*geom.LineString, error) {
	var (
		geometry geom.T
		err      error
	)

	if s.WKBHexTrajectory == "" && s.WKTTrajectory == "" {
		return nil, errors.New("specify either WKTTrajectory, or WKBHexTrajectory")
	}

	if s.WKBHexTrajectory != "" {
		geometry, err = ewkbhex.Decode(s.WKBHexTrajectory)
		if err != nil {
			return nil, err
		}
	}

	if s.WKTTrajectory != "" {
		geometry, err = geomwkt.Unmarshal(s.WKTTrajectory)
		if err != nil {
			return nil, err
		}
	}

	line, ok := geometry.(*geom.LineString)
	if !ok {
		return nil, errors.New("geometry was not a valid linestring")
	}

	return line, nil
}

// Resample runs resampling based on given config
func (s *Instance) Resample() (string, error) {
	switch s.Metric {
	case Distance:
		return s.resampleDistance()
	case Time:
		return s.resampleTime()
	default:
		return "", errors.New("metric was not specified to a valid metric")
	}
}

func (s *Instance) resampleDistance() (string, error) {
	parsed, err := s.parse2DTrajectory()
	if err != nil {
		return "", err
	}

	sampled := resample.ToInterval(parsed, geo.DistanceHaversine, float64(s.SampleRate))
	return wkt.MarshalString(sampled), nil
}

// resampleTime resamples trajectory based on s.SampleRate given in hours.
// Extracts the first position within intervals based on sample rate
func (s *Instance) resampleTime() (string, error) {
	var err error

	trajectory, err := s.parse3DTrajectory()
	if err != nil {
		return "", err
	}

	intervals := s.getTimeIntervals(trajectory)
	reducedCoords := []geom.Coord{}

	coords := trajectory.Coords()

	// within each interval add the first coord to reducedCoords
	for _, interval := range intervals {
		var first *geom.Coord

		// find coord first within interval
		for i := range coords {
			coordInterval := s.roundTime(int64(coords[i][2]))
			if coordInterval == interval {
				first = &coords[i]
				break
			}
		}

		if first != nil {
			reducedCoords = append(reducedCoords, *first)
		}
	}

	// if the last coord wasn't in the reduced coords, add it
	lastReduced := reducedCoords[len(reducedCoords)-1]
	if !lastReduced.Equal(geom.XYZ, coords[len(coords)-1]) {
		reducedCoords = append(reducedCoords, coords[len(coords)-1])
	}

	if len(reducedCoords) <= 1 {
		return "", errors.New("too few points in sampled trajectory")
	}

	reduced, err := geom.NewLineString(geom.XYZ).SetCoords(reducedCoords)
	if err != nil {
		return "", err
	}

	return geomwkt.Marshal(reduced)
}

// time sampling helpers

func (s *Instance) roundTime(ts int64) time.Time {
	return time.Unix(ts, 0).UTC().Round(time.Duration(s.SampleRate) * time.Hour)
}

func (s *Instance) getTimeIntervals(trajectory *geom.LineString) []time.Time {
	times := make(map[string]time.Time)
	for _, coord := range trajectory.Coords() {
		rounded := s.roundTime(int64(coord[2]))
		times[rounded.Format("2006.01.02:15:04")] = rounded
	}

	ret := make([]time.Time, 0, len(times))
	for _, value := range times {
		ret = append(ret, value)
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Before(ret[j])
	})

	return ret
