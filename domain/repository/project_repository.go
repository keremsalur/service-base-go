package repository

type IProjectRepository[T any, ID any] interface {
	IBaseRepository[T, ID]
}
