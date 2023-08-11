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

type IKVStore[T any] interface {
	Get(k string) []byte
	Set(k string, v T) error
}

type KVStore[T any] struct {
	Storage   map[string]T
	CachePath *string
}

func (kv *KVStore[T]) Get(k string) T {
	return kv.Storage[k]
}

func (kv *KVStore[T]) Set(k string, v T) {
	kv.Storage[k] = v
}

func (kv *KVStore[T]) CacheStorage() error {
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

func loadStorageFromFile[T any](path *string) map[string]T {
	var cachepath string
	if path == nil {
		cachepath = DefaultCachePath
	} else {
		cachepath = *path
	}

	res := make(map[string]T)
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

func InitKVStore[T any](cachepath *string) *KVStore[T] {
	storage := loadStorageFromFile[T](cachepath)

	return &KVStore[T]{
		Storage: storage,
	}
}
