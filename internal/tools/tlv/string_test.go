package tlv

import (
	"bytes"
	"encoding/binary"
	"io"
	"reflect"
	"testing"
)

func createStringTLV(str string) TypeLengthValue {
	typ := StringType

	blen := make([]byte, 4)
	binary.BigEndian.PutUint32(blen, uint32(len(str)))

	buf := []byte{typ}
	buf = append(buf, append(blen, []byte(str)...)...)

	return TypeLengthValue(buf)
}

func TestString_ReadFrom(t *testing.T) {
	str := "test_string"
	strTlv := createStringTLV(str)
	invalidStrTlv := append([]byte{BinaryType}, strTlv[1:]...)

	reader := bytes.NewReader(strTlv)
	invalidReader := bytes.NewReader(invalidStrTlv)

	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		s       *String
		args    args
		want    int64
		want2   String
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "should read to correct string",
			s:    new(String),
			args: args{
				r: reader,
			},
			want2: String(str),
			want:  int64(5 + len(str)),
		},
		{
			name: "should return err for invalid type",
			s:    new(String),
			args: args{
				r: invalidReader,
			},
			want:    1,
			want2:   *new(String),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.ReadFrom(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("String.ReadFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("String.ReadFrom() = %v, want %v", got, tt.want)
			}
			if *tt.s != String(tt.want2) {
				t.Errorf("String.ReadFrom() = %v, want %v", *tt.s, String(tt.want2))
			}
		})
	}
}

func TestString_WriteTo(t *testing.T) {
	s := "test_string"
	str := String(s)
	getStr := func(tlv TypeLengthValue) string {
		len := tlv.GetLength()
		if len > 0 {
			return string(tlv.GetValue())
		}

		return ""
	}

	tests := []struct {
		name    string
		s       *String
		want    int64
		wantW   string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "should write string correctly",
			s:       &str,
			want:    int64(5 + len(str)),
			wantW:   s,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			got, err := tt.s.WriteTo(w)
			if (err != nil) != tt.wantErr {
				t.Errorf("String.WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("String.WriteTo() = %v, want %v", got, tt.want)
			}
			if gotW := getStr(w.Bytes()); gotW != tt.wantW {
				t.Errorf("String.WriteTo() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestString_FromTLV(t *testing.T) {
	type args struct {
		tlv TypeLengthValue
	}
	tests := []struct {
		name    string
		s       *String
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.FromTLV(tt.args.tlv); (err != nil) != tt.wantErr {
				t.Errorf("String.FromTLV() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestString_ToTLV(t *testing.T) {
	tests := []struct {
		name    string
		s       *String
		want    TypeLengthValue
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.ToTLV()
			if (err != nil) != tt.wantErr {
				t.Errorf("String.ToTLV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("String.ToTLV() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_String(t *testing.T) {
	tests := []struct {
		name string
		s    *String
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("String.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
