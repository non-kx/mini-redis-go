package db

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	DefaultCachePath = "./tmp/kv_store_cache.json"
)

type IKVStore interface {
	Get(k string) []byte
	Set(k string, v []byte) error
}

type KVStore struct {
	Storage   map[string][]byte
	CachePath *string
}

func (kv *KVStore) Get(k string) []byte {
	return kv.Storage[k]
}

func (kv *KVStore) Set(k string, v []byte) {
	kv.Storage[k] = v
}

func (kv *KVStore) CacheStorage() error {
	var cachepath string
	if kv.CachePath == nil {
		cachepath = DefaultCachePath
	} else {
		cachepath = *kv.CachePath
	}

	cachedirarr := strings.Split(cachepath, "/")
	if len(cachedirarr) > 0 {
		cachedirarr = cachedirarr[:len(cachedirarr)-1]
	}
	cachedir := strings.Join(cachedirarr, "/")

	err := os.MkdirAll(cachedir, os.ModePerm)
	if err != nil {
		log.Println(err)
	}

	json, err := json.Marshal(kv.Storage)
	if err != nil {
		log.Println(err)
		return err
	}

	err = os.WriteFile(cachepath, []byte(json), os.ModePerm)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func loadStorageFromFile(path *string) map[string][]byte {
	var cachepath string
	if path == nil {
		cachepath = DefaultCachePath
	} else {
		cachepath = *path
	}

	res := make(map[string][]byte)
	dat, err := os.ReadFile(cachepath)
	if err != nil {
		fmt.Println("There is a problem load from file, start with empty storage")
		return res
	}

	err = json.Unmarshal(dat, &res)
	if err != nil {
		return res
	}

	return res
}

func InitKVStore(cachepath *string) *KVStore {
	storage := loadStorageFromFile(cachepath)

	return &KVStore{
		Storage: storage,
	}
}
