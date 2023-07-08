package cnnvd

import (
	"github.com/y4ney/collect-cnnvd-vuln/internal/model"
	"github.com/y4ney/collect-cnnvd-vuln/internal/utils"
	"reflect"
	"testing"
)

func TestReqProduct_Fetch(t *testing.T) {
	var products []*model.Product
	_ = utils.ReadFile("./testdata/product.json", &products)
	type fields struct {
		ProductKeyword string
	}
	type args struct {
		retry int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*model.Product
		wantErr bool
	}{
		// TODO 官方bug
		{
			name:    "test for product",
			fields:  fields{""},
			args:    args{5},
			want:    products,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ReqProduct{
				ProductKeyword: tt.fields.ProductKeyword,
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
