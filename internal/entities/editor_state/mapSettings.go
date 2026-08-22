package editor_state

// MapSettings holds the map dimension options.
type MapSettings struct {
	MapSize              int  `json:"mapSize"`
	ExperimentalMapSizes bool `json:"experimentalMapSizes"`
}
