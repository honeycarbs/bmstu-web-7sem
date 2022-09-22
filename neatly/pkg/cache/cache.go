package cache

type Repository interface {
	GetIterator() Iterator

	Get(uuid []byte) ([]byte, error)

	Set(key []byte, val []byte, expireIn int) error

	Delete(key []byte) (affected bool)

	EntryCount() (entryCount int64)

	HitCount() int64

	MissCount() int64
}
