package async

import (
	"errors"
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
				f: func(a ...any) (any, error) {
					str1, ok := a[0].(string)
					if !ok {
						return "", errors.New("Error no string")
					}

					str2, ok := a[1].(string)
					if !ok {
						return "", errors.New("Error no string")
					}

					return str1 + str2, nil
				},
				args: []any{"s1", "s2"},
			},
			want:  make(ResChan),
			want1: make(ErrChan),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Async(tt.args.f, tt.args.args...)
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("DoAsync() got = %v, want %v", got, tt.want)
			// }
			// if !reflect.DeepEqual(got1, tt.want1) {
			// 	t.Errorf("DoAsync() got1 = %v, want %v", got1, tt.want1)
			// }

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
