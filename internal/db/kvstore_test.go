package db

import (
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
