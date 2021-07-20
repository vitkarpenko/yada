package bot

import (
	"reflect"
	"testing"
)

func Test_tokenize(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "ok",
			args: args{
				message: "о КАРТОШКА $%^@&! и  	omae WA moU    shindeiru",
			},
			want: []string{"о", "картошка", "и", "omae", "wa", "mou", "shindeiru"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tokenize(tt.args.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tokenize() = %v, want %v", got, tt.want)
			}
		})
	}
}
