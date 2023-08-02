package db

import (
	"reflect"
	"testing"
)

func TestKVStore_Get(t *testing.T) {
	type fields struct {
		Storage map[string]*any
	}
	type args struct {
		k string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &KVStore{
				Storage: tt.fields.Storage,
			}
			if got := kv.Get(tt.args.k); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KVStore.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKVStore_Set(t *testing.T) {
	type fields struct {
		Storage map[string]*any
	}
	type args struct {
		k string
		v *any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &KVStore{
				Storage: tt.fields.Storage,
			}
			kv.Set(tt.args.k, tt.args.v)
		})
	}
}

func Test_InitKVStore(t *testing.T) {
	type args struct {
		cachepath *string
	}
	tests := []struct {
		name    string
		args    args
		want    *KVStore
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InitKVStore(tt.args.cachepath)
			if (err != nil) != tt.wantErr {
				t.Errorf("initKVStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initKVStore() = %v, want %v", got, tt.want)
			}
		})
	}
}
