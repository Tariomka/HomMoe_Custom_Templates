package file_system

// IPathResolutionService turns requested locations and typed-in names into
// paths that are safe to act on. It owns every "where does this actually point"
// decision so callers never assemble a path themselves.
type IPathResolutionService interface {
	ResolveStartDirectory(preferred string) string
	ParentDirectory(current string) string
	ResolveSaveTarget(directory, filename, requiredSuffix string) (string, bool)
	PathExists(path string) bool
	DirectoryExists(path string) bool
}
