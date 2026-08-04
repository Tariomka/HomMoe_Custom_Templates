package repositories

type IFileRepository[T any] interface {
	Load(filePath string) (T, error)
	Save(filePath string, entity T) (string, error)
}
