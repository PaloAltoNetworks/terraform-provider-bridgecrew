package bridgecrew

import (
	"os"
	"reflect"
	"testing"
)

func Test_loadFileContent(t *testing.T) {
	type args struct {
		v string
	}

	os.Remove("hello-world.txt")
	contents := []byte("hello\ngo\n")
	os.WriteFile("hello-world.txt", contents, 0644)

	d1, _ := os.ReadFile("hello-world.txt")

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{name: "File Exists", args: args{"hello-world.txt"}, want: d1, wantErr: false},
		{name: "File doesn't exist", args: args{"jimbo.txt"}, want: d1, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadFileContent(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadFileContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.wantErr == true {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadFileContent() = %v, want %v", got, tt.want)
			}
		})
	}
}
