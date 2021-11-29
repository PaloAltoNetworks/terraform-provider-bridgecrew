package bridgecrew

import (
	"reflect"
	"testing"
)

func Test_loadFileContent(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name    string `json:"name,omitempty"`
		args    args   `json:"args"`
		want    []byte `json:"want,omitempty"`
		wantErr bool   `json:"want_err,omitempty"`
	}{
		//{name: "jim", args: {"hjhgkjfhgfkghfjhgfjhgf"}}, want: make([]byte,1,1), wantErr: true},
		//{name: "jim", args: (), want: make([]byte,1,1), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadFileContent(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadFileContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadFileContent() = %v, want %v", got, tt.want)
			}
		})
	}
}
