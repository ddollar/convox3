package storage

type QueryOptions struct {
	Forward *bool
	Index   *string
	Limit   *int64
	Start   map[string]string
}

type Interface interface {
	Delete(table, id string) error
	Get(table, id string, item interface{}) error
	GetBatch(table string, ids []string, items interface{}) error
	GetIndex(table, index string, key map[string]string, items interface{}) error
	List(table string, items interface{}) error
	Put(table string, item interface{}) error
	Query(table string, key map[string]string, opts QueryOptions, items interface{}) error
}
