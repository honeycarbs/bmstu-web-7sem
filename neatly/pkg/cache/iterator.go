package cache

type Entry struct {
	Key   []byte
	Value []byte
}

type Iterator interface {
	Next() *Entry
}
