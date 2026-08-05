package repositories

type IFileRepository[T any] interface {
	Load(filePath string) (T, error)
	Save(directory string, filename string, entity T) (string, error)
}
