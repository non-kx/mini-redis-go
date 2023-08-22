package db

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKVStore_NewKVStore(t *testing.T) {
	kvstore := NewKVStore[string](nil)

	assert.NotNil(t, kvstore)
	assert.NotNil(t, kvstore.Storage)
}

func TestKVStore_Get(t *testing.T) {
	testkey := "test_key"
	testval := "test_val"
	kvstore := NewKVStore[string](nil)
	kvstore.Storage[testkey] = testval

	val := kvstore.Get(testkey)

	assert.Equal(t, val, testval)
}

func TestKVStore_Set(t *testing.T) {
	testkey := "test_key"
	testval := "test_val"
	kvstore := NewKVStore[string](nil)

	val := kvstore.Get(testkey)

	assert.Equal(t, "", val)

	kvstore.Set(testkey, testval)
	val = kvstore.Get(testkey)

	assert.Equal(t, testval, val)
}

func TestKVStore_Delete(t *testing.T) {
	testkey := "test_key"
	testval := "test_val"
	kvstore := NewKVStore[string](nil)

	kvstore.Set(testkey, testval)
	val := kvstore.Get(testkey)

	assert.Equal(t, testval, val)

	kvstore.Delete(testkey)
	val = kvstore.Get(testkey)

	assert.Equal(t, "", val)
}

func TestLoadFromNilFile(t *testing.T) {
	storage := loadStorageFromFile[string](nil)

	assert.Equal(t, 0, len(storage))
}

func TestLoadFromFile(t *testing.T) {
	path, err := filepath.Abs("./test/test_cache.json")

	assert.Nil(t, err)

	storage := loadStorageFromFile[string](&path)

	assert.Greater(t, len(storage), 0)
}

func TestCacheStorageToFile(t *testing.T) {
	path, _ := filepath.Abs("./test/test_cache.json")
	kvstore := &KVStore[string]{
		Storage:   map[string]string{"test": "test", "a": "a", "b": "b"},
		CachePath: &path,
	}

	err := kvstore.CacheStorage()

	assert.Nil(t, err)
}
