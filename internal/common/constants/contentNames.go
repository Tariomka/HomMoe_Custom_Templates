package constants

const (
	HubContentName       = "mandatory_content_hub"
	neutralContentPrefix = "mandatory_content_neutral_"
	sideContentPrefix    = "mandatory_content_side_"
)

func GetNeutralContentNameFor(label string) string {
	return neutralContentPrefix + label
}

func GetSideContentNameFor(label string) string {
	return sideContentPrefix + label
}
