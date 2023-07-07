package cnnvd

import (
	"github.com/y4ney/collect-cnnvd-vuln/internael/model"
	"github.com/y4ney/collect-cnnvd-vuln/internael/utils"
	"reflect"
	"testing"
)

func TestReqHazardLevel_Fetch(t *testing.T) {
	var hazardLevel []*model.HazardLevel
	_ = utils.ReadFile("./testdata/hazard_level.json", &hazardLevel)
	type args struct {
		retry int
	}
	tests := []struct {
		name    string
		args    args
		want    []*model.HazardLevel
		wantErr bool
	}{
		{
			name:    "test for hazard level",
			args:    args{5},
			want:    hazardLevel,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ReqHazardLevel{}
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
