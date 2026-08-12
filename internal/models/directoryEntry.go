package models

// DirectoryEntry is one entry of a directory listing: a file, a subdirectory or
// a volume root. It is the transport type between the filesystem services and
// whatever presents the listing, so no consumer needs to touch [os.DirEntry].
type DirectoryEntry struct {
	Name  string
	Path  string
	IsDir bool
}
