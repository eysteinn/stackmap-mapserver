package types

import (
	"strings"
	"time"
)

type UniqueProducts struct {
	Products []string `json:"product_names"`
}
type UniqueProjects struct {
	Projects []string `json:"project_names"`
}

type ProductTimes struct {
	Product string      `json:"product"`
	Times   []time.Time `json:"times"`
}

type SQLData struct {
	SQLHost string
	SQLDB   string
	SQLUser string
	SQLPass string
}
type MapLayerData struct {
	ProductTimes
	SQLData
	Project     string
	StartRange  time.Time
	EndRange    time.Time
	DefaultTime time.Time
}

func (self *MapLayerData) TimeRangeString() string {
	return self.StartRange.Format(time.RFC3339) + "/" + self.EndRange.Format(time.RFC3339)
}

func (self *MapLayerData) AllTimesString() string {
	strtimes := make([]string, len(self.Times))
	for idx, tmp := range self.Times {
		//strtimes = append(strtimes, tmp.Format(time.RFC3339))
		strtimes[idx] = tmp.Format(time.RFC3339)
	}
	return strings.Join(strtimes, ",")
}

func (self *MapLayerData) DefaultTimeString() string {
	return self.DefaultTime.Format(time.RFC3339)
}
