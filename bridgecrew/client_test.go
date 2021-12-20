package bridgecrew

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func Test_authClient(t *testing.T) {
	type args struct {
		params    RequestParams
		configure ProviderConfig
	}

	tests := []struct {
		name    string
		args    args
		want    *http.Client
		want1   *http.Request
		want2   diag.Diagnostics
		want3   bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, got3, err := authClient(tt.args.params, tt.args.configure)
			if (err != nil) != tt.wantErr {
				t.Errorf("authClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authClient() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("authClient() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("authClient() got2 = %v, want %v", got2, tt.want2)
			}
			if got3 != tt.want3 {
				t.Errorf("authClient() got3 = %v, want %v", got3, tt.want3)
			}
		})
	}
}
