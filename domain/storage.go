package domain

import "database/sql"

type DatabaseInstance = sql.DB
type CacheInstance = interface{}
type FileStorageInstance = interface{}

type IDataStorage[C, D, F any] interface {
	GetCacheInstance() IStorage[C]
	GetDatabaseInstance() IStorage[D]
	GetFileStorageInstance() IStorage[F]
}

type IStorage[T any] interface {
	Connect() (T, error)
	Disconnect() error
}
