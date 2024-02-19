package eversendSdk

import (
	"reflect"
	"testing"
)

func TestNewEversendApp(t *testing.T) {
	type args struct {
		clientId     string
		clientSecret string
	}
	tests := []struct {
		name string
		args args
		want *Eversend
	}{
		{
			name: "Test NewEversendApp",
			args: args{
				clientId:     "client_id",
				clientSecret: "client_secret",
			},
			want: &Eversend{
				clientId:     "client_id",
				clientSecret: "client_secret",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEversendApp(tt.args.clientId, tt.args.clientSecret); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEversendApp() = %v, want %v", got, tt.want)
			}
		})
	}
}
