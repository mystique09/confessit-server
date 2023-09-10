package domain

import "database/sql"

type DatabaseInstance = sql.DB
type CacheInstance = interface{}
type FileStorageInstance = interface{}

type IDataStorage interface {
	GetCacheInstance() IStorage[DatabaseInstance]
	GetDatabaseInstance() IStorage[CacheInstance]
	GetFileStorageInstance() IStorage[FileStorageInstance]
}

type IStorage[T any] interface {
	Connect() (T, error)
	Disconnect() error
}
