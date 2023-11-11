package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"testing"

	"google.golang.org/grpc"

	"github.com/DmitriiKumancev/rusprofile-grpc/pkg"
)

func TestMain(m *testing.M) {
	go startGRPCServer()
	go startHTTPServer()

	os.Exit(m.Run())
}

func equal(got *pkg.Company, want *pkg.Company) bool {
	if got == nil && want == nil {
		return true
	}

	if reflect.DeepEqual(got, &pkg.Company{}) && (want == nil) {
		return true
	}

	if (got == nil) || (want == nil) {
		return false
	}

	return got.Name == want.Name &&
		got.INN == want.INN &&
		got.KPP == want.KPP &&
		got.CEO == want.CEO
}

func Test_startGRPCServer(t *testing.T) {
	conn, err := grpc.Dial(fmt.Sprintf(":%d", *grpcPort), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("couldn't open gRPC connection: %v", err)
	}

	client := pkg.NewCompanyFinderClient(conn)

	type args struct {
		ctx context.Context
		inn *pkg.INN
	}
	tests := []struct {
		name    string
		args    args
		want    *pkg.Company
		wantErr bool
	}{
		{
			name: "Xsolla",
			args: args{
				ctx: context.Background(),
				inn: &pkg.INN{INN: "5902879646"},
			},
			want: &pkg.Company{
				INN:  "5902879646",
				KPP:  "590201001",
				Name: "ООО \"Иксолла\"",
				CEO:  "Чемоданова Валентина Игоревна",
			},
			wantErr: false,
		},
		{
			name: "Fail when incorrect INN is used",
			args: args{
				ctx: context.Background(),
				inn: &pkg.INN{INN: "123"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.ByINN(tt.args.ctx, tt.args.inn)

			if (err != nil) != tt.wantErr {
				t.Errorf("ByINN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !equal(got, tt.want) {
				t.Errorf("ByINN() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_startHTTPServer(t *testing.T) {
	url := fmt.Sprintf("http://localhost:%d/v1/company", *httpPort)

	type args struct {
		inn *pkg.INN
	}
	tests := []struct {
		name     string
		args     args
		want     *pkg.Company
		wantCode int
	}{
		{
			name: "Xsolla",
			args: args{
				inn: &pkg.INN{INN: "5902879646"},
			},
			want: &pkg.Company{
				INN:  "5902879646",
				KPP:  "590201001",
				Name: "ООО \"Иксолла\"",
				CEO:  "Чемоданова Валентина Игоревна",
			},
			wantCode: 200,
		},
		{
			name: "Fail when incorrect INN is used",
			args: args{
				inn: &pkg.INN{INN: "123"},
			},
			want:     nil,
			wantCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(fmt.Sprintf("%s/%s", url, tt.args.inn.INN))
			if err != nil {
				t.Errorf("request error: %v", err)
			}

			if tt.wantCode != resp.StatusCode {
				t.Errorf("wanted status code %d, got %d", tt.wantCode, resp.StatusCode)
			}

			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					t.Errorf("couldn't close response body: %v", err)
				}
			}(resp.Body)

			b, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("couldn't read response body: %v", err)
			}

			var got pkg.Company
			if err := json.Unmarshal(b, &got); err != nil {
				t.Errorf("couldn't unmarshal from JSON: %v", err)
			}

			if !equal(&got, tt.want) {
				t.Errorf("ByINN() got = %v, want %v", &got, tt.want)
			}
		})
	}
}
