package editor_state_v1

// MapSettings is the frozen v1 snapshot - see editorState.go. Never edit.
type MapSettings struct {
	MapSize              int  `json:"mapSize"`
	ExperimentalMapSizes bool `json:"experimentalMapSizes"`
}
