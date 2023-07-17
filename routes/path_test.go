package routes

import (
	"testing"

	"github.com/pandodao/mtg/mtgpack"
	"github.com/shopspring/decimal"
)

func TestPath_String(t *testing.T) {
	type fields struct {
		Amount decimal.Decimal
		Routes Routes
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "empty",
			fields: fields{
				Amount: decimal.Zero,
				Routes: Routes{},
			},
			want: "0:",
		},
		{
			name: "one",
			fields: fields{
				Amount: decimal.NewFromInt(1),
				Routes: Routes{1},
			},
			want: "1:1",
		},
		{
			name: "two",
			fields: fields{
				Amount: decimal.NewFromInt(1),
				Routes: Routes{1, 2},
			},
			want: "1:1,2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Path{
				Amount: tt.fields.Amount,
				Routes: tt.fields.Routes,
			}
			if got := p.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParsePath(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    Path
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				s: "0:",
			},
			want: Path{
				Amount: decimal.Zero,
				Routes: Routes{},
			},
			wantErr: false,
		},
		{
			name: "one",
			args: args{
				s: "1:1",
			},
			want: Path{
				Amount: decimal.NewFromInt(1),
				Routes: Routes{1},
			},
			wantErr: false,
		},
		{
			name: "invalid amount",
			args: args{
				s: "a:1",
			},
			want:    Path{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePath(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Equal(tt.want) {
				t.Errorf("ParsePath() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPath_Equal(t *testing.T) {
	type fields struct {
		Amount decimal.Decimal
		Routes Routes
	}
	type args struct {
		other Path
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "empty",
			fields: fields{
				Amount: decimal.Zero,
				Routes: Routes{},
			},
			args: args{
				other: Path{},
			},
			want: true,
		},
		{
			name: "one",
			fields: fields{
				Amount: decimal.NewFromInt(1),
				Routes: Routes{1},
			},
			args: args{
				other: Path{
					Amount: decimal.NewFromInt(1),
					Routes: Routes{1},
				},
			},
			want: true,
		},
		{
			name: "amount not equal",
			fields: fields{
				Amount: decimal.NewFromInt(1),
				Routes: Routes{1},
			},
			args: args{
				other: Path{
					Amount: decimal.NewFromInt(2),
					Routes: Routes{1},
				},
			},
			want: false,
		},
		{
			name: "routes not equal",
			fields: fields{
				Amount: decimal.NewFromInt(1),
				Routes: Routes{1},
			},
			args: args{
				other: Path{
					Amount: decimal.NewFromInt(1),
					Routes: Routes{2},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Path{
				Amount: tt.fields.Amount,
				Routes: tt.fields.Routes,
			}
			if got := p.Equal(tt.args.other); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPath_EncodeMtg(t *testing.T) {
	type fields struct {
		Amount decimal.Decimal
		Routes Routes
	}
	type args struct {
		enc *mtgpack.Encoder
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "empty",
			fields: fields{
				Amount: decimal.Zero,
				Routes: Routes{},
			},
			args: args{
				enc: mtgpack.NewEncoder(),
			},
			wantErr: false,
		},
		{
			name: "one",
			fields: fields{
				Amount: decimal.NewFromInt(1),
				Routes: Routes{1},
			},
			args: args{
				enc: mtgpack.NewEncoder(),
			},
			wantErr: false,
		},
		{
			name: "two",
			fields: fields{
				Amount: decimal.NewFromInt(1),
				Routes: Routes{1, 2},
			},
			args: args{
				enc: mtgpack.NewEncoder(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b1 := tt.args.enc.Bytes()

			p := Path{
				Amount: tt.fields.Amount,
				Routes: tt.fields.Routes,
			}

			if err := p.EncodeMtg(tt.args.enc); (err != nil) != tt.wantErr {
				t.Errorf("EncodeMtg() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				b2 := tt.args.enc.Bytes()
				if got, want := len(b2)-len(b1), 8+1+len(tt.fields.Routes)*2; got != want {
					t.Errorf("EncodeMtg() got = %v, want %v", got, want)
				}
			}
		})
	}
}
