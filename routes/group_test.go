package routes

import (
	"reflect"
	"testing"

	"github.com/shopspring/decimal"
)

func TestGroup_Sum(t *testing.T) {
	tests := []struct {
		name string
		g    Group
		want decimal.Decimal
	}{
		{
			name: "empty",
			g:    Group{},
			want: decimal.Zero,
		},
		{
			name: "one",
			g:    Group{Path{Amount: decimal.NewFromInt(1)}},
			want: decimal.NewFromInt(1),
		},
		{
			name: "two",
			g:    Group{Path{Amount: decimal.NewFromInt(1)}, Path{Amount: decimal.NewFromInt(2)}},
			want: decimal.NewFromInt(3),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.Sum(); !got.Equal(tt.want) {
				t.Errorf("Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_String(t *testing.T) {
	tests := []struct {
		name string
		g    Group
		want string
	}{
		{
			name: "empty",
			g:    Group{},
			want: "",
		},
		{
			name: "one",
			g: Group{Path{
				Amount: decimal.NewFromInt(1),
				Routes: Routes{1, 2},
			}},
			want: "1:1,2",
		},
		{
			name: "two",
			g: Group{
				Path{
					Amount: decimal.NewFromInt(1),
					Routes: Routes{1, 2},
				},
				Path{
					Amount: decimal.NewFromInt(2),
					Routes: Routes{3, 4},
				},
			},
			want: "1:1,2|2:3,4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseGroup(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    Group
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				s: "",
			},
			want:    Group{},
			wantErr: false,
		},
		{
			name: "one",
			args: args{
				s: "1:1,2",
			},
			want: Group{Path{
				Amount: decimal.NewFromInt(1),
				Routes: Routes{1, 2},
			}},
			wantErr: false,
		},
		{
			name: "two",
			args: args{
				s: "1:1,2|2:3,4",
			},
			want: Group{
				Path{
					Amount: decimal.NewFromInt(1),
					Routes: Routes{1, 2},
				},
				Path{
					Amount: decimal.NewFromInt(2),
					Routes: Routes{3, 4},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseGroup(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseGroup() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_Scan(t *testing.T) {
	type args struct {
		src interface{}
	}
	tests := []struct {
		name    string
		g       Group
		args    args
		wantErr bool
	}{
		{
			name: "empty",
			g:    Group{},
			args: args{
				src: "",
			},
			wantErr: false,
		},
		{
			name: "one",
			g: Group{Path{
				Amount: decimal.NewFromInt(1),
				Routes: Routes{1, 2},
			}},
			args: args{
				src: "1:1,2",
			},
			wantErr: false,
		},
		{
			name: "legacy empty",
			g:    Group{},
			args: args{
				src: "null",
			},
			wantErr: false,
		},
		{
			name: "legacy one",
			g: Group{Path{
				Routes: Routes{1, 2},
			}},
			args: args{
				src: "[1,2]",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.g.Scan(tt.args.src); (err != nil) != tt.wantErr {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
