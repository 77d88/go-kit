package xdbutil

import (
	"reflect"
	"testing"
)

func TestParseConnection(t *testing.T) {
	type args struct {
		connStr string
	}
	tests := []struct {
		name    string
		args    args
		want    *ConnectionInfo
		wantErr bool
	}{
		{
			name: "mysql",
			args: args{
				connStr: "mysql://root:123456@127.0.0.1:3306/test",
			},
			want: &ConnectionInfo{
				Protocol: "mysql",
				Username: "root",
				Password: "123456",
				Host:     "127.0.0.1",
				Port:     3306,
				Database: "test",
			},
			wantErr: false,
		},
		{
			name: "postgresql",
			args: args{
				connStr: "postgresql://root:123456@127.0.0.1:5432/test",
			},
			want: &ConnectionInfo{
				Protocol: "postgresql",
				Username: "root",
				Password: "123456",
				Host:     "127.0.0.1",
				Port:     5433,
				Database: "test",
			},
		},
		{
			name: "reids",
			args: args{
				connStr: "redis://default:123456@127.0.0.1:6379/0",
			},
			want: &ConnectionInfo{
				Protocol: "redis",
				Username: "",
				Password: "123456",
				Host:     "127.0.0.1",
				Port:     6379,
				Database: "0",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseConnection(tt.args.connStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseConnection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseConnection() got = %v, want %v", got, tt.want)
			}
		})
	}
}
