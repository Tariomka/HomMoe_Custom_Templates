package zone

// Orientation controls how the bounding shape is fit and how the layout is rotated.
type Orientation struct {
	Mode                 string `json:"mode,omitempty"`
	ZeroAngleZone        string `json:"zeroAngleZone,omitempty"`
	BaseAngleMin         int    `json:"baseAngleMin,omitempty"`
	BaseAngleMax         int    `json:"baseAngleMax,omitempty"`
	RandomAngleAmplitude int    `json:"randomAngleAmplitude,omitempty"`
	RandomAngleStep      int    `json:"randomAngleStep,omitempty"`
}
