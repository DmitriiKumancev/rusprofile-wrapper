package pkg

import (
	"context"
	"reflect"
	"testing"
)

func TestCompanyFinder_ByINN(t *testing.T) {
	type args struct {
		ctx context.Context
		inn *INN
	}
	tests := []struct {
		name    string
		args    args
		want    *Company
		wantErr bool
	}{
		{
			name: "Xsolla",
			args: args{
				ctx: context.Background(),
				inn: &INN{INN: "5902879646"},
			},
			want: &Company{
				INN:  "5902879646",
				KPP:  "590201001",
				Name: "ООО \"Иксолла\"",
				Ceo:  "Чемоданова Валентина Игоревна",
			},
			wantErr: false,
		},
		{
			name: "Fail when incorrect INN is used",
			args: args{
				ctx: context.Background(),
				inn: &INN{INN: "123"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CompanyFinder{}

			got, err := c.ByINN(tt.args.ctx, tt.args.inn)

			if (err != nil) != tt.wantErr {
				t.Errorf("ByINN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ByINN() got = %v, want %v", got, tt.want)
			}
		})
	}
}
