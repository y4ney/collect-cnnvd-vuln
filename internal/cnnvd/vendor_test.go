package cnnvd

import (
	"github.com/y4ney/collect-cnnvd-vuln/internael/model"
	"github.com/y4ney/collect-cnnvd-vuln/internael/utils"
	"reflect"
	"testing"
)

func TestReqVendor_Fetch(t *testing.T) {
	var data []*model.Vendor
	_ = utils.ReadFile("./testdata/vendor.json", &data)
	type fields struct {
		VendorKeyword string
	}
	type args struct {
		retry int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*model.Vendor
		wantErr bool
	}{
		// TODO 测试未通过
		{
			name:    "test for vendor",
			fields:  fields{""},
			args:    args{retry: 5},
			want:    data,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ReqVendor{
				VendorKeyword: tt.fields.VendorKeyword,
			}
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
