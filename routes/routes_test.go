package routes

import (
	"math"
	"reflect"
	"testing"

	"github.com/pandodao/mtg/mtgpack"
)

func TestRoutes_String(t *testing.T) {
	tests := []struct {
		name string
		r    Routes
		want string
	}{
		{
			name: "empty",
			r:    Routes{},
			want: "",
		},
		{
			name: "one",
			r:    Routes{1},
			want: "1",
		},
		{
			name: "two",
			r:    Routes{1, 2},
			want: "1,2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseRoutes(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    Routes
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				s: "",
			},
			want:    Routes{},
			wantErr: false,
		},
		{
			name: "one",
			args: args{
				s: "1",
			},
			want:    Routes{1},
			wantErr: false,
		},
		{
			name: "two",
			args: args{
				s: "1,2",
			},
			want:    Routes{1, 2},
			wantErr: false,
		},
		{
			name: "invalid",
			args: args{
				s: "1,invalid,3",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRoutes(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRoutes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseRoutes() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoutes_Cmp(t *testing.T) {
	type args struct {
		other Routes
	}
	tests := []struct {
		name string
		r    Routes
		args args
		want int
	}{
		{
			name: "short length",
			r:    Routes{1},
			args: args{
				other: Routes{2, 1},
			},
			want: 1,
		},
		{
			name: "equal length",
			r:    Routes{2, 1},
			args: args{
				other: Routes{1, 2},
			},
			want: -1,
		},
		{
			name: "same",
			r:    Routes{1, 2},
			args: args{
				other: Routes{1, 2},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Cmp(tt.args.other); got != tt.want {
				t.Errorf("Cmp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoutes_HashString(t *testing.T) {
	tests := []struct {
		name string
		r    Routes
	}{
		{
			name: "empty",
			r:    Routes{},
		},
		{
			name: "one",
			r:    Routes{1},
		},
		{
			name: "two",
			r:    Routes{1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.r.HashString()
			r := ParseHashedRoutes(s)
			if !reflect.DeepEqual(r, tt.r) {
				t.Errorf("HashString() = %v, want %v", r, tt.r)
			}
		})
	}
}

func TestRoutes_EncodeMtg(t *testing.T) {
	type args struct {
		enc *mtgpack.Encoder
	}
	tests := []struct {
		name    string
		r       Routes
		args    args
		wantErr bool
	}{
		{
			name: "empty",
			r:    Routes{},
			args: args{
				enc: mtgpack.NewEncoder(),
			},
			wantErr: false,
		},
		{
			name: "one",
			r:    Routes{1},
			args: args{
				enc: mtgpack.NewEncoder(),
			},
			wantErr: false,
		},
		{
			name: "two",
			r:    Routes{1, 2},
			args: args{
				enc: mtgpack.NewEncoder(),
			},
			wantErr: false,
		},
		{
			name: "overflow",
			r:    Routes{1, int64(math.MaxUint16) + 1, 3},
			args: args{
				enc: mtgpack.NewEncoder(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b1 := tt.args.enc.Bytes()

			if err := tt.r.EncodeMtg(tt.args.enc); (err != nil) != tt.wantErr {
				t.Errorf("EncodeMtg() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				b2 := tt.args.enc.Bytes()
				if got, want := len(b2)-len(b1), len(tt.r)*2+1; got != want {
					t.Errorf("EncodeMtg() len = %v, want %v", got, want)
				}
			}
		})
	}
}

func TestRoutes_DecodeMtg(t *testing.T) {
	type args struct {
		dec *mtgpack.Decoder
	}
	tests := []struct {
		name    string
		r       Routes
		args    args
		wantErr bool
	}{
		{
			name: "empty",
			r:    Routes{},
			args: args{
				dec: mtgpack.NewDecoder([]byte{0}),
			},
			wantErr: false,
		},
		{
			name: "one",
			r:    Routes{1},
			args: args{
				dec: mtgpack.NewDecoder([]byte{1, 0, 1}),
			},
			wantErr: false,
		},
		{
			name: "two",
			r:    Routes{1, 2},
			args: args{
				dec: mtgpack.NewDecoder([]byte{2, 0, 1, 0, 2}),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.DecodeMtg(tt.args.dec); (err != nil) != tt.wantErr {
				t.Errorf("DecodeMtg() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
