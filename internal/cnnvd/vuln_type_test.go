package cnnvd

import (
	"github.com/y4ney/collect-cnnvd-vuln/internal/model"
	"github.com/y4ney/collect-cnnvd-vuln/internal/utils"
	"reflect"
	"testing"
)

func TestReqVulType_Fetch(t *testing.T) {
	var vulnTypes []*model.VulnType
	_ = utils.ReadFile("./testdata/vuln-type.json", &vulnTypes)
	type args struct {
		retry int
	}
	tests := []struct {
		name    string
		args    args
		want    []*model.VulnType
		wantErr bool
	}{
		{
			name:    "test for vuln type",
			args:    args{5},
			want:    vulnTypes,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ReqVulType{}
			got, err := r.Fetch(tt.args.retry)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fetch() got = %v, want %v", got, tt.want)
			}
		})
	}
}
