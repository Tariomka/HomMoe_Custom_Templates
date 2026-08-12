package connection_editor

// connectionPairKey identifies an unordered pair of zones, so that every
// connection running between the same two zones lands in the same bucket
// regardless of which end it was authored from.
type connectionPairKey struct {
	from string
	to   string
}
