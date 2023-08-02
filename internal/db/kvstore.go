package db

type IKVStore interface {
	Get(k string) *any
	Set(k string, v *any) error
}

type KVStore struct {
	Storage map[string]*any
}

func (kv *KVStore) Get(k string) any {
	return kv.Storage[k]
}

func (kv *KVStore) Set(k string, v *any) {
	kv.Storage[k] = v
}

func InitKVStore(cachepath *string) (*KVStore, error) {
	if cachepath != nil {
		return nil, nil
	}

	return &KVStore{
		Storage: make(map[string]*any),
	}, nil
}
