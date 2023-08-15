package async

import (
	"reflect"
	"testing"
)

func TestAwait(t *testing.T) {
	type args struct {
		rc ResChan
		ec ErrChan
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Await(tt.args.rc, tt.args.ec)
			if (err != nil) != tt.wantErr {
				t.Errorf("Await() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Await() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAsync(t *testing.T) {
	type args struct {
		f    Anyfunc
		args []any
	}
	tests := []struct {
		name  string
		args  args
		want  ResChan
		want1 ErrChan
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Async(tt.args.f, tt.args.args...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Async() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Async() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
