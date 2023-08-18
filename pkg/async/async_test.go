package async

import (
	"testing"
)

type args struct {
	f    Anyfunc
	args []interface{}
}

type doAsyncTest struct {
	name  string
	args  args
	want  ResChan
	want1 ErrChan
}

func TestDoAsync(t *testing.T) {
	tests := []doAsyncTest{
		// TODO: Add test cases.
		{
			name: "should return string result",
			args: args{
				f: func() (any, error) {
					str, err := func() (string, error) {
						str1 := "s1"
						str2 := "s2"

						return str1 + str2, nil
					}()
					if err != nil {
						return nil, err
					}
					return str, err
				},
			},
			want:  make(ResChan),
			want1: make(ErrChan),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Async(tt.args.f)
			res, err := Await(got, got1)

			if res != "s1s2" {
				t.Errorf("DoAsync() got = %v, want %v", res, "s1s2")
			}
			if err != nil {
				t.Errorf("DoAsync() err = %v, want %v", err, "s1s2")
			}
		})
	}
}
